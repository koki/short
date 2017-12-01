package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	Recreate       bool                `json:"recreate,omitempty"`
	MaxUnavailable *intstr.IntOrString `json:"max_unavailable,omitempty"`
	MaxSurge       *intstr.IntOrString `json:"max_extra,omitempty"`

	MinReadySeconds         int32  `json:"min_ready,omitempty"`
	RevisionHistoryLimit    *int32 `json:"max_revs,omitempty"`
	Paused                  bool   `json:"paused,omitempty"`
	ProgressDeadlineSeconds *int32 `json:"progress_deadline,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	TemplateMetadata *PodTemplateMeta `json:"pod_meta,omitempty"`
	PodTemplate      `json:",inline"`

	// Status
	DeploymentStatus `json:",inline"`
}

type DeploymentStatus struct {
	ObservedGeneration int64                    `json:"generation_observed,omitempty"`
	Replicas           DeploymentReplicasStatus `json:"replicas_status,omitempty"`
	Conditions         []DeploymentCondition    `json:"condition,omitempty"`
	CollisionCount     *int32                   `json:"hash_collisions,omitempty"`
}

type DeploymentReplicasStatus struct {
	Total       int32 `json:"total,omitempty"`
	Updated     int32 `json:"updated,omitempty"`
	Ready       int32 `json:"ready,omitempty"`
	Available   int32 `json:"available,omitempty"`
	Unavailable int32 `json:"unavailable,omitempty"`
}

type DeploymentConditionType string

const (
	DeploymentAvailable      DeploymentConditionType = "available"
	DeploymentProgressing    DeploymentConditionType = "progressing"
	DeploymentReplicaFailure DeploymentConditionType = "replica-failure"
)

type DeploymentCondition struct {
	Type               DeploymentConditionType `json:"type"`
	Status             ConditionStatus         `json:"status"`
	LastUpdateTime     metav1.Time             `json:"timestamp,omitempty"`
	LastTransitionTime metav1.Time             `json:"last_change,omitempty"`
	Reason             string                  `json:"reason,omitempty"`
	Message            string                  `json:"message,omitempty"`
}
