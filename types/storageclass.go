package types

type StorageClassWrapper struct {
	StorageClass `json:"storage_class,omitempty"`
}

type StorageClass struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Provisioner          string                         `json:"provisioner,omitempty"`
	Parameters           map[string]string              `json:"params, omitempty"`
	Reclaim              *PersistentVolumeReclaimPolicy `json:"reclaim,omitempty"`
	MountOptions         []string                       `json:"mount_opts,omitempty"`
	AllowVolumeExpansion *bool                          `json:"allow_expansion,omitempty"`
	VolumeBindingMode    *VolumeBindingMode             `json:"binding_mode,omitempty"`
}

type VolumeBindingMode string

const (
	VolumeBindingImmediate            VolumeBindingMode = "immediate"
	VolumeBindingWaitForFirstConsumer VolumeBindingMode = "wait-for-first-consumer"
)
