package cmd

import (
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
	Verbose    string

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
	RootCmd.PersistentFlags().StringVarP(&Opts.Verbose, "verbose", "v", "info", "log level (debug,info,trace,warn,error,fatal,panic)")
	RootCmd.PersistentFlags().BoolVarP(&Opts.MultiThreadingMode, "thread", "t", false, "multi threading mode, work pair with -n")
	RootCmd.PersistentFlags().IntVarP(&Opts.ThreadingNum, "thread-num", "n", 16, "threading num, default 16")

	RootCmd.PersistentFlags().BoolVarP(&Opts.SkipKubeDNSCheck, "skip-kube-dns-check", "k", false, "skip kube-dns check, force check if current environment is matched kube-dns schema")
}

var RootCmd = &cobra.Command{
	Use:   "k8spider",
	Short: "k8spider is a tool to discover k8s services",
	Long:  "k8spider is a tool to discover k8s services",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(Opts.Verbose)
		if Opts.DnsServer != "" {
			pkg.NetResolver = pkg.WarpDnsServer(Opts.DnsServer)
		}
		if Opts.SkipKubeDNSCheck == false { // Not Skip
			if pkg.CheckKubeDNS() {
				log.Warn("current environment is not a kubernetes cluster")
				os.Exit(1)
			}
		}
		pkg.Zone = Opts.Zone
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func SetLogLevel(level string) {
	switch level {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
