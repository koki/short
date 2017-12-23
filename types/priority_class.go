package types

type PriorityClassWrapper struct {
	PriorityClass PriorityClass `json:"priority_class"`
}

type PriorityClass struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Value         int32  `json:"priority,omitempty"`
	GlobalDefault bool   `json:"default",omitempty`
	Description   string `json:"description,omitempty"`
}
