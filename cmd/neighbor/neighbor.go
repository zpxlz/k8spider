package neighbor

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	cmdx "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/mutli"
	"github.com/esonhugh/k8spider/pkg/printer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Opts = struct {
	NamespaceWordlist string
	NamespaceList     []string
	PodCidr           string
}{}

func init() {
	cmdx.RootCmd.AddCommand(NeighborCmd)
	NeighborCmd.Flags().StringVar(&Opts.NamespaceWordlist, "ns-file", "", "namespace wordlist file")
	NeighborCmd.Flags().StringSliceVar(&Opts.NamespaceList, "ns", []string{}, "namespace list")
	NeighborCmd.Flags().StringVarP(&Opts.PodCidr, "pod-cidr", "p", defaultPodCidr(), "pod cidr list, watch out for the network interface name, default is eth0")
}

func defaultPodCidr() string {
	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		if i.Name == "eth0" {
			addrs, _ := i.Addrs()
			if addrs != nil || len(addrs) > 0 {
				ip := strings.Split(addrs[0].String(), "/")[0]
				return fmt.Sprintf("%v/16", ip)
			}
		}
	}
	return "10.0.0.1/16"
}

var NeighborCmd = &cobra.Command{
	Use:     "neighbor",
	Short:   "neighbor is a tool to discover k8s pod and available ip in subnet (require k8s coredns with pod verified config)",
	Aliases: []string{"n", "nei"},
	Run: func(cmd *cobra.Command, args []string) {
		if !pkg.TestPodVerified() {
			log.Fatalf("k8s coredns with pod verified config could not be set")
		}
		if Opts.NamespaceWordlist != "" {
			f, e := os.OpenFile(Opts.NamespaceWordlist, os.O_RDONLY, 0666)
			if e != nil {
				log.Fatalf("open file %v failed: %v", Opts.NamespaceWordlist, e)
			}
			defer f.Close()
			fileScanner := bufio.NewScanner(f)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				Opts.NamespaceList = append(Opts.NamespaceList, fileScanner.Text())
			}
		}
		log.Tracef("namespace list: %v", Opts.NamespaceList)
		ipNets, err := pkg.ParseStringToIPNet(Opts.PodCidr)
		if err != nil {
			log.Warnf("ParseStringToIPNet failed: %v", err)
			return
		}
		r := RunMultiThread(Opts.NamespaceList, ipNets, cmdx.Opts.ThreadingNum)
		printer.PrintResult(r, cmdx.Opts.OutputFile)
	},
}

func RunMultiThread(ns []string, net *net.IPNet, num int) (finalRecord define.Records) {
	scan := mutli.ScanNeighbor(ns, net, num)
	for r := range scan {
		finalRecord = append(finalRecord, r...)
	}
	if len(finalRecord) == 0 {
		log.Warn("ScanSubnet Found Nothing")
		return
	}
	return
}
