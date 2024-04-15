package subnet

import (
	"net"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/mutli"
	"github.com/esonhugh/k8spider/pkg/printer"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	command.RootCmd.AddCommand(SubNetCmd)
}

var SubNetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "subnet is a tool to discover k8s available ip in subnet",
	Run: func(cmd *cobra.Command, args []string) {
		if command.Opts.Cidr == "" {
			log.Warn("cidr is required")
			return
		}
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
	},
}

func Run(net *net.IPNet) (records define.Records) {
	records = scanner.ScanSubnet(net)
	if records == nil || len(records) == 0 {
		log.Warnf("ScanSubnet Found Nothing")
		return
	}
	return
}

func RunMultiThread(net *net.IPNet, num int) (finalRecord define.Records) {
	scan := mutli.NewSubnetScanner(num)
	for r := range scan.ScanSubnet(net) {
		finalRecord = append(finalRecord, r...)
	}
	if len(finalRecord) == 0 {
		log.Warn("ScanSubnet Found Nothing")
		return
	}
	return
}
