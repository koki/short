package types

type EndpointsWrapper struct {
	Endpoints Endpoints `json:"endpoints,omitempty"`
}

type Endpoints struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Subsets []EndpointSubset `json:"subsets,omitempty"`
}

type EndpointSubset struct {
	Addresses         []EndpointAddress `json:"addrs,omitempty"`
	NotReadyAddresses []EndpointAddress `json:"unready_addrs,omitempty"`
	Ports             []string          `json:"ports,omitempty"`
}

type EndpointAddress struct {
	IP       string           `json:"ip,omitempty"`
	Hostname string           `json:"hostname,omitempty"`
	Nodename *string          `json:"node,omitempty"`
	Target   *ObjectReference `json:"target,omitempty"`
}
