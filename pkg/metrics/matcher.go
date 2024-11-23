package metrics

import "errors"

type MatchRules []*MetricMatcher

func GenerateMatchRules() MatchRules {
	return make(MatchRules, 0)
}

func DefaultMatchRules() MatchRules {
	return []*MetricMatcher{
		NewMetricMatcher("configmap").AddLabel("namespace").AddLabel("configmap"),
		NewMetricMatcher("secret").AddLabel("namespace").AddLabel("secret"),

		NewMetricMatcher("node").AddLabel("node").AddLabel("kernel_version").
			AddLabel("os_image").AddLabel("container_runtime_version").
			AddLabel("provider_id").AddLabel("internal_ip"),

		NewMetricMatcher("pod").AddLabel("namespace").AddLabel("pod").AddLabel("node").
			AddLabel("host_ip").AddLabel("pod_ip"),
		NewMetricMatcher("container").SetHeader("kube_pod_container_info").AddLabel("namespace").
			AddLabel("pod").AddLabel("container").AddLabel("image"),
		NewMetricMatcher("cronjob").AddLabel("namespace").AddLabel("cronjob").
			AddLabel("schedule"),

		NewMetricMatcher("service_account").SetHeader("kube_pod_service_account").
			AddLabel("namespace").AddLabel("pod").AddLabel("service_account"),

		NewMetricMatcher("service").AddLabel("namespace").AddLabel("service").
			AddLabel("cluster_ip"),
		NewMetricMatcher("endpoint_address").SetHeader("kube_endpoint_address").
			AddLabel("namespace").AddLabel("endpoint").AddLabel("ip"),
		NewMetricMatcher("endpoint_port").SetHeader("kube_endpoint_ports").
			AddLabel("namespace").AddLabel("endpoint").AddLabel("port_number"),
	}
}

func (m MatchRules) Compile() error {
	var err error = nil
	for i := range m {
		e := m[i].Compile()
		if e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}

func (m MatchRules) Match(target string) (*MetricMatcher, error) {
	for _, r := range m {
		_, e := r.Match(target)
		if e != nil {
			continue
		} else {
			return r, nil
		}
	}
	return nil, errors.New("no match found")
}
