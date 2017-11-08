package types

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/golang/glog"
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
	ExternalName string `json:"externalName,omitempty"`

	// ClusterIP services:

	PodLabels   map[string]string      `json:"podLabels,omitempty"`
	ExternalIPs []IPAddr               `json:"externalIPs,omitempty"`
	Port        *ServicePort           `json:"port,omitempty"`
	Ports       map[string]ServicePort `json:"ports,omitempty"`
	ClusterIP   ClusterIP              `json:"clusterIP,omitempty"`

	PublishNotReadyAddresses bool                  `json:"publishNotReadyAddresses,omitempty"`
	ExternalTrafficPolicy    ExternalTrafficPolicy `json:"externalTrafficPolicy,omitempty"`
	ClientIPAffinity         *intstr.IntOrString   `json:"clientIPAffinitySeconds,omitempty"`

	// LoadBalancer services:

	LoadBalancer *LoadBalancer `json:"loadBalancer,omitempty"`
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

type Protocol string

const (
	ProtocolUDP Protocol = "UDP"
	ProtocolTCP Protocol = "TCP"
)

type ServicePort struct {
	Expose int32

	// PodPort is a port or the name of a containerPort.
	PodPort intstr.IntOrString

	// NodePort is optional. 0 is empty.
	NodePort int32

	// Protocol is optional. "" is empty.
	Protocol Protocol
}

type ExternalTrafficPolicy string

const (
	ExternalTrafficPolicyNil     ExternalTrafficPolicy = ""
	ExternalTrafficPolicyLocal   ExternalTrafficPolicy = "Local"
	ExternalTrafficPolicyCluster ExternalTrafficPolicy = "Cluster"
)

func ClientIPAffinitySeconds(s int) *intstr.IntOrString {
	x := intstr.FromInt(s)
	return &x
}

func ClientIPAffinityDefault() *intstr.IntOrString {
	x := intstr.FromString("Default")
	return &x
}

type LoadBalancer struct {
	IP                  IPAddr `json:"ip,omitempty"`
	Allowed             []CIDR `json:"allowed,omitempty"`
	HealthCheckNodePort int32  `json:"healthCheckNodePort,omitempty"`

	// From Service.Status:

	Ingress []Ingress `json:"ingress,omitempty"`
}

func (p *ServicePort) InitFromString(s string) error {
	segments := strings.Split(s, ":")
	l := len(segments)
	if l < 2 {
		glog.Error("Sections for Expose & Pod port are both required.")
		return fmt.Errorf("too few sections in (%s)", s)
	}

	if l > 4 {
		glog.Error("Too many sections for ServicePort")
		return fmt.Errorf("too many sections in (%s)", s)
	}

	expose, err := strconv.ParseInt(segments[0], 10, 32)
	if err != nil {
		glog.Error("Expose should be a port number.")
		return err
	}

	p.Expose = int32(expose)

	p.PodPort = intstr.Parse(segments[1])

	if l >= 3 {
		nodePort, err := strconv.ParseInt(segments[2], 10, 32)
		if err != nil {
			glog.Error("NodePort should be a port number.")
			return err
		}
		p.NodePort = int32(nodePort)
	}

	if l == 4 {
		switch segments[3] {
		case "TCP":
			p.Protocol = ProtocolTCP
		case "UDP":
			p.Protocol = ProtocolUDP
		default:
			return fmt.Errorf("unrecognized Protocol (%s)", segments[3])
		}
	}

	return nil
}

func (p *ServicePort) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		glog.Error("Expected a string for ServicePort")
		return err
	}

	return p.InitFromString(s)
}

func (p ServicePort) String() string {
	s := fmt.Sprintf("%d:%s", p.Expose, p.PodPort.String())

	if p.NodePort == 0 {
		return s
	}

	s = fmt.Sprintf("%s:%d", s, p.NodePort)

	if p.Protocol == "" {
		return s
	}

	return fmt.Sprintf("%s:%s", s, p.Protocol)
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
