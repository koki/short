package types

type Affinity struct {
	NodeAffinity    string `json:"node,omitempty"`
	PodAffinity     string `json:"pod,omitempty"`
	PodAntiAffinity string `json:"!pod,omitempty"`
	Topology        string `json:"topology,omitempty"`
	Namespace       string `json:"namespace,omitempty"`
}
