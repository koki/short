package types

import (
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CRDWrapper struct {
	CRD CustomResourceDefinition `json:"crd"`
}

type CustomResourceDefinition struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// Spec::CRDSpec
	//   Group::string, Version::string, Names::CRDNames
	CRDMeta    CRDMeta                 `json:"meta,omitempty"`
	Scope      CRDResourceScope        `json:"scope,omitempty"`
	Validation *apiext.JSONSchemaProps `json:"validation"`

	// Status
	Conditions []CRDCondition `json:"conditions,omitempty"`
	Accepted   CRDName        `json:"accepted,omitempty"`
}

type CRDMeta struct {
	Group   string `json:"group,omitempty"`
	Version string `json:"version,omitempty"`

	CRDName `json:",inline"`
}

type CRDName struct {
	// Plural is the plural name of the resource to serve.  It must match the name of the CustomResourceDefinition-registration
	// too: plural.group and it must be all lowercase.
	Plural string `json:"plural,omitempty"`
	// Singular is the singular name of the resource.  It must be all lowercase  Defaults to lowercased <kind>
	Singular string `json:"singular,omitempty"`
	// ShortNames are short names for the resource.  It must be all lowercase.
	ShortNames []string `json:"short,omitempty"`
	// Kind is the serialized kind of the resource.  It is normally CamelCase and singular.
	Kind string `json:"kind,omitempty"`
	// ListKind is the serialized kind of the list for this resource.  Defaults to <kind>List.
	ListKind string `json:"list,omitempty"`
}

// ResourceScope is an enum defining the different scopes available to a custom resource
type CRDResourceScope string

const (
	CRDClusterScoped   CRDResourceScope = "cluster"
	CRDNamespaceScoped CRDResourceScope = "ns"
)

// CRDConditionType is a valid value for CRDCondition.Type
type CRDConditionType string

const (
	CRDEstablished   CRDConditionType = "established"
	CRDNamesAccepted CRDConditionType = "names-accepted"
	CRDTerminating   CRDConditionType = "terminating"
)

// CRDCondition contains details for the current condition of this pod.
type CRDCondition struct {
	Type   CRDConditionType `json:"type"`
	Status ConditionStatus  `json:"status"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"last_change,omitempty"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	Reason  string `json:"reason"`
	Message string `json:"msg"`
}
