package types

import (
	"k8s.io/apimachinery/pkg/util/intstr"
)

type DaemonSetWrapper struct {
	DaemonSet DaemonSet `json:"daemon_set"`
}

type DaemonSet struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	OnDelete       bool                `json:"replace_on_delete,omitempty"`
	MaxUnavailable *intstr.IntOrString `json:"max_unavailable,omitempty"`

	MinReadySeconds      int32  `json:"min_ready,omitempty"`
	RevisionHistoryLimit *int32 `json:"max_revs,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	TemplateMetadata *PodTemplateMeta `json:"pod_meta,omitempty"`
	PodTemplate      `json:",inline"`

	// Status
	DaemonSetStatus `json:",inline"`
}

type DaemonSetStatus struct {
	ObservedGeneration   int64  `json:"generation_observed,omitempty"`
	NumNodesScheduled    int32  `json:"num_nodes_scheduled,omitempty"`
	NumNodesMisscheduled int32  `json:"num_nodes_misscheduled,omitempty"`
	NumNodesDesired      int32  `json:"num_nodes_desired,omitempty"`
	NumReady             int32  `json:"num_ready,omitempty"`
	NumUpdated           int32  `json:"num_updated,omitempty"`
	NumAvailable         int32  `json:"num_available,omitempty"`
	NumUnavailable       int32  `json:"num_unavailable,omitempty"`
	CollisionCount       *int32 `json:"hash_collisions,omitempty"`
}
