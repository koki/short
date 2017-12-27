package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_Namespace_to_Kube_Namespace(kokiWrapper *types.NamespaceWrapper) (*v1.Namespace, error) {
	var err error
	kubeNamespace := &v1.Namespace{}
	kokiNamespace := kokiWrapper.Namespace

	kubeNamespace.Name = kokiNamespace.Name
	kubeNamespace.Namespace = kokiNamespace.Namespace
	if len(kokiNamespace.Version) == 0 {
		kubeNamespace.APIVersion = "v1"
	} else {
		kubeNamespace.APIVersion = kokiNamespace.Version
	}
	kubeNamespace.Kind = "Namespace"
	kubeNamespace.ClusterName = kokiNamespace.Cluster
	kubeNamespace.Labels = kokiNamespace.Labels
	kubeNamespace.Annotations = kokiNamespace.Annotations

	spec, err := revertNamespaceSpec(kokiNamespace)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "Namespace spec")
	}
	kubeNamespace.Spec = spec

	status, err := revertNamespaceStatus(kokiNamespace)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "Namespace status")
	}
	kubeNamespace.Status = status

	return kubeNamespace, nil
}

func revertNamespaceStatus(kokiNamespace types.Namespace) (v1.NamespaceStatus, error) {
	var kubeStatus v1.NamespaceStatus

	if kokiNamespace.Phase == "" {
		return kubeStatus, nil
	}

	var phase v1.NamespacePhase
	switch kokiNamespace.Phase {
	case types.NamespaceActive:
		phase = v1.NamespaceActive
	case types.NamespaceTerminating:
		phase = v1.NamespaceTerminating
	default:
		return kubeStatus, serrors.InvalidValueErrorf(kokiNamespace.Phase, "Invalid namespace phase")
	}
	kubeStatus.Phase = phase

	return kubeStatus, nil
}

func revertNamespaceSpec(kokiNamespace types.Namespace) (v1.NamespaceSpec, error) {
	var kubeSpec v1.NamespaceSpec
	var kubeFinalizers []v1.FinalizerName

	for i := range kokiNamespace.Finalizers {
		kokiFinalizer := kokiNamespace.Finalizers[i]

		kubeFinalizer, err := revertFinalizer(kokiFinalizer)
		if err != nil {
			return kubeSpec, err
		}

		kubeFinalizers = append(kubeFinalizers, kubeFinalizer)
	}
	kubeSpec.Finalizers = kubeFinalizers

	return kubeSpec, nil
}

func revertFinalizer(kokiFinalizer types.FinalizerName) (v1.FinalizerName, error) {
	if kokiFinalizer == "" {
		return "", nil
	}

	switch kokiFinalizer {
	case types.FinalizerKubernetes:
		return v1.FinalizerKubernetes, nil
	}

	return "", serrors.InvalidValueErrorf(kokiFinalizer, "unrecognized value")
}
