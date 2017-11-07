package types

type Action struct {
	Command []string   `json:"command,omitempty"`
	Net     *NetAction `json:"net,omitempty"`
}

type NetAction struct {
	Headers []string `json:"headers,omitempty"`
	URL     string   `json:"url,omitempty"`
}
