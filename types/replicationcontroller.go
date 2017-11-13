package types

import (
	"k8s.io/api/core/v1"
)

type ReplicationControllerWrapper struct {
	ReplicationController ReplicationController `json:"replication_controller"`
}

type ReplicationController struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Replicas        *int32 `json:"replicas,omitempty"`
	MinReadySeconds int32  `json:"ready_seconds,omitempty"`

	// Selector and the Template's Labels are expected to be equal
	// if both exist, so we standardize on using the Template's labels.
	Template *Pod `json:"template,omitempty"`

	Status *v1.ReplicationControllerStatus `json:"status,omitempty"`
}
