package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
	"github.com/koki/short/util/intbool"
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

	kubeService.Spec.Type, err = revertServiceType(kokiService.Type)
	if err != nil {
		return nil, err
	}

	kubeService.Spec.Selector = kokiService.Selector
	kubeService.Spec.ClusterIP = string(kokiService.ClusterIP)
	kubeService.Spec.ExternalIPs = revertExternalIPs(kokiService.ExternalIPs)
	kubeService.Spec.SessionAffinity, kubeService.Spec.SessionAffinityConfig = revertSessionAffinity(kokiService.ClientIPAffinity)
	if err != nil {
		return nil, err
	}

	kubeService.Spec.PublishNotReadyAddresses = kokiService.PublishNotReadyAddresses

	kubeService.Spec.ExternalTrafficPolicy, err = revertExternalTrafficPolicy(kokiService.ExternalTrafficPolicy)
	if err != nil {
		return nil, err
	}

	kubePorts, err := revertPorts(kokiService)
	if err != nil {
		return nil, err
	}
	kubeService.Spec.Ports = kubePorts

	if kokiService.Type == types.ClusterIPServiceTypeLoadBalancer {
		kubeService.Spec.HealthCheckNodePort = kokiService.HealthCheckNodePort
		kubeService.Spec.LoadBalancerIP = string(kokiService.LoadBalancerIP)
		kubeService.Spec.LoadBalancerSourceRanges = revertLoadBalancerSources(kokiService.Allowed)
		kubeService.Status.LoadBalancer.Ingress = revertIngress(kokiService.Ingress)
	}

	return kubeService, nil
}

func revertServiceType(kokiType types.ClusterIPServiceType) (v1.ServiceType, error) {
	if len(kokiType) == 0 {
		return v1.ServiceTypeClusterIP, nil
	}
	switch kokiType {
	case types.ClusterIPServiceTypeDefault:
		return v1.ServiceTypeClusterIP, nil
	case types.ClusterIPServiceTypeNodePort:
		return v1.ServiceTypeNodePort, nil
	case types.ClusterIPServiceTypeLoadBalancer:
		return v1.ServiceTypeLoadBalancer, nil
	default:
		return "", util.InvalidInstanceError(kokiType)
	}
}

func revertPort(name string, kokiPort *types.ServicePort, kokiNodePort int32) (*v1.ServicePort, error) {
	kubePort := &v1.ServicePort{}
	kubePort.Port = kokiPort.Expose
	if kokiPort.PodPort != nil {
		kubePort.TargetPort = *kokiPort.PodPort
	}

	if len(name) > 0 {
		kubePort.Name = name
	}

	if kokiNodePort > 0 {
		kubePort.NodePort = kokiNodePort
	}

	kubePort.Protocol = kokiPort.Protocol

	return kubePort, nil
}

func revertPorts(kokiService *types.Service) (kubePorts []v1.ServicePort, err error) {
	if kokiService.Port != nil {
		kubePort, err := revertPort("", kokiService.Port, kokiService.NodePort)
		if err != nil {
			return nil, err
		}

		kubePorts = []v1.ServicePort{*kubePort}
		return kubePorts, nil
	}

	kubePorts = make([]v1.ServicePort, 0, len(kokiService.Ports))
	for _, kokiPort := range kokiService.Ports {
		kubePort, err := revertPort(kokiPort.Name, &kokiPort.Port, kokiPort.NodePort)
		if err != nil {
			return nil, err
		}

		kubePorts = append(kubePorts, *kubePort)
	}

	return kubePorts, nil
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
		return "", util.InvalidInstanceError(policy)
	}
}

func revertIngress(kokiIngress []types.Ingress) []v1.LoadBalancerIngress {
	if len(kokiIngress) == 0 {
		return nil
	}

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

func revertLoadBalancerSources(kokiCIDRs []types.CIDR) []string {
	strs := make([]string, len(kokiCIDRs))
	for i, cidr := range kokiCIDRs {
		strs[i] = string(cidr)
	}

	return strs
}

func revertSessionAffinity(kokiAffinity *intbool.IntOrBool) (v1.ServiceAffinity, *v1.SessionAffinityConfig) {
	if kokiAffinity == nil {
		return "", nil
	}

	switch kokiAffinity.Type {
	case intbool.Int:
		sessionAffinityConfig := &v1.SessionAffinityConfig{
			&v1.ClientIPConfig{util.Int32Ptr(kokiAffinity.IntVal)},
		}
		return v1.ServiceAffinityClientIP, sessionAffinityConfig
	case intbool.Bool:
		if kokiAffinity.BoolVal {
			return v1.ServiceAffinityClientIP, nil
		}
	}

	return "", nil
}

func revertExternalIPs(kokiAddrs []types.IPAddr) []string {
	strs := make([]string, len(kokiAddrs))
	for i, addr := range kokiAddrs {
		strs[i] = string(addr)
	}

	return strs
}
