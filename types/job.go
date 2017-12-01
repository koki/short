package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobWrapper struct {
	Job Job `json:"job"`
}

type Job struct {
	Version string `json:"version,omitempty"`

	PodTemplateMeta `json:",inline"`
	JobTemplate     `json:",inline"`
	JobStatus       `json:",inline"`
}

type JobTemplate struct {
	Parallelism *int32 `json:"parallelism,omitempty"`
	Completions *int32 `json:"completions,omitempty"`
	MaxRetries  *int32 `json:"max_retries,omitempty"`

	ActiveDeadlineSeconds *int64 `json:"active_deadline,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	ManualSelector *bool `json:"select_manually,omitempty"`

	TemplateMetadata *PodTemplateMeta `json:"pod_meta,omitempty"`
	PodTemplate      `json:",inline"`
}

type JobStatus struct {
	Conditions []JobCondition `json:"condition,omitempty"`
	StartTime  *metav1.Time   `json:"start_time,omitempty"`
	EndTime    *metav1.Time   `json:"end_time,omitempty"`
	Running    *int32         `json:"running,omitempty"`
	Successful *int32         `json:"successful,omitempty"`
	Failed     *int32         `json:"failed,omitempty"`
}

type JobConditionType string

const (
	JobComplete JobConditionType = "Complete"
	JobFailed   JobConditionType = "Failed"
)

type JobCondition struct {
	Type               JobConditionType `json:"type"`
	Status             ConditionStatus  `json:"status"`
	LastProbeTime      metav1.Time      `json:"timestamp,omitempty"`
	LastTransitionTime metav1.Time      `json:"last_change,omitempty"`
	Reason             string           `json:"reason,omitempty"`
	Message            string           `json:"message,omitempty"`
}
