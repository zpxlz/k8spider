package all

import (
	"net"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/mutli"
	"github.com/esonhugh/k8spider/pkg/post"
	"github.com/esonhugh/k8spider/pkg/printer"
	"github.com/esonhugh/k8spider/pkg/scanner"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	command.RootCmd.AddCommand(AllCmd)
}

var AllCmd = &cobra.Command{
	Use:   "all",
	Short: "all is a tool to discover k8s services and available ip in subnet",
	Run: func(cmd *cobra.Command, args []string) {
		if command.Opts.Cidr == "" {
			log.Warn("cidr is required")
			return
		}
		// Wildcard
		records := scanner.DumpWildCard(command.Opts.Zone)
		if records != nil {
			printer.PrintResult(records, command.Opts.OutputFile)
		}
		// AXFR Dumping
		records, err := scanner.DumpAXFR(dns.Fqdn(command.Opts.Zone), "ns.dns."+command.Opts.Zone+":53")
		if err == nil {
			printer.PrintResult(records, command.Opts.OutputFile)
		} else {
			log.Errorf("Transfer failed: %v", err)
		}
		// Service Discovery
		ipNets, err := pkg.ParseStringToIPNet(command.Opts.Cidr)
		if err != nil {
			log.Warnf("ParseStringToIPNet failed: %v", err)
			return
		}
		var finalRecord define.Records
		if command.Opts.MultiThreadingMode {
			finalRecord = RunMultiThread(ipNets, command.Opts.ThreadingNum)
		} else {
			finalRecord = Run(ipNets)
		}
		printer.PrintResult(finalRecord, command.Opts.OutputFile)
		PostRun(finalRecord)
	},
}

func Run(net *net.IPNet) (finalRecord define.Records) {
	var records define.Records = scanner.ScanSubnet(net)
	if records == nil || len(records) == 0 {
		log.Warnf("ScanSubnet Found Nothing")
		return
	}
	records = scanner.ScanSvcForPorts(records)
	return records
}

func RunMultiThread(net *net.IPNet, count int) (finalRecord define.Records) {
	scan := mutli.ScanAll(net, count)
	for r := range scan {
		finalRecord = append(finalRecord, r...)
	}
	return
}

func PostRun(finalRecord define.Records) {
	log.Info("Extract Namespaces: ")
	list := post.RecordsDumpNameSpace(finalRecord, command.Opts.Zone)
	for _, ns := range list {
		log.Infof("Namespace: %s", ns)
	}
	log.Info("Extract Service: ")
	list = post.RecordsDumpFullService(finalRecord, command.Opts.Zone)
	for _, svc := range list {
		log.Infof("Service: %s", svc)
	}
}
