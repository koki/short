package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/json"
	serrors "github.com/koki/structurederrors"
)

type ReplicaSetWrapper struct {
	ReplicaSet ReplicaSet `json:"replica_set"`
}

type ReplicaSet struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Replicas        *int32 `json:"replicas,omitempty"`
	MinReadySeconds int32  `json:"ready_seconds,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	TemplateMetadata *PodTemplateMeta `json:"pod_meta,omitempty"`
	PodTemplate      `json:",inline"`

	// Status
	ReplicaSetStatus `json:",inline"`
}

type RSSelector struct {
	Shorthand string
	Labels    map[string]string
}

type ReplicaSetStatus struct {
	ObservedGeneration int64                    `json:"generation_observed,omitempty"`
	Replicas           ReplicaSetReplicasStatus `json:"replicas_status,omitempty"`
	Conditions         []ReplicaSetCondition    `json:"condition,omitempty"`
}

type ReplicaSetReplicasStatus struct {
	Total        int32 `json:"total,omitempty"`
	FullyLabeled int32 `json:"fully_labeled,omitempty"`
	Ready        int32 `json:"ready,omitempty"`
	Available    int32 `json:"available,omitempty"`
}

type ReplicaSetConditionType string

const (
	ReplicaSetReplicaFailure ReplicaSetConditionType = "replica-failure"
)

// ReplicaSetCondition describes the state of a replica set at a certain point.
type ReplicaSetCondition struct {
	Type               ReplicaSetConditionType `json:"type"`
	Status             ConditionStatus         `json:"status"`
	LastTransitionTime metav1.Time             `json:"last_change,omitempty"`
	Reason             string                  `json:"reason,omitempty"`
	Message            string                  `json:"message,omitempty"`
}

func (s *RSSelector) UnmarshalJSON(data []byte) error {
	var str string
	strErr := json.Unmarshal(data, &str)
	if strErr == nil {
		s.Shorthand = str
		return nil
	}

	labels := map[string]string{}
	dictErr := json.Unmarshal(data, &labels)
	if dictErr != nil {
		return serrors.InvalidValueForTypeErrorf(string(data), s, "couldn't parse JSON as string or dictionary: (%s), (%s)", strErr.Error(), dictErr.Error())
	}

	s.Labels = labels
	return nil
}

func (s RSSelector) MarshalJSON() ([]byte, error) {
	if len(s.Shorthand) > 0 {
		b, err := json.Marshal(s.Shorthand)
		if err != nil {
			return nil, serrors.InvalidInstanceContextErrorf(err, s, "marshalling shorthand string to JSON")
		}

		return b, nil
	}

	b, err := json.Marshal(s.Labels)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, s, "marshalling labels dictionary to JSON")
	}

	return b, nil
}
