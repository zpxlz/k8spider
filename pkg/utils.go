package pkg

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

var NetResolver *net.Resolver = net.DefaultResolver

var Zone string // Zone is the domain name of the cluster

func ParseStringToIPNet(s string) (ipnet *net.IPNet, err error) {
	_, ipnet, err = net.ParseCIDR(s)
	return
}

func ParseIPNetToIPs(ipv4Net *net.IPNet) (ips []net.IP) {
	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	// loop through addresses as uint32
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip)
	}
	return
}

func PTRRecord(ip net.IP) []string {
	names, err := NetResolver.LookupAddr(context.Background(), ip.String())
	if err != nil {
		log.Debugf("LookupAddr failed: %v", err)
		return nil
	}
	return names
}

func SRVRecord(svcDomain string) (string, []*net.SRV, error) {
	cname, srvs, err := NetResolver.LookupSRV(context.Background(), "", "", svcDomain)
	return cname, srvs, err
}

func ARecord(domain string) (ips []net.IP, err error) {
	// ips, err = NetResolver.LookupIP()
	ips, err = NetResolver.LookupIP(context.Background(), "ip", domain)
	return
}

func IPtoPodHostName(ip, namespace string) string {
	return fmt.Sprintf("%s.%s.pod.%s", strings.ReplaceAll(ip, ".", "-"), namespace, Zone)
}

func TestPodVerified() bool {
	iplist := []string{
		"8.8.8.8",
		"1.1.1.1",
		"114.114.114.114",
	}
	for _, ip := range iplist {
		targetHostName := IPtoPodHostName(ip, "kube-system")
		log.Tracef("test if record %v is ip: %v", targetHostName, ip)
		ips, err := ARecord(targetHostName)
		if err != nil {
			continue
		}
		// all of this ip should not return correct ip if verified is set.
		for _, i := range ips {
			if i.String() == ip {
				return false
			}
		}
	}
	return true
}

func TXTRecord(domain string) (txts []string, err error) {
	txts, err = NetResolver.LookupTXT(context.Background(), domain)
	return
}

// https://github.com/kubernetes/dns/blob/master/docs/specification.md
// CheckKubernetes checks if the current environment is a kubernetes cluster
func CheckKubernetes() bool {
	_, err := ARecord("kubernetes.default.svc." + Zone)
	if err != nil {
		return false
	}
	t, err := TXTRecord("dns-version." + Zone)
	if err != nil {
		return false
	}
	log.Infof("dns-version: %v", strings.Join(t, ","))
	return true
}
