package types

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CronJobWrapper struct {
	CronJob CronJob `json:"cron_job"`
}

type CronJob struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Schedule                string `json:"schedule,omitempty"`
	Suspend                 *bool  `json:"suspend,omitempty"`
	StartingDeadlineSeconds *int64 `json:"start_deadline,omitempty"`

	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrency,omitempty"`

	MaxSuccessHistory *int32 `json:"max_success_history,omitempty"`
	MaxFailureHistory *int32 `json:"max_failure_history,omitempty"`

	TemplateMetadata *PodTemplateMeta `json:"job_meta,omitempty"`
	JobTemplate      `json:",inline"`

	// Status
	CronJobStatus `json:",inline"`
}

type ConcurrencyPolicy string

const (
	AllowConcurrent   ConcurrencyPolicy = "allow"
	ForbidConcurrent  ConcurrencyPolicy = "forbid"
	ReplaceConcurrent ConcurrencyPolicy = "replace"
)

type CronJobStatus struct {
	Active        []v1.ObjectReference `json:"active,omitempty"`
	LastScheduled *metav1.Time         `json:"last_scheduled,omitempty"`
}
