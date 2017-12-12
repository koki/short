package converters

import (
	"strconv"
	"strings"

	"k8s.io/api/core/v1"
	kubeTypes "k8s.io/apimachinery/pkg/types"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_Endpoints_to_Kube_v1_Endpoints(endpoints *types.EndpointsWrapper) (*v1.Endpoints, error) {
	kubeEndpoints := &v1.Endpoints{}
	kokiEndpoints := &endpoints.Endpoints

	kubeEndpoints.Name = kokiEndpoints.Name
	kubeEndpoints.Namespace = kokiEndpoints.Namespace
	kubeEndpoints.APIVersion = "v1"
	kubeEndpoints.Kind = "Endpoints"
	kubeEndpoints.ClusterName = kokiEndpoints.Cluster
	kubeEndpoints.Labels = kokiEndpoints.Labels
	kubeEndpoints.Annotations = kokiEndpoints.Annotations

	subsets, err := revertSubsets(kokiEndpoints.Subsets)
	if err != nil {
		return nil, err
	}
	kubeEndpoints.Subsets = subsets

	return kubeEndpoints, nil
}

func revertSubsets(subsets []types.EndpointSubset) ([]v1.EndpointSubset, error) {
	var kubeSubsets []v1.EndpointSubset

	for i := range subsets {
		subset := subsets[i]

		addrs, err := revertEndpointAddresses(subset.Addresses)
		if err != nil {
			return nil, err
		}

		unreadyAddrs, err := revertEndpointAddresses(subset.NotReadyAddresses)
		if err != nil {
			return nil, err
		}

		ports, err := revertEndpointPorts(subset.Ports)
		if err != nil {
			return nil, err
		}

		kubeSubset := v1.EndpointSubset{
			Addresses:         addrs,
			NotReadyAddresses: unreadyAddrs,
			Ports:             ports,
		}

		kubeSubsets = append(kubeSubsets, kubeSubset)
	}

	return kubeSubsets, nil
}

func revertEndpointPorts(ports []string) ([]v1.EndpointPort, error) {
	var kubePorts []v1.EndpointPort

	for i := range ports {
		port := ports[i]

		fields := strings.Split(port, ":")

		protocol := v1.ProtocolTCP
		portVal := 0
		name := ""
		if len(fields) >= 2 {
			if strings.ToLower(fields[0]) == "udp" {
				protocol = v1.ProtocolUDP
			} else if strings.ToLower(fields[0]) != "tcp" {
				return nil, serrors.InvalidValueErrorf(fields[0], "invalid protocol")
			}

			portStr := strings.TrimPrefix(fields[1], "//")
			portV, err := strconv.Atoi(portStr)
			if err != nil {
				return nil, err
			}
			portVal = portV
		}

		if len(fields) == 3 {
			name = fields[2]
		}

		if len(fields) != 2 && len(fields) != 3 {
			return nil, serrors.InvalidValueErrorf(port, "invalid endpoints port format")
		}

		kubePort := v1.EndpointPort{
			Protocol: protocol,
			Port:     int32(portVal),
			Name:     name,
		}

		kubePorts = append(kubePorts, kubePort)
	}

	return kubePorts, nil
}

func revertEndpointAddresses(addrs []types.EndpointAddress) ([]v1.EndpointAddress, error) {
	var kubeAddrs []v1.EndpointAddress

	for i := range addrs {
		addr := addrs[i]

		target, err := revertTarget(addr.Target)
		if err != nil {
			return nil, err
		}

		kubeAddr := v1.EndpointAddress{
			IP:        addr.IP,
			Hostname:  addr.Hostname,
			NodeName:  addr.Nodename,
			TargetRef: target,
		}

		kubeAddrs = append(kubeAddrs, kubeAddr)
	}

	return kubeAddrs, nil
}

func revertTarget(target *types.ObjectReference) (*v1.ObjectReference, error) {
	if target == nil {
		return nil, nil
	}

	return &v1.ObjectReference{
		Kind:            target.Kind,
		Namespace:       target.Namespace,
		Name:            target.Name,
		UID:             kubeTypes.UID(target.UID),
		APIVersion:      target.Version,
		ResourceVersion: target.ResourceVersion,
		FieldPath:       target.FieldPath,
	}, nil
}
