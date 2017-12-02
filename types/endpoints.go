package types

type EndpointWrapper struct {
	Endpoint Endpoint `json:"endpoint,omitempty"`
}

type Endpoint struct {
	Subsets []EndpointSubset `json:"subsets"`
}

type EndpointSubset struct {
	Addresses         []EndpointAddress `json:"addrs,omitempty"`
	NotReadyAddresses []EndpointAddress `json:"unready_addrs,omitempty"`
	Ports             []EndpointPort    `json:"ports,omitempty"`
}

type EndpointAddress struct {
	IP       string           `json:"ip,omitempty"`
	Hostname string           `json:"hostname,omitempty"`
	Nodename *string          `json:"node,omitempty"`
	Target   *ObjectReference `json:"target,omitempty"`
}

type EndpointPort struct {
	Name     string   `json:"name,omitempty"`
	Port     int32    `json:"port,omitempty"`
	Protocol Protocol `json:"protocol,omitempty"`
}
