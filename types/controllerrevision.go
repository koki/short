package types

import (
	"k8s.io/apimachinery/pkg/runtime"
)

type ControllerRevisionWrapper struct {
	ControllerRevision ControllerRevision `json:"controller_revision"`
}

type ControllerRevision struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Data runtime.RawExtension `json:"data,omitempty"`

	// Revision indicates the revision of the state represented by Data.
	Revision int64 `json:"revision"`
}
