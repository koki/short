package converters

import (
	"fmt"
	"reflect"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_Service_To_Kube_v1_Service(service *types.ServiceWrapper) (*v1.Service, error) {
	var err error
	kubeService := &v1.Service{}
	kokiService := &service.Service

	kubeService.Name = kokiService.Name
	kubeService.Namespace = kokiService.Namespace
	kubeService.APIVersion = kokiService.Version
	kubeService.Kind = "Service"
	kubeService.ClusterName = kokiService.Cluster
	kubeService.Labels = kokiService.Labels
	kubeService.Annotations = kokiService.Annotations

	if len(kokiService.ExternalName) > 0 {
		kubeService.Spec.Type = v1.ServiceTypeExternalName
		kubeService.Spec.ExternalName = kokiService.ExternalName
		return kubeService, nil
	}

	kubeService.Spec.Type = v1.ServiceTypeClusterIP

	kubeService.Spec.Selector = kokiService.PodLabels
	kubeService.Spec.ClusterIP = string(kokiService.ClusterIP)
	kubeService.Spec.ExternalIPs = revertExternalIPs(kokiService.ExternalIPs)
	err = revertSessionAffinityInto(kokiService.ClientIPAffinity, &kubeService.Spec)
	if err != nil {
		return nil, err
	}

	kubeService.Spec.PublishNotReadyAddresses = kokiService.PublishNotReadyAddresses

	kubeService.Spec.ExternalTrafficPolicy, err = revertExternalTrafficPolicy(kokiService.ExternalTrafficPolicy)
	if err != nil {
		return nil, err
	}

	err = revertPortsInto(kokiService, &kubeService.Spec)
	if err != nil {
		return nil, err
	}

	if kokiService.LoadBalancer != nil {
		kubeService.Spec.Type = v1.ServiceTypeLoadBalancer
		revertLoadBalancerInto(kokiService.LoadBalancer, kubeService)
	}

	return kubeService, nil
}

func revertPort(name string, kokiPort *types.ServicePort) (*v1.ServicePort, error) {
	kubePort := &v1.ServicePort{}
	kubePort.Port = kokiPort.Expose
	kubePort.TargetPort = kokiPort.PodPort

	if len(name) > 0 {
		kubePort.Name = name
	}

	if kokiPort.NodePort != 0 {
		kubePort.NodePort = kokiPort.NodePort
	}

	if kokiPort.Protocol != "" {
		switch kokiPort.Protocol {
		case types.ProtocolTCP:
			kubePort.Protocol = v1.ProtocolTCP
		case types.ProtocolUDP:
			kubePort.Protocol = v1.ProtocolUDP
		default:
			return nil, fmt.Errorf(
				"unrecognized ServicePort Protocol (%#v)",
				kokiPort.Protocol)
		}
	}

	return kubePort, nil
}

func revertPortsInto(kokiService *types.Service, into *v1.ServiceSpec) error {
	if kokiService.Port != nil {
		kubePort, err := revertPort("", kokiService.Port)
		if err != nil {
			return err
		}

		into.Ports = []v1.ServicePort{*kubePort}
		if kubePort.NodePort != 0 {
			into.Type = v1.ServiceTypeNodePort
		}
	} else {
		kubePorts := make([]v1.ServicePort, 0, len(kokiService.Ports))
		for name, kokiPort := range kokiService.Ports {
			kubePort, err := revertPort(name, &kokiPort)
			if err != nil {
				return err
			}

			kubePorts = append(kubePorts, *kubePort)
			if kubePort.NodePort != 0 {
				into.Type = v1.ServiceTypeNodePort
			}
		}

		into.Ports = kubePorts
	}

	return nil
}

func revertExternalTrafficPolicy(policy types.ExternalTrafficPolicy) (v1.ServiceExternalTrafficPolicyType, error) {
	switch policy {
	case types.ExternalTrafficPolicyNil:
		return "", nil
	case types.ExternalTrafficPolicyLocal:
		return v1.ServiceExternalTrafficPolicyTypeLocal, nil
	case types.ExternalTrafficPolicyCluster:
		return v1.ServiceExternalTrafficPolicyTypeCluster, nil
	default:
		return "", fmt.Errorf("unrecognized ExternalTrafficPolicy (%s)", policy)
	}
}

func revertLoadBalancerInto(kokiLB *types.LoadBalancer, into *v1.Service) {
	into.Spec.HealthCheckNodePort = kokiLB.HealthCheckNodePort

	into.Spec.LoadBalancerIP = string(kokiLB.IP)
	into.Spec.LoadBalancerSourceRanges = revertLoadBalancerSources(kokiLB.Allowed)

	into.Status.LoadBalancer.Ingress = revertIngress(kokiLB.Ingress)
}

func revertIngress(kokiIngress []types.Ingress) []v1.LoadBalancerIngress {
	kubeIngress := make([]v1.LoadBalancerIngress, len(kokiIngress))
	for index, singleKokiIngress := range kokiIngress {
		if singleKokiIngress.IP != nil {
			kubeIngress[index] = v1.LoadBalancerIngress{IP: singleKokiIngress.IP.String()}
		} else {
			kubeIngress[index] = v1.LoadBalancerIngress{Hostname: singleKokiIngress.Hostname}
		}
	}

	return kubeIngress
}

func revertLoadBalancerSources(kokiCidrs []types.CIDR) []string {
	strs := make([]string, len(kokiCidrs))
	for i, cidr := range kokiCidrs {
		strs[i] = string(cidr)
	}

	return strs
}

func revertSessionAffinityInto(kokiAffinity *intstr.IntOrString, into *v1.ServiceSpec) error {
	if kokiAffinity == nil {
		return nil
	}

	if reflect.DeepEqual(kokiAffinity, types.ClientIPAffinityDefault()) {
		into.SessionAffinity = "ClientIP"
		return nil
	}

	switch kokiAffinity.Type {
	case intstr.Int:
		into.SessionAffinity = "ClientIP"
		into.SessionAffinityConfig = &v1.SessionAffinityConfig{
			&v1.ClientIPConfig{util.Int32Ptr(kokiAffinity.IntVal)},
		}
		return nil
	default:
		return fmt.Errorf("unrecognized ClientIPAffinity (%#v)", kokiAffinity)
	}
}

func revertExternalIPs(kokiAddrs []types.IPAddr) []string {
	strs := make([]string, len(kokiAddrs))
	for i, addr := range kokiAddrs {
		strs[i] = string(addr)
	}

	return strs
}
