package types

import (
	"k8s.io/apimachinery/pkg/util/intstr"
)

type IngressWrapper struct {
	Ingress Ingress `json:"ingress"`
}

type Ingress struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// Backend::*IngressBackend
	ServiceName string              `json:"backend,omitempty"`
	ServicePort *intstr.IntOrString `json:"backend_port,omitempty"`

	TLS   []IngressTLS  `json:"tls,omitempty"`
	Rules []IngressRule `json:"rules,omitempty"`

	// Status::IngressStatus LoadBalancer::LoadBalancerStatus
	LoadBalancerIngress []LoadBalancerIngress `json:"endpoints,omitempty"`
}

type IngressTLS struct {
	Hosts      []string `json:"hosts,omitempty"`
	SecretName string   `json:"secret,omitempty"`
}

type IngressRule struct {
	Host string `json:"host,omitempty"`

	// inline::IngressRuleValue HTTP::*HTTPIngressRuleValue
	Paths []HTTPIngressPath `json:"paths"`
}

type HTTPIngressPath struct {
	Path string `json:"path,omitempty"`

	// Backend::IngressBackend
	ServiceName string             `json:"service"`
	ServicePort intstr.IntOrString `json:"port"`
}
