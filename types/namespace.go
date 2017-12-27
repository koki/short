package types

type NamespaceWrapper struct {
	Namespace `json:"namespace"`
}

type Namespace struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Finalizers []FinalizerName `json:"finalizers,omitempty"`

	Phase NamespacePhase `json:"phase,omitempty"`
}

type NamespacePhase string

const (
	NamespaceActive      NamespacePhase = "Active"
	NamespaceTerminating NamespacePhase = "Terminating"
)

type FinalizerName string

const (
	FinalizerKubernetes FinalizerName = "kubernetes"
)
