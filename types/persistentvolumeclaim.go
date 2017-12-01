package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PersistentVolumeClaimWrapper struct {
	PersistentVolumeClaim `json:"pvc,omitempty"`
}

type PersistentVolumeClaim struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	StorageClass *string                      `json:"storage_class,omitempty"`
	Volume       string                       `json:"volume,omitempty"`
	AccessModes  []PersistentVolumeAccessMode `json:"access_modes,omitempty"`
	Storage      string                       `json:"storage,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	PersistentVolumeClaimStatus `json:",inline"`
}

type PersistentVolumeAccessMode string

const (
	ReadWriteOnce PersistentVolumeAccessMode = "rw_once"
	ReadOnlyMany  PersistentVolumeAccessMode = "ro_many"
	ReadWriteMany PersistentVolumeAccessMode = "rw_many"
)

type PersistentVolumeClaimStatus struct {
	Phase       PersistentVolumeClaimPhase       `json:"phase,omitempty"`
	AccessModes []PersistentVolumeAccessMode     `json:"access_modes,omitempty"`
	Storage     string                           `json:"storage,omitempty"`
	Conditions  []PersistentVolumeClaimCondition `json:"condition,omitempty"`
}

type PersistentVolumeClaimCondition struct {
	Type               PersistentVolumeClaimConditionType `json:"type,omitempty"`
	Status             ConditionStatus                    `json:"status"`
	LastProbeTime      metav1.Time                        `json:"timestamp,omitempty"`
	LastTransitionTime metav1.Time                        `json:"last_change,omitempty"`
	Reason             string                             `json:"reason,omitempty"`
	Message            string                             `json:"message,omitempty"`
}

type PersistentVolumeClaimPhase string

const (
	ClaimPending PersistentVolumeClaimPhase = "pending"
	ClaimBound   PersistentVolumeClaimPhase = "bound"
	ClaimLost    PersistentVolumeClaimPhase = "lost"
)

type PersistentVolumeClaimConditionType string

const (
	PersistentVolumeClaimResizing PersistentVolumeClaimConditionType = "resizing"
)
