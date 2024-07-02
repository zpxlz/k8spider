package scanner

import (
	"net"

	"github.com/esonhugh/k8spider/pkg"
	log "github.com/sirupsen/logrus"
)

func ScanPodExist(ip net.IP, ns string) bool {
	targetHostName := pkg.IPtoPodHostName(ip.String(), ns)
	ips, err := pkg.ARecord(targetHostName)
	if err != nil {
		log.Tracef("ScanPodExist %v failed: %v", ip.String(), err)
		return false
	}
	for _, i := range ips {
		if i.String() == ip.String() {
			return true
		}
	}
	return false
}
