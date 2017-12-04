package converters

import (
	"fmt"

	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_v1_Endpoints_to_Koki_Endpoints(kubeEndpoints *v1.Endpoints) (*types.EndpointsWrapper, error) {
	var err error
	kokiWrapper := &types.EndpointsWrapper{}
	kokiEndpoints := &kokiWrapper.Endpoints

	kokiEndpoints.Name = kubeEndpoints.Name
	kokiEndpoints.Namespace = kubeEndpoints.Namespace
	kokiEndpoints.Version = kubeEndpoints.APIVersion
	kokiEndpoints.Cluster = kubeEndpoints.ClusterName
	kokiEndpoints.Labels = kubeEndpoints.Labels
	kokiEndpoints.Annotations = kubeEndpoints.Annotations

	subsets, err := convertSubsets(kubeEndpoints.Subsets)
	if err != nil {
		return nil, err
	}
	kokiEndpoints.Subsets = subsets

	return kokiWrapper, nil
}

func convertSubsets(subsets []v1.EndpointSubset) ([]types.EndpointSubset, error) {
	var kokiSubsets []types.EndpointSubset

	for i := range subsets {
		subset := subsets[i]
		addrs, err := convertEndpointAddresses(subset.Addresses)
		if err != nil {
			return nil, err
		}

		unreadyAddrs, err := convertEndpointAddresses(subset.NotReadyAddresses)
		if err != nil {
			return nil, err
		}

		ports, err := convertEndpointPorts(subset.Ports)
		if err != nil {
			return nil, err
		}

		kokiSubset := types.EndpointSubset{
			Addresses:         addrs,
			NotReadyAddresses: unreadyAddrs,
			Ports:             ports,
		}

		kokiSubsets = append(kokiSubsets, kokiSubset)
	}
	return kokiSubsets, nil
}

func convertEndpointPorts(ports []v1.EndpointPort) ([]string, error) {
	var kokiPorts []string

	for i := range ports {
		port := ports[i]

		protocol := convertProtocol(port.Protocol)

		if protocol == "" {
			protocol = "tcp"
		}

		kokiPort := fmt.Sprintf("%s://%d", protocol, port.Port)

		if port.Name != "" {
			kokiPort = fmt.Sprintf("%s:%s", kokiPort, port.Name)
		}

		kokiPorts = append(kokiPorts, kokiPort)
	}

	return kokiPorts, nil
}

func convertEndpointAddresses(addrs []v1.EndpointAddress) ([]types.EndpointAddress, error) {
	var endpointAddrs []types.EndpointAddress

	for i := range addrs {
		addr := addrs[i]

		target, err := convertTarget(addr.TargetRef)
		if err != nil {
			return nil, err
		}

		endpointAddr := types.EndpointAddress{
			IP:       addr.IP,
			Hostname: addr.Hostname,
			Nodename: addr.NodeName,
			Target:   target,
		}

		endpointAddrs = append(endpointAddrs, endpointAddr)
	}

	return endpointAddrs, nil
}

func convertTarget(target *v1.ObjectReference) (*types.ObjectReference, error) {
	if target == nil {
		return nil, nil
	}

	kokiTarget := &types.ObjectReference{}

	kokiTarget.Kind = target.Kind
	kokiTarget.Namespace = target.Namespace
	kokiTarget.Name = target.Name
	kokiTarget.UID = string(target.UID)
	kokiTarget.Version = target.APIVersion
	kokiTarget.ResourceVersion = target.ResourceVersion
	kokiTarget.FieldPath = target.FieldPath

	return kokiTarget, nil
}
