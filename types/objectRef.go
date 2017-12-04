package types

type ObjectReference struct {
	Kind            string `json:"kind,omitempty"`
	Namespace       string `json:"namespace,omitempty"`
	Name            string `json:"name,omitempty"`
	UID             string `json:"uid,omitempty"`
	Version         string `json:"version,omitempty"`
	ResourceVersion string `json:"resource_version,omitempty"`
	FieldPath       string `json:"field_path,omitempty"`
}
