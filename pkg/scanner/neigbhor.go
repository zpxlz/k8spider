package scanner

import (
	"github.com/esonhugh/k8spider/pkg"
)

func TestPodVerified(zone string) bool {
	iplist := []string{
		"8.8.8.8",
		"1.1.1.1",
		"114.114.114.114",
	}
	for _, ip := range iplist {
		targetHostName := pkg.IPtoPodHostName(ip, zone)
		ips, err := pkg.ARecord(targetHostName)
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
