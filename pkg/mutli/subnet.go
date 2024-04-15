package mutli

import (
	"net"
	"sync"
	"time"

	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
)

type SubnetScanner struct {
	wg    *sync.WaitGroup
	count int
}

func NewSubnetScanner(threading ...int) *SubnetScanner {
	if len(threading) == 0 {
		return &SubnetScanner{
			wg: new(sync.WaitGroup),
		}
	} else {
		return &SubnetScanner{
			wg:    new(sync.WaitGroup),
			count: threading[0],
		}
	}
}

func (s *SubnetScanner) ScanSubnet(subnet *net.IPNet) <-chan []define.Record {
	if subnet == nil {
		log.Debugf("subnet is nil")
		return nil
	}
	out := make(chan []define.Record, 100)
	go func() {
		log.Debugf("splitting subnet into 16 pices")
		// if subnets, err := pkg.SubnetShift(subnet, 4); err != nil {
		if subnets, err := pkg.SubnetInto(subnet, s.count); err != nil {
			log.Errorf("Subnet split into %v failed, fallback to single mode, reason: %v", s.count, err)
			go s.scan(subnet, out)
		} else {
			log.Debugf("Subnet split into %v success", len(subnets))
			for _, sn := range subnets {
				go s.scan(sn, out)
			}
		}
		time.Sleep(10 * time.Millisecond) // wait for all goroutines to start
		s.wg.Wait()
		close(out)
	}()
	return out
}

func (s *SubnetScanner) scan(subnet *net.IPNet, to chan []define.Record) {
	s.wg.Add(1)
	// to <- scanner.ScanSubnet(subnet)
	for _, ip := range pkg.ParseIPNetToIPs(subnet) {
		to <- scanner.ScanSingleIP(ip)
	}
	s.wg.Done()
}
