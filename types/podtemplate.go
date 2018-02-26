package types

import ()

type PodTemplateWrapper struct {
	PodTemplate PodTemplateResource `json:"pod_template"`
}

type PodTemplateResource struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	TemplateMetadata PodTemplateMeta `json:"pod_meta,omitempty"`
	PodTemplate      `json:",inline"`
}
