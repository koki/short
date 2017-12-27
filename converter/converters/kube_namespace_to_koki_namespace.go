package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_Namespace_to_Koki_Namespace(kubeNamespace *v1.Namespace) (*types.NamespaceWrapper, error) {
	kokiWrapper := &types.NamespaceWrapper{}
	kokiNamespace := &kokiWrapper.Namespace

	kokiNamespace.Name = kubeNamespace.Name
	kokiNamespace.Namespace = kubeNamespace.Namespace
	kokiNamespace.Version = kubeNamespace.APIVersion
	kokiNamespace.Cluster = kubeNamespace.ClusterName
	kokiNamespace.Labels = kubeNamespace.Labels
	kokiNamespace.Annotations = kubeNamespace.Annotations

	finalizers, err := convertNamespaceSpec(kubeNamespace.Spec)
	if err != nil {
		return nil, err
	}
	kokiNamespace.Finalizers = finalizers

	phase, err := convertNamespaceStatus(kubeNamespace.Status)
	if err != nil {
		return nil, err
	}
	kokiNamespace.Phase = phase

	return kokiWrapper, nil
}

func convertNamespaceSpec(kubeSpec v1.NamespaceSpec) ([]types.FinalizerName, error) {
	var kokiFinalizers []types.FinalizerName

	for i := range kubeSpec.Finalizers {
		kubeFinalizer := kubeSpec.Finalizers[i]

		var kokiFinalizer types.FinalizerName
		switch kubeFinalizer {
		case v1.FinalizerKubernetes:
			kokiFinalizer = types.FinalizerKubernetes
		default:
			return nil, serrors.InvalidValueErrorf(kubeFinalizer, "unrecognized finalizer")
		}

		kokiFinalizers = append(kokiFinalizers, kokiFinalizer)
	}

	return kokiFinalizers, nil
}

func convertNamespaceStatus(kubeStatus v1.NamespaceStatus) (types.NamespacePhase, error) {
	switch kubeStatus.Phase {
	case v1.NamespaceActive:
		return types.NamespaceActive, nil
	case v1.NamespaceTerminating:
		return types.NamespaceTerminating, nil
	case "":
		return "", nil
	}

	return "", serrors.InvalidValueErrorf(kubeStatus.Phase, "invalid phase")
}
