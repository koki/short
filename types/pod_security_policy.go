package types

type PodSecurityPolicyWrapper struct {
	PodSecurityPolicy PodSecurityPolicy `json:"pod_security_policy"`
}

type PodSecurityPolicy struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Privileged          bool     `json:"privileged,omitempty"`
	AllowCapabilities   []string `json:"cap_allow,omitempty"`
	DenyCapabilities    []string `json:"cap_deny,omitempty"`
	DefaultCapabilities []string `json:"cap_default,omitempty"`

	VolumePlugins  []string        `json:"vol_plugins,omitempty"`
	HostMode       []HostMode      `json:"host_mode,omitempty"`
	HostPortRanges []HostPortRange `json:"host_port_ranges,omitempty"`

	SELinux     SELinuxPolicy `json:"selinux_policy,omitempty"`
	UIDPolicy   UIDPolicy     `json:"uid_policy,omitempty"`
	GIDPolicy   GIDPolicy     `json:"gid_policy,omitempty"`
	FSGIDPolicy GIDPolicy     `json:"fsgid_policy,omitempty"`

	ReadOnlyRootFS         bool  `json:"rootfs_ro,omitempty"`
	AllowEscalation        *bool `json:"allow_escalation,omitempty"`
	AllowEscalationDefault *bool `json:"allow_escalation_default,omitempty"`

	AllowedHostPaths   []string `json:"host_paths_allowed,omitempty"`
	AllowedFlexVolumes []string `json:"flex_volumes_allowed,omitempty"`
}

type SELinuxPolicy struct {
	Policy  SELinuxPolicyType `json:"policy,omitempty"`
	SELinux `json:",inline"`
}

type SELinuxPolicyType string

const (
	SELinuxPolicyAny  SELinuxPolicyType = "*"
	SELinuxPolicyMust SELinuxPolicyType = "must_be"
)

type UIDPolicy struct {
	Policy UIDPolicyType `json:"policy,omitempty"`
	Ranges []IDRange     `json:"ranges,omitempty"`
}

type UIDPolicyType string

const (
	UIDPolicyAny     UIDPolicyType = "*"
	UIDPolicyMust    UIDPolicyType = "must_be"
	UIDPolicyNonRoot UIDPolicyType = "non_root"
)

type GIDPolicy struct {
	Policy GIDPolicyType `json:"policy,omitempty"`
	Ranges []IDRange     `json:"ranges,omitempty"`
}

type GIDPolicyType string

const (
	GIDPolicyAny  GIDPolicyType = "*"
	GIDPolicyMust GIDPolicyType = "must_be"
)

type IDRange struct {
	Min int64 `json:"min,omitempty"`
	Max int64 `json:"max,omitempty"`
}

type HostPortRange struct {
	Min int32 `json:"min,omitempty"`
	Max int32 `json:"max,omitempty"`
}
