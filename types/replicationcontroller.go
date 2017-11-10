package types

import (
	"k8s.io/api/core/v1"
)

type ReplicationControllerWrapper struct {
	ReplicationController ReplicationController `json:"replicationController"`
}

type ReplicationController struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Replicas        *int32            `json:"replicas,omitempty"`
	MinReadySeconds int32             `json:"minReadySeconds,omitempty"`
	PodLabels       map[string]string `json:"podLabels,omitempty"`
	Template        *Pod              `json:"template,omitempty"`

	Status *v1.ReplicationControllerStatus `json:"status,omitempty"`
}
