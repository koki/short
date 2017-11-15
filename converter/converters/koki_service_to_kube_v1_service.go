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

	kubeService.Spec.Type = v1.ServiceTypeClusterIP

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

	hasNodePort, kubePorts, err := revertPorts(kokiService)
	if err != nil {
		return nil, err
	}
	kubeService.Spec.Ports = kubePorts
	if hasNodePort {
		kubeService.Spec.Type = v1.ServiceTypeNodePort
	}

	if kokiService.HasLoadBalancer() {
		kubeService.Spec.Type = v1.ServiceTypeLoadBalancer
		kubeService.Spec.HealthCheckNodePort = kokiService.HealthCheckNodePort
		kubeService.Spec.LoadBalancerIP = string(kokiService.LoadBalancerIP)
		kubeService.Spec.LoadBalancerSourceRanges = revertLoadBalancerSources(kokiService.Allowed)
		kubeService.Status.LoadBalancer.Ingress = revertIngress(kokiService.Ingress)
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

	kubePort.Protocol = kokiPort.Protocol

	return kubePort, nil
}

// Set the Service's Ports and its Type (if there are NodePorts).
func revertPorts(kokiService *types.Service) (hasNodePort bool, kubePorts []v1.ServicePort, err error) {
	if kokiService.Port != nil {
		kubePort, err := revertPort("", kokiService.Port)
		if err != nil {
			return false, nil, err
		}

		kubePorts = []v1.ServicePort{*kubePort}
		hasNodePort = kubePort.NodePort != 0
		return hasNodePort, kubePorts, nil
	}

	kubePorts = make([]v1.ServicePort, 0, len(kokiService.Ports))
	for name, kokiPort := range kokiService.Ports {
		kubePort, err := revertPort(name, &kokiPort)
		if err != nil {
			return false, nil, err
		}

		kubePorts = append(kubePorts, *kubePort)
		if kubePort.NodePort != 0 {
			hasNodePort = true
		}
	}

	return hasNodePort, kubePorts, nil
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
		return "", util.PrettyTypeError(policy, "unrecognized policy")
	}
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
