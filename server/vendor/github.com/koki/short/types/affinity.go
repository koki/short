package types

type Affinity struct {
	NodeAffinity    string   `json:"node,omitempty"`
	PodAffinity     string   `json:"pod,omitempty"`
	PodAntiAffinity string   `json:"anti_pod,omitempty"`
	Topology        string   `json:"topology,omitempty"`
	Namespaces      []string `json:"namespaces,omitempty"`
}
