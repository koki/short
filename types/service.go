package types

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/golang/glog"

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

	Selector    map[string]string      `json:"selector,omitempty"`
	ExternalIPs []IPAddr               `json:"external_ips,omitempty"`
	Port        *ServicePort           `json:"port,omitempty"`
	Ports       map[string]ServicePort `json:"ports,omitempty"`
	ClusterIP   ClusterIP              `json:"cluster_ip,omitempty"`

	PublishNotReadyAddresses bool                  `json:"unready_endpoints,omitempty"`
	ExternalTrafficPolicy    ExternalTrafficPolicy `json:"route_policy,omitempty"`
	ClientIPAffinity         *intbool.IntOrBool    `json:"stickiness,omitempty"`

	// LoadBalancer services:
	LoadBalancer        bool      `json:"lb,omitempty"`
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

type ServicePort struct {
	Expose int32

	// PodPort is a port or the name of a containerPort.
	PodPort intstr.IntOrString

	// NodePort is optional. 0 is empty.
	NodePort int32

	// Protocol is optional. "" is empty.
	Protocol v1.Protocol
}

type ExternalTrafficPolicy string

const (
	ExternalTrafficPolicyNil     ExternalTrafficPolicy = ""
	ExternalTrafficPolicyLocal   ExternalTrafficPolicy = "node-local"
	ExternalTrafficPolicyCluster ExternalTrafficPolicy = "cluster-wide"
)

func (s *Service) HasLoadBalancer() bool {
	if s.LoadBalancer {
		return true
	}

	if len(s.LoadBalancerIP) > 0 {
		return true
	}

	if len(s.Allowed) > 0 {
		return true
	}

	if s.HealthCheckNodePort > 0 {
		return true
	}

	if len(s.Ingress) > 0 {
		return true
	}

	return false
}

func (p *ServicePort) InitFromString(str string) error {
	matches := protocolPortRegexp.FindStringSubmatch(str)
	if len(matches) > 0 {
		p.Protocol = v1.Protocol(matches[1])
		str = matches[2]
	} else {
		p.Protocol = v1.ProtocolTCP
	}

	segments := strings.Split(str, ":")
	l := len(segments)
	if l < 2 {
		glog.Error("Sections for Expose & Pod port are both required.")
		return fmt.Errorf("too few sections in (%s)", str)
	}
	if l > 4 {
		glog.Error("Too many sections for ServicePort")
		return fmt.Errorf("too many sections in (%s)", str)
	}

	expose, err := strconv.ParseInt(segments[0], 10, 32)
	if err != nil {
		return util.PrettyTypeError(p, str)
	}
	p.Expose = int32(expose)

	p.PodPort = intstr.Parse(segments[1])

	if l == 3 {
		nodePort, err := strconv.ParseInt(segments[2], 10, 32)
		if err != nil {
			return util.PrettyTypeError(p, str)
		}
		p.NodePort = int32(nodePort)
	}

	return nil
}

func appendColonIntSegment(str string, i int32) string {
	if len(str) == 0 {
		return fmt.Sprintf("%d", i)
	}

	return fmt.Sprintf("%s:%d", str, i)
}

func (p *ServicePort) String() string {
	str := fmt.Sprintf("%d:%s", p.Expose, p.PodPort.String())
	if p.NodePort > 0 {
		str = appendColonIntSegment(str, p.NodePort)
	}

	if len(p.Protocol) == 0 || p.Protocol == v1.ProtocolTCP {
		// No need to specify protocol
		return str
	}

	return fmt.Sprintf("%s://%s", p.Protocol, str)
}

func (p *ServicePort) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		glog.Error("Expected a string for ServicePort")
		return err
	}

	p.InitFromString(s)
	return nil
}

func (p ServicePort) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

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
		glog.Error("Expected a string for Ingress")
		return err
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
	return json.Marshal(i.String())
}

func (l LoadBalancer) IsZero() bool {
	return reflect.DeepEqual(l, LoadBalancer{})
}

func (s *Service) SetLoadBalancer(lb *LoadBalancer) {
	if lb == nil {
		return
	}

	if lb.IsZero() {
		s.LoadBalancer = true
		return
	}

	s.LoadBalancerIP = lb.IP
	s.Ingress = lb.Ingress
	s.Allowed = lb.Allowed
	s.HealthCheckNodePort = lb.HealthCheckNodePort
}
