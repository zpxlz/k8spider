package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/esonhugh/k8spider/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Opts = struct {
	Cidr       string
	DnsServer  string
	SvcDomains []string
	Zone       string
	OutputFile string
	Verbose    int

	MultiThreadingMode bool
	ThreadingNum       int

	SkipKubeDNSCheck bool
}{}

func init() {

	RootCmd.PersistentFlags().StringVarP(&Opts.Cidr, "cidr", "c", os.Getenv("KUBERNETES_SERVICE_HOST")+"/16", "cidr like: 192.168.0.1/16")

	RootCmd.PersistentFlags().StringVarP(&Opts.DnsServer, "dns-server", "d", "", "dns server")
	RootCmd.PersistentFlags().IntVarP(&pkg.DnsTimeout, "dns-timeout", "i", 2, "dns timeout")

	RootCmd.PersistentFlags().StringSliceVarP(&Opts.SvcDomains, "svc-domains", "s", []string{}, "service domains, like: kubernetes.default,etcd.default don't add zone like svc.cluster.local")

	RootCmd.PersistentFlags().StringVarP(&Opts.Zone, "zone", "z", "cluster.local", "zone")

	RootCmd.PersistentFlags().StringVarP(&Opts.OutputFile, "output-file", "o", "", "output file")

	RootCmd.PersistentFlags().CountVarP(&Opts.Verbose, "verbose", "v", "log level (-v debug,-vv trace, info")

	RootCmd.PersistentFlags().BoolVarP(&Opts.MultiThreadingMode, "thread", "t", false, "multi threading mode, work pair with -n")
	RootCmd.PersistentFlags().IntVarP(&Opts.ThreadingNum, "thread-num", "n", 16, "threading num, default 16")

	RootCmd.PersistentFlags().BoolVarP(&Opts.SkipKubeDNSCheck, "skip-kube-dns-check", "k", false, "skip kube-dns check, force check if current environment is matched kube-dns schema")
}

var RootCmd = &cobra.Command{
	Use:   "k8spider",
	Short: "k8spider is a tool to discover k8s services",
	Long:  "k8spider is Powerful+Fast+Low Privilege Kubernetes service discovery tools via kubernetes DNS service. Currently supported service ip-port BruteForcing / AXFR Domain Transfer Dump / Coredns WildCard Dump / Pod Verified IP discovery\n\nTopics\n",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set Log Levels
		SetLogLevel(Opts.Verbose)
		// Set pkg global config
		pkg.Zone = Opts.Zone

		// debug print option data
		opt, _ := json.MarshalIndent(Opts, "", "  ")
		log.Tracef("Opts: %v", string(opt))

		if Opts.DnsServer != "" {
			pkg.NetResolver = pkg.WarpDnsServer(Opts.DnsServer)
		}

		// Check if current environment is a kubernetes cluster
		// If the command is whereisdns, which means DNS is not sure , so skip this check!
		// If SkipKubeDNSCheck is true, skip this check!
		if Opts.SkipKubeDNSCheck == false {
			if cmd.Use != "whereisdns" {
				if !pkg.CheckKubeDNS() {
					log.Warn("current environment is not a kubernetes cluster")
					os.Exit(1)
				}
			}
		} else {
			log.Tracef("kubernetes environment checking bypassed")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var LevelMap = []log.Level{
	log.InfoLevel,  // 0
	log.DebugLevel, // 1
	log.TraceLevel, // 2
}

func SetLogLevel(level int) {
	if level > 2 {
		level = 2
	}
	log.SetLevel(LevelMap[level])
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
