package types

type Port struct {
	Name     string `json:"name,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	IP       string `json:"ip,omitempty"`
	PortMap  string `json:"port_map,omitempty"`
}
