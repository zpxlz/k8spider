package metrics

import "encoding/json"

type Resource struct {
	Namespace string              `json:"namespace"`
	Type      string              `json:"type"`
	Name      string              `json:"name"`
	Spec      map[string][]string `json:"spec"`
}

func NewResource(t string) *Resource {
	return &Resource{
		Type: t,
		Spec: make(map[string][]string, 4),
	}
}

func (r *Resource) AddLabelSpec(l Label) {
	if l.Value == "" {
		return
	}
	if _, ok := r.Spec[l.Key]; !ok {
		r.Spec[l.Key] = make([]string, 0)
	}
	r.Spec[l.Key] = append(r.Spec[l.Key], l.Value)
}

func (r *Resource) AddSpec(key string, value string) {
	r.AddLabelSpec(Label{Key: key, Value: value})
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

		resourceType := m.Name
		if m.Name == "endpoint_address" || m.Name == "endpoint_port" {
			for i, c := range res {
				if m.FindLabel("namespace") == c.Namespace && m.FindLabel("endpoint") == c.Name {
					resource = res[i]
					addFlag = false
				}
			}
			resourceType = "endpoint"
		}

		if addFlag {
			resource = NewResource(resourceType)
		}

		resource.Namespace = m.FindLabel("namespace")
		resource.Name = m.FindLabel(resourceType)

		// merge endpoint_address and endpoint_port
		for _, l := range m.Labels {
			if l.Key != "namespace" && l.Key != resourceType {
				resource.AddLabelSpec(l)
			}
		}
		if addFlag {
			res = append(res, resource)
		}
	}
	return res
}
