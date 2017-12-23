package types

type PodPresetWrapper struct {
	PodPreset PodPreset `json:"pod_preset"`
}

type PodPreset struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	Env          []Env             `json:"env,omitempty"`
	Volumes      map[string]Volume `json:"volumes,omitempty"`
	VolumeMounts []VolumeMount     `json:"mounts,omitempty"`
}
