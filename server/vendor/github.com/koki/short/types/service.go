package types

import (
	"encoding/json"
	"net"

	"github.com/koki/short/util"
	"github.com/koki/short/util/intbool"
)

type ServiceWrapper struct {
	Service Service `json:"service"`
}

type Service struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// ExternalName services only.
	ExternalName string `json:"cname,omitempty"`

	// ClusterIP services:
	Type ClusterIPServiceType `json:"type,omitempty"`

	Selector    map[string]string `json:"selector,omitempty"`
	ExternalIPs []IPAddr          `json:"external_ips,omitempty"`

	Port     *ServicePort       `json:"port,omitempty"`
	NodePort int32              `json:"node_port,omitempty"`
	Ports    []NamedServicePort `json:"ports,omitempty"`

	ClusterIP ClusterIP `json:"cluster_ip,omitempty"`

	PublishNotReadyAddresses bool                  `json:"unready_endpoints,omitempty"`
	ExternalTrafficPolicy    ExternalTrafficPolicy `json:"route_policy,omitempty"`
	ClientIPAffinity         *intbool.IntOrBool    `json:"stickiness,omitempty"`

	// LoadBalancer services:
	LoadBalancerIP      IPAddr    `json:"lb_ip,omitempty"`
	Allowed             []CIDR    `json:"lb_client_ips,omitempty"`
	HealthCheckNodePort int32     `json:"healthcheck_port,omitempty"`
	Ingress             []Ingress `json:"endpoints,omitempty"`
}

// LoadBalancer helper type.
type LoadBalancer struct {
	IP                  IPAddr
	Allowed             []CIDR
	HealthCheckNodePort int32
	Ingress             []Ingress
}

type ClusterIPServiceType string

const (
	ClusterIPServiceTypeDefault      ClusterIPServiceType = "cluster-ip"
	ClusterIPServiceTypeNodePort     ClusterIPServiceType = "node-port"
	ClusterIPServiceTypeLoadBalancer ClusterIPServiceType = "load-balancer"
)

type IPAddr string
type CIDR string

type ClusterIP string

const (
	ClusterIPNil  ClusterIP = ""
	ClusterIPNone ClusterIP = "None"
)

func ClusterIPAddr(a IPAddr) ClusterIP {
	return ClusterIP(string(a))
}

type Ingress struct {
	IP       net.IP
	Hostname string
}

type ExternalTrafficPolicy string

const (
	ExternalTrafficPolicyNil     ExternalTrafficPolicy = ""
	ExternalTrafficPolicyLocal   ExternalTrafficPolicy = "node-local"
	ExternalTrafficPolicyCluster ExternalTrafficPolicy = "cluster-wide"
)

func (i *Ingress) InitFromString(s string) {
	ip := net.ParseIP(s)
	if ip != nil {
		i.IP = ip
		return
	}

	i.Hostname = s
}

func (i *Ingress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return util.InvalidValueErrorf(string(data), "expected a string for Ingress")
	}

	i.InitFromString(s)
	return nil
}

func (i Ingress) String() string {
	if i.IP != nil {
		return i.IP.String()
	}

	return i.Hostname
}

func (i Ingress) MarshalJSON() ([]byte, error) {
	str := i.String()
	b, err := json.Marshal(str)
	if err != nil {
		return nil, util.InvalidInstanceErrorf(i, "couldn't marshal from string (%s) to JSON: %s", str, err.Error())
	}

	return b, nil
}

func (s *Service) SetLoadBalancer(lb *LoadBalancer) {
	if lb == nil {
		return
	}

	s.LoadBalancerIP = lb.IP
	s.Ingress = lb.Ingress
	s.Allowed = lb.Allowed
	s.HealthCheckNodePort = lb.HealthCheckNodePort
}
