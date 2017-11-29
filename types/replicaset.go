package types

import (
	"encoding/json"

	apps "k8s.io/api/apps/v1beta2"

	"github.com/koki/short/util"
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

	Status *apps.ReplicaSetStatus `json:"status,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	TemplateMetadata *PodTemplateMeta `json:"pod_meta,omitempty"`
	PodTemplate      `json:",inline"`
}

type RSSelector struct {
	Shorthand string
	Labels    map[string]string
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
		return util.InvalidValueForTypeErrorf(string(data), s, "couldn't parse JSON as string or dictionary: (%s), (%s)", strErr.Error(), dictErr.Error())
	}

	s.Labels = labels
	return nil
}

func (s RSSelector) MarshalJSON() ([]byte, error) {
	if len(s.Shorthand) > 0 {
		b, err := json.Marshal(s.Shorthand)
		if err != nil {
			return nil, util.InvalidInstanceErrorf(s, "couldn't marshal shorthand string to JSON: %s", err.Error())
		}

		return b, nil
	}

	b, err := json.Marshal(s.Labels)
	if err != nil {
		return nil, util.InvalidInstanceErrorf(s, "couldn't marshal labels dictionary to JSON: %s", err.Error())
	}

	return b, nil
}
