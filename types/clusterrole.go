package types

import (
	"fmt"
	"strings"

	"github.com/koki/json"
	serrors "github.com/koki/structurederrors"
)

type ClusterRoleWrapper struct {
	ClusterRole ClusterRole `json:"cluster_role"`
}

// ClusterRole is a cluster level, logical grouping of PolicyRules that can be referenced as a unit by a RoleBinding or ClusterRoleBinding.
type ClusterRole struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Rules []PolicyRule `json:"rules"`

	// AggregationRule.ClusterRoleSelectors :: []metav1.LabelSelector
	AggregationRule []string `json:"aggregation,omitempty"`
}

type PolicyRule struct {
	Verbs []string `json:"verbs"`

	APIGroups     []string `json:"groups,omitempty"`
	Resources     []string `json:"resources,omitempty"`
	ResourceNames []string `json:"resource_names,omitempty"`

	NonResourceURLs []string `json:"non_resource_urls,omitempty"`
}

type ClusterRoleBindingWrapper struct {
	ClusterRoleBinding ClusterRoleBinding `json:"cluster_role_binding"`
}

// ClusterRoleBinding references a ClusterRole, but not contain it.  It can reference a ClusterRole in the global namespace,
// and adds who information via Subject.
type ClusterRoleBinding struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// Subjects holds references to the objects the role applies to.
	Subjects []Subject `json:"subjects"`

	// RoleRef can only reference a ClusterRole in the global namespace.
	// If the RoleRef cannot be resolved, the Authorizer must return an error.
	RoleRef RoleRef `json:"role"`
}

type Subject struct {
	// Kind of object being referenced. Values defined by this API group are "User", "Group", and "ServiceAccount".
	Kind string `json:"kind"`
	// APIGroup holds the API group of the referenced subject.
	// Defaults to "" for ServiceAccount subjects.
	// Defaults to "rbac.authorization.k8s.io" for User and Group subjects.
	APIGroup string `json:"apiGroup,omitempty"`
	Name     Name   `json:"name"`
	// Namespace of the referenced object.  If the object kind is non-namespace, such as "User" or "Group", and this value is not empty
	// the Authorizer should report an error.
	Namespace string `json:"namespace,omitempty"`
}

// RoleRef contains information that points to the role being used
type RoleRef struct {
	APIGroup string `json:"apiGroup"`
	Kind     string `json:"kind"`
	Name     Name   `json:"name"`
}

func (r RoleRef) GroupKind() string {
	if len(r.APIGroup) > 0 {
		return fmt.Sprintf("%s.%s", r.APIGroup, r.Kind)
	}

	return r.Kind
}

func (r RoleRef) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%s:%s", r.GroupKind(), EscapeName(r.Name))
	return json.Marshal(str)
}

func (r *RoleRef) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	segments := SplitAtUnescapedColons(str)
	if len(segments) != 2 {
		return serrors.InvalidValueForTypeErrorf(str, r, "expected 'group.kind:name'")
	}

	r.Name = UnescapeName(segments[1])

	splitAt := strings.LastIndex(segments[0], ".")
	if splitAt >= 0 {
		r.APIGroup = segments[0][:splitAt]
		r.Kind = segments[0][splitAt+1:]
	} else {
		return serrors.InvalidValueForTypeErrorf(str, r, "expected 'group.kind:name'")
	}

	return nil
}

func (r Subject) GroupKind() string {
	if len(r.APIGroup) > 0 {
		return fmt.Sprintf("%s.%s", r.APIGroup, r.Kind)
	}

	return r.Kind
}

func (r Subject) MarshalJSON() ([]byte, error) {
	var str string
	if len(r.Namespace) > 0 {
		str = fmt.Sprintf("%s:%s:%s", r.GroupKind(), r.Namespace, EscapeName(r.Name))
	} else {
		str = fmt.Sprintf("%s:%s", r.GroupKind(), EscapeName(r.Name))
	}
	return json.Marshal(str)
}

func (r *Subject) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	segments := SplitAtUnescapedColons(str)
	if len(segments) != 3 && len(segments) != 2 {
		return serrors.InvalidValueForTypeErrorf(str, r, "expected '[group.kind|kind]:name' OR '[group.kind|kind]:namespace:name'")
	}

	splitAt := strings.LastIndex(segments[0], ".")
	if splitAt >= 0 {
		r.APIGroup = segments[0][:splitAt]
		r.Kind = segments[0][splitAt+1:]
	} else {
		r.Kind = segments[0]
	}

	if len(segments) == 2 {
		r.Name = UnescapeName(segments[1])
	} else { // len is 3
		r.Namespace = segments[1]
		r.Name = UnescapeName(segments[2])
	}

	return nil
}
