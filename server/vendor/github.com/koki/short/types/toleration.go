package types

type Selector string

type Toleration struct {
	ExpiryAfter *int64 `json:"expiry_after,omitempty"`
	Selector    `json:",inline"`
}
