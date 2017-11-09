package converters

import (
	"net"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Kube_v1_Service_to_Koki_Service(kubeService *v1.Service) (*types.ServiceWrapper, error) {
	var err error
	kokiWrapper := &types.ServiceWrapper{}
	kokiService := &kokiWrapper.Service

	kokiService.Name = kubeService.Name
	kokiService.Namespace = kubeService.Namespace
	kokiService.Version = kubeService.APIVersion
	kokiService.Cluster = kubeService.ClusterName
	kokiService.Labels = kubeService.Labels
	kokiService.Annotations = kubeService.Annotations

	if kubeService.Spec.Type == v1.ServiceTypeExternalName {
		kokiService.ExternalName = kubeService.Spec.ExternalName
		return kokiWrapper, nil
	}

	kokiService.PodLabels = kubeService.Spec.Selector
	kokiService.ClusterIP = types.ClusterIP(kubeService.Spec.ClusterIP)
	kokiService.ExternalIPs = convertExternalIPs(kubeService.Spec.ExternalIPs)
	kokiService.ClientIPAffinity = convertSessionAffinity(&kubeService.Spec)
	kokiService.PublishNotReadyAddresses = kubeService.Spec.PublishNotReadyAddresses

	kokiService.ExternalTrafficPolicy, err = convertExternalTrafficPolicy(kubeService.Spec.ExternalTrafficPolicy)
	if err != nil {
		return nil, err
	}

	kokiPort, kokiPorts, err := convertPorts(kubeService.Spec.Ports)
	if err != nil {
		return nil, err
	}
	kokiService.Port = kokiPort
	kokiService.Ports = kokiPorts

	if kubeService.Spec.Type == v1.ServiceTypeLoadBalancer {
		kokiService.LoadBalancer, err = convertLoadBalancer(kubeService)
		if err != nil {
			return nil, err
		}
	}

	return kokiWrapper, nil
}

func convertIngress(kubeIngress []v1.LoadBalancerIngress) ([]types.Ingress, error) {
	if len(kubeIngress) == 0 {
		return nil, nil
	}

	kokiIngress := make([]types.Ingress, len(kubeIngress))
	for index, singleKubeIngress := range kubeIngress {
		if len(singleKubeIngress.IP) > 0 {
			ip := net.ParseIP(singleKubeIngress.IP)
			if ip == nil {
				return nil, util.PrettyTypeError(singleKubeIngress, "invalid IP")
			}

			kokiIngress[index] = types.Ingress{IP: ip}
		} else {
			kokiIngress[index] = types.Ingress{
				Hostname: singleKubeIngress.Hostname,
			}
		}
	}

	return kokiIngress, nil
}

func convertLoadBalancerSources(kubeSources []string) []types.CIDR {
	kokiCIDRs := make([]types.CIDR, len(kubeSources))
	for i, kubeSource := range kubeSources {
		kokiCIDRs[i] = types.CIDR(kubeSource)
	}

	return kokiCIDRs
}

func convertLoadBalancer(kubeService *v1.Service) (*types.LoadBalancer, error) {
	var err error
	kokiLB := &types.LoadBalancer{}
	kokiLB.IP = types.IPAddr(kubeService.Spec.LoadBalancerIP)
	kokiLB.Ingress, err = convertIngress(kubeService.Status.LoadBalancer.Ingress)
	if err != nil {
		return nil, err
	}
	kokiLB.Allowed = convertLoadBalancerSources(kubeService.Spec.LoadBalancerSourceRanges)
	kokiLB.HealthCheckNodePort = kubeService.Spec.HealthCheckNodePort

	return kokiLB, nil
}

func convertPort(kubePort v1.ServicePort) (*types.ServicePort, error) {
	kokiPort := &types.ServicePort{}
	kokiPort.Expose = kubePort.Port
	kokiPort.PodPort = kubePort.TargetPort

	kokiPort.NodePort = kubePort.NodePort
	if kubePort.Protocol != "" {
		switch kubePort.Protocol {
		case v1.ProtocolTCP:
			kokiPort.Protocol = types.ProtocolTCP
		case v1.ProtocolUDP:
			kokiPort.Protocol = types.ProtocolUDP
		default:
			return nil, util.PrettyTypeError(kubePort, "unrecognized protocol")
		}
	}

	return kokiPort, nil
}

func convertPorts(kubePorts []v1.ServicePort) (*types.ServicePort, map[string]types.ServicePort, error) {
	if len(kubePorts) == 1 && len(kubePorts[0].Name) == 0 {
		// Just one port, and it's unnamed
		kokiPort, err := convertPort(kubePorts[0])
		return kokiPort, nil, err
	}

	kokiPorts := make(map[string]types.ServicePort, len(kubePorts))
	for _, kubePort := range kubePorts {
		kokiPort, err := convertPort(kubePort)
		if err != nil {
			return nil, nil, err
		}

		kokiPorts[kubePort.Name] = *kokiPort
	}

	return nil, kokiPorts, nil
}

func convertExternalTrafficPolicy(kubePolicy v1.ServiceExternalTrafficPolicyType) (types.ExternalTrafficPolicy, error) {
	switch kubePolicy {
	case "":
		return types.ExternalTrafficPolicyNil, nil
	case v1.ServiceExternalTrafficPolicyTypeLocal:
		return types.ExternalTrafficPolicyLocal, nil
	case v1.ServiceExternalTrafficPolicyTypeCluster:
		return types.ExternalTrafficPolicyCluster, nil
	default:
		return "", util.PrettyTypeError(kubePolicy, "unrecognized value")
	}
}

// Returns koki ClientIPAffinitySeconds.
func convertSessionAffinity(kubeSpec *v1.ServiceSpec) *intstr.IntOrString {
	if kubeSpec.SessionAffinity == v1.ServiceAffinityClientIP {
		if kubeSpec.SessionAffinityConfig != nil && kubeSpec.SessionAffinityConfig.ClientIP != nil && kubeSpec.SessionAffinityConfig.ClientIP.TimeoutSeconds != nil {
			return types.ClientIPAffinitySeconds(
				int(*kubeSpec.SessionAffinityConfig.ClientIP.TimeoutSeconds))
		}

		return types.ClientIPAffinityDefault()
	}

	return nil
}

func convertExternalIPs(kubeIPs []string) []types.IPAddr {
	kokiIPs := make([]types.IPAddr, len(kubeIPs))
	for i, kubeIP := range kubeIPs {
		kokiIPs[i] = types.IPAddr(kubeIP)
	}

	return kokiIPs
}
