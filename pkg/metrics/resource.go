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
		Spec: make(map[string]string, 4),
	}
}

func (r *Resource) AddLabelSpec(l Label) {
	r.Spec[l.Key] = l.Value
}

func (r *Resource) AddSpec(key string, value string) {
	r.Spec[key] = value
}

type ResourceList []*Resource

func (r *Resource) JSON() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (rl *ResourceList) JSON() string {
	var res = ""
	for _, r := range *rl {
		res += r.JSON() + "\n"
	}
	return res
}

func ConvertToResource(r []*MetricMatcher) []*Resource {
	var res []*Resource
	for _, m := range r {
		var resource *Resource
		var addFlag = true
		if m.Name == "endpoint_address" || m.Name == "endpoint_port" {
			for i, c := range res {
				if m.FindLabel("namespace") == c.Namespace && m.FindLabel("endpoint") == c.Name {
					resource = res[i]
					addFlag = false
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
		resource.Namespace = m.FindLabel("namespace")

		// merge endpoint_address and endpoint_port
		for _, l := range m.Labels {
			if l.Key != "namespace" && l.Key != resource.Type {
				resource.AddLabelSpec(l)
			}
		}
		if addFlag {
			res = append(res, resource)
		}
	}
	return res
}
