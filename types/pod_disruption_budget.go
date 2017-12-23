package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/util/floatstr"
)

type PodDisruptionBudgetWrapper struct {
	PodDisruptionBudget PodDisruptionBudget `json:"pdb"`
}

type PodDisruptionBudget struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	MaxEvictionsAllowed *floatstr.FloatOrString `json:"max_evictions,omitempty"`
	MinPodsRequired     *floatstr.FloatOrString `json:"min_pods,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	// Status
	PodDisruptionBudgetStatus `json:",inline"`
}

type PodDisruptionBudgetStatus struct {
	ObservedGeneration    int64                  `json:"generation_observed,omitempty"`
	DisruptedPods         map[string]metav1.Time `json:"disrupted_pods,omitempty"`
	PodDisruptionsAllowed int32                  `json:"allowed_disruptions,omitempty"`
	CurrentHealthy        int32                  `json:"current_healthy_pods,omitempty"`
	DesiredHealthy        int32                  `json:"desired_healthy_pods,omitempty"`
	ExpectedPods          int32                  `json:"expected_pods,omitempty"`
}
