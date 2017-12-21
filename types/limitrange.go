package types

import (
	"k8s.io/api/core/v1"
)

type LimitRangeWrapper struct {
	LimitRange `json:"limit_range"`
}

type LimitRange struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// Spec::LimitRangeSpec
	Limits []LimitRangeItem `json:"limits"`
}

type LimitRangeItem struct {
	// Type of resource that this limit applies to.
	Type LimitType `json:"kind,omitempty"`
	// Max usage constraints on this kind by resource name.
	Max v1.ResourceList `json:"max,omitempty"`
	// Min usage constraints on this kind by resource name.
	Min v1.ResourceList `json:"min,omitempty"`
	// Default resource requirement limit value by resource name
	//   (if resource limit is omitted)
	Default v1.ResourceList `json:"default_max,omitempty"`
	// default resource requirement request value by resource name
	// (if resource request is omitted)
	DefaultRequest v1.ResourceList `json:"default_min,omitempty"`
	// MaxLimitRequestRatio represents the max burst for the named resource.
	MaxLimitRequestRatio v1.ResourceList `json:"max_burst_ratio,omitempty"`
}

type LimitType string

const (
	LimitTypePod                   LimitType = "pod"
	LimitTypeContainer             LimitType = "container"
	LimitTypePersistentVolumeClaim LimitType = "pvc"
)
