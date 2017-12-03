package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Selector map[string]string `json:"selector,omitempty"`

	// Template fields
	TemplateMetadata *PodTemplateMeta `json:"pod_meta,omitempty"`
	PodTemplate      `json:",inline"`

	// Status
	ReplicationControllerStatus `json:",inline"`
}

type ReplicationControllerStatus struct {
	ObservedGeneration int64                                `json:"generation_observed,omitempty"`
	Replicas           *ReplicationControllerReplicasStatus `json:"replicas_status,omitempty"`
	Conditions         []ReplicationControllerCondition     `json:"condition,omitempty"`
}

type ReplicationControllerReplicasStatus struct {
	Total        int32 `json:"total,omitempty"`
	FullyLabeled int32 `json:"fully_labeled,omitempty"`
	Ready        int32 `json:"ready,omitempty"`
	Available    int32 `json:"available,omitempty"`
}

type ReplicationControllerConditionType string

const (
	ReplicationControllerReplicaFailure ReplicationControllerConditionType = "replica-failure"
)

// ReplicationControllerCondition describes the state of a replica set at a certain point.
type ReplicationControllerCondition struct {
	Type               ReplicationControllerConditionType `json:"type"`
	Status             ConditionStatus                    `json:"status"`
	LastTransitionTime metav1.Time                        `json:"last_change,omitempty"`
	Reason             string                             `json:"reason,omitempty"`
	Message            string                             `json:"message,omitempty"`
}
