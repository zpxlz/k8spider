package post

import (
	"strings"

	"github.com/esonhugh/k8spider/define"
	"github.com/miekg/dns"
)

func RecordsDumpFullService(r []define.Record, zone string) []string {
	var result []string
	for _, record := range r {
		if record.SvcDomain != "" && IsServiceFormat(record.SvcDomain, zone) {
			result = append(result, record.SvcDomain)
		}
		for _, srv := range record.SrvRecords {
			for _, s := range srv.Srv {
				if IsServiceFormat(s.Target, zone) {
					result = append(result, s.Target)
				}
			}
		}
		R := ExtraParser(record.Extra)
		if R != "" {
			if IsServiceFormat(R, zone) {
				result = append(result, R)
			}
		}
	}
	return UniqueSlice(result)
}

func IsServiceFormat(domain string, zone string) bool {
	zonelen := len(strings.Split(dns.Fqdn(zone), "."))
	dn := ReverseSlice(strings.Split(dns.Fqdn(domain), "."))
	if len(dn) > 4 {
		if dn[zonelen] == "svc" { // Check if it is a service domain
			return true
		}
	}
	return false
}

func RecordsDumpNameSpace(r []define.Record, zone string) []string {
	result := RecordsDumpFullService(r, zone)
	for i, record := range result {
		result[i] = GetNamespaceFromDomain(record, zone)
	}
	return UniqueSlice(result)
}

func ExtraParser(record string) string {
	return strings.Split(record, " ")[0]
}

func GetNamespaceFromDomain(domain string, zone string) string {
	zonelen := len(strings.Split(dns.Fqdn(zone), "."))
	dn := ReverseSlice(strings.Split(dns.Fqdn(domain), "."))
	if len(dn) > 4 {
		if dn[zonelen] == "svc" { // Check if it is a service domain
			return dn[1+zonelen]
		}
	}
	return ""
}
