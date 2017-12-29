package types

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
