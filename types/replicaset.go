package types

import (
	apps "k8s.io/api/extensions/v1beta1"
)

type ReplicaSetWrapper struct {
	ReplicaSet ReplicaSet `json:"replicaSet"`
}

type ReplicaSet struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Replicas        *int32 `json:"replicas,omitempty"`
	MinReadySeconds int32  `json:"minReadySeconds,omitempty"`
	PodSelector     string `json:"podSelector,omitempty"`
	Template        *Pod   `json:"template,omitempty"`

	Status *apps.ReplicaSetStatus `json:"status,omitempty"`
}
