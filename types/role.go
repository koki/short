package types

type RoleWrapper struct {
	Role Role `json:"role"`
}

// Role is a namespaced, logical grouping of PolicyRules that can be referenced as a unit by a RoleBinding.
type Role struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Rules []PolicyRule `json:"rules"`
}

type RoleBindingWrapper struct {
	RoleBinding RoleBinding `json:"role_binding"`
}

// RoleBinding references a role, but does not contain it.  It can reference a Role in the same namespace or a
// ClusterRole in the global namespace. It adds who information via Subjects and namespace information by
// which namespace it exists in.  RoleBindings in a given namespace only have effect in that namespace.
type RoleBinding struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// Subjects holds references to the objects the role applies to.
	Subjects []Subject `json:"subjects"`

	// RoleRef can reference a Role in the current namespace or a ClusterRole in the global namespace.
	// If the RoleRef cannot be resolved, the Authorizer must return an error.
	RoleRef RoleRef `json:"role"`
}

