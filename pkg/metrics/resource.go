package metrics

import "encoding/json"

type Resource struct {
	Namespace string            `json:"namespace"`
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Spec      map[string]string `json:"spec"`
}

func NewResource(t string) *Resource {
	return &Resource{
		Type: t,
		Spec: make(map[string]string),
	}
}

func (r *Resource) AddLabelSpec(l Label) {
	r.Spec[l.Key] = l.Value
}

func (r *Resource) AddSpec(key string, value string) {
	r.Spec[key] = value
}

type ResourceList []*Resource

func (rl *ResourceList) JSON() string {
	b, _ := json.Marshal(rl)
	return string(b)
}

func ConvertToResource(r []*MetricMatcher) []*Resource {
	var res []*Resource
	for _, m := range r {
		var resource *Resource
		if m.Name == "endpoint_address" || m.Name == "endpoint_port" {
			for i, c := range res {
				if m.FindLabel("namespace") == c.Namespace && m.FindLabel("endpoint") == c.Name {
					resource = res[i]
				} else {
					resource = NewResource("endpoint")
				}
			}
		} else {
			resource = NewResource(m.Name)
		}

		resource.Namespace = m.FindLabel("namespace")
		if m.Name == "endpoint_address" || m.Name == "endpoint_port" {
			resource.Name = m.FindLabel("endpoint")
		} else {
			resource.Name = m.FindLabel(m.Name)
		}
		// merge endpoint_address and endpoint_port
		for _, l := range m.Labels {
			if l.Key != "namespace" && l.Key != m.Name {
				resource.AddLabelSpec(l)
			}
		}
		res = append(res, resource)
	}
	return res
}
