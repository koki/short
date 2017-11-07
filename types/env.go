package types

type EnvStr string

type Env struct {
	EnvStr   `json:"val,inline"` /*embedded key or key=value or prefix*/
	From     string              `json:"from,omitempty"`
	Required *bool               `json:"required,omitempty"`
}
