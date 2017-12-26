package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type APIServiceWrapper struct {
	APIService APIService `json:"api_service"`
}

type APIService struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Service      string `json:"service,omitempty"`
	GroupVersion string `json:"group_version,omitempty"`
	TLSVerify    bool   `json:"tls_verify,omitempty"`

	CABundle []byte `json:"ca_bundle,omitempty"`

	MinGroupPriority int32 `json:"min_group_priority,omitempty"`
	VersionPriority  int32 `json:"version_priority,omitempty"`

	Conditions []APIServiceCondition `json:"conditions,omitempty"`
}

type APIServiceCondition struct {
	Type               APIServiceConditionType   `json:"type,omitempty"`
	Status             APIServiceConditionStatus `json:"status,omitempty"`
	LastTransitionTime metav1.Time               `json:"last_change,omitempty"`
	Reason             string                    `json:"reason,omitempty"`
	Message            string                    `json:"msg,omitempty"`
}

type APIServiceConditionStatus string

const (
	APIServiceConditionTrue    APIServiceConditionStatus = "True"
	APIServiceConditionFalse   APIServiceConditionStatus = "False"
	APIServiceConditionUnknown APIServiceConditionStatus = "Unknown"
)

type APIServiceConditionType string

const (
	APIServiceAvailable APIServiceConditionType = "Available"
)
