package types

type Toleration struct {
	GracePeriod int    `json:"grace_period,omitempty"`
	Selector    string `json:"selector,omitempty"`
}
