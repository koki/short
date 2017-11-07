package types

type Lifecycle struct {
	OnStart *Action `json:"on_start,omitempty"`
	PreStop *Action `json:"pre_stop,omitempty"`
}
