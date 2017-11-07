package types

type Action struct {
	Command []string   `json:"command,omitempty"`
	Net     *NetAction `json:"net,omitempty"`
}

type NetAction struct {
	Headers []string   `json:"headers,omitempty"`
	URL     string     `json:"url,omitempty"`
	Method  HTTPMethod `json:"method,omitempty"`
}

type HTTPMethod string

const (
	HTTPGetMethod    HTTPMethod = "GET"
	HTTPPutMethod    HTTPMethod = "PUT"
	HTTPPostMethod   HTTPMethod = "POST"
	HTTPHeadMethod   HTTPMethod = "HEAD"
	HTTPDeleteMethod HTTPMethod = "DELETE"
)
