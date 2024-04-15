package mutli

import (
	"net"

	"github.com/esonhugh/k8spider/define"
)

func ScanAll(subnet *net.IPNet, num int) (result <-chan []define.Record) {
	subs := NewSubnetScanner(num)
	result = ScanServiceWithChan(subs.ScanSubnet(subnet))
	return result
}
