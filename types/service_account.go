package types

type ServiceAccountWrapper struct {
	ServiceAccount ServiceAccount `json:"service_account"`
}

type ServiceAccount struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Secrets          []ObjectReference `json:"secrets,omitempty"`
	ImagePullSecrets []string          `json:"registry_secrets,omitempty"`

	AutomountServiceAccountToken *bool `json:"auto,omitempty"`
}
