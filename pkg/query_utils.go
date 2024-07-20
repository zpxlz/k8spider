package pkg

import (
	"context"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	DnsTimeout  = 2
	NetResolver = DefaultResolver()
	Zone        string // Zone is the domain name of the cluster
)

type SpiderResolver struct {
	dns string
	ctx context.Context
	r   *net.Resolver
}

func DefaultResolver() *SpiderResolver {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(DnsTimeout)*time.Second) // I don't think if a inside cluster dns query has more than 2s latency.
	return &SpiderResolver{
		dns: "<default-dns>",
		r:   net.DefaultResolver,
		ctx: ctx,
	}
}

func WarpDnsServer(dnsServer string) *SpiderResolver {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(DnsTimeout)*time.Second)
	return &SpiderResolver{
		dns: dnsServer,
		r: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, network, dnsServer)
			},
		},
		ctx: ctx,
	}
}

func (s *SpiderResolver) CurrentDNS() string {
	return s.dns
}

func (s *SpiderResolver) PTRRecord(ip net.IP) []string {
	names, err := s.r.LookupAddr(s.ctx, ip.String())
	if err != nil {
		log.Debugf("LookupAddr failed: %v", err)
		return nil
	}
	return names
}

func PTRRecord(ip net.IP) []string {
	return NetResolver.PTRRecord(ip)
}

func (s *SpiderResolver) SRVRecord(svcDomain string) (string, []*net.SRV, error) {
	cname, srvs, err := s.r.LookupSRV(s.ctx, "", "", svcDomain)
	return cname, srvs, err
}

func (s *SpiderResolver) CustomSRVRecord(svcDomain string, service, proto string) (string, []*net.SRV, error) {
	cname, srvs, err := s.r.LookupSRV(s.ctx, service, proto, svcDomain)
	return cname, srvs, err
}

func SRVRecord(svcDomain string) (string, []*net.SRV, error) {
	return NetResolver.SRVRecord(svcDomain)
}

func (s *SpiderResolver) ARecord(domain string) ([]net.IP, error) {
	return s.r.LookupIP(s.ctx, "ip", domain)
}

func ARecord(domain string) (ips []net.IP, err error) {
	return NetResolver.ARecord(domain)
}

func (s *SpiderResolver) TXTRecord(domain string) ([]string, error) {
	return s.r.LookupTXT(s.ctx, domain)
}

func TXTRecord(domain string) (txts []string, err error) {
	return NetResolver.TXTRecord(domain)
}
