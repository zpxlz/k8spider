package mutli

import (
	"net"
	"sync"

	"github.com/esonhugh/k8spider/define"
)

type NeighborScanner struct {
	wg    *sync.WaitGroup
	count int
}

func NewNeighborScanner(threading ...int) *NeighborScanner {
	if len(threading) == 0 {
		return &NeighborScanner{
			wg: new(sync.WaitGroup),
		}
	} else {
		return &NeighborScanner{
			wg:    new(sync.WaitGroup),
			count: threading[0],
		}
	}
}

func (s *NeighborScanner) ScanNeighbor(subnet *net.IPNet) <-chan []define.Record {
	return nil // todo: implement this
}
