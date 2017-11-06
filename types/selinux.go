package types

type SELinux struct {
	Level string `json:"level,omitempty"`
	Role  string `json:"role,omitempty"`
	Type  string `json:"type,omitempty"`
	User  string `json:"user,omitempty"`
}
