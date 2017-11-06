package types

type Probe struct {
	Action
	Delay           int `json:"delay,omitempty"`
	Interval        int `json:"interval,omitempty"`
	MinCountSuccess int `json:"min_count_success,omitempty"`
	MinCountFailure int `json:"min_count_fail,omitempty"`
	Timeout         int `json:"timeout,omitempty"`
}
