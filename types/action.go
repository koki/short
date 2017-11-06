package types

import (
	"net/url"
)

type Action struct {
	Command []string  `json:"command,omitempty"`
	Net     NetAction `json:"net,omitempty"`
}

type NetAction struct {
	Headers map[string]string `json:"headers,omitempty"`
	URL     url.URL           `json:"net,omitempty"`
	Scheme  HttpScheme        `json:"scheme,omitempty"`
}

type HttpScheme string

const (
	HttpGetScheme    HttpScheme = "GET"
	HttpPutScheme    HttpScheme = "PUT"
	HttpPostScheme   HttpScheme = "POST"
	HttpHeadScheme   HttpScheme = "HEAD"
	HttpDeleteScheme HttpScheme = "DELETE"
)
