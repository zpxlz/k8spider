package wildcard

import (
	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/pkg/printer"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	command.RootCmd.AddCommand(WildCardCmd)
}

var WildCardCmd = &cobra.Command{
	Use:   "wild",
	Short: "wild is a tool to abuse wildcard feature in kubernetes service discovery",
	Run: func(cmd *cobra.Command, args []string) {
		if command.Opts.Zone == "" {
			log.Warn("zone can't empty")
			return
		}
		printer.PrintResult(scanner.DumpWildCard(command.Opts.Zone), command.Opts.OutputFile)
	},
}
