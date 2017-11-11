package types

import (
	apps "k8s.io/api/extensions/v1beta1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

type DeploymentWrapper struct {
	Deployment Deployment `json:"deployment"`
}

type Deployment struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Replicas       *int32              `json:"replicas,omitempty"`
	Template       Pod                 `json:"template"`
	Recreate       bool                `json:"recreate,omitempty"`
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
	MaxSurge       *intstr.IntOrString `json:"maxSurge,omitempty"`

	MinReadySeconds         int32  `json:"minReadySeconds,omitempty"`
	RevisionHistoryLimit    *int32 `json:"revisionHistoryLimit,omitempty"`
	Paused                  bool   `json:"paused,omitempty"`
	ProgressDeadlineSeconds *int32 `json:"progressDeadlineSeconds,omitempty"`

	Status *apps.DeploymentStatus `json:"status,omitempty"`
}
