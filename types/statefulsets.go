package types

type StatefulSetWrapper struct {
	StatefulSet StatefulSet `json:"stateful_set"`
}

type StatefulSet struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Replicas  *int32 `json:"replicas,omitempty"`
	OnDelete  bool   `json:"replace_on_delete,omitempty"`
	Partition *int32 `json:"partition,omitempty"`

	RevisionHistoryLimit *int32                  `json:"max_revs,omitempty"`
	PodManagementPolicy  PodManagementPolicyType `json:"pod_policy,omitempty"`
	Service              string                  `json:"service,omitempty"`
	PVCs                 []PersistentVolumeClaim `json:"pvcs,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	TemplateMetadata *PodTemplateMeta `json:"pod_meta,omitempty"`
	PodTemplate      `json:",inline"`

	// Status
	StatefulSetStatus `json:",inline"`
}

type PodManagementPolicyType string

const (
	OrderedReadyPodManagement PodManagementPolicyType = "ordered"
	ParallelPodManagement                             = "parallel"
)

type StatefulSetStatus struct {
	ObservedGeneration int64  `json:"generation_observed,omitempty"`
	Replicas           int32  `json:"replicas,omitempty"`
	ReadyReplicas      int32  `json:"ready,omitempty"`
	CurrentReplicas    int32  `json:"current,omitempty"`
	UpdatedReplicas    int32  `json:"updated,omitempty"`
	Revision           string `json:"rev,omitempty"`
	UpdateRevision     string `json:"update_rev,omitempty"`
	CollisionCount     *int32 `json:"hash_collisions,omitempty"`
}
