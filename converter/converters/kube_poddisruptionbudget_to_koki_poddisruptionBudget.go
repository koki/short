package converters

import (
	"fmt"

	policyv1beta1 "k8s.io/api/policy/v1beta1"

	"github.com/koki/short/types"
	"github.com/koki/short/util/floatstr"
)

func Convert_Kube_PodDisruptionBudget_to_Koki_PodDisruptionBudget(kubePDB *policyv1beta1.PodDisruptionBudget) (*types.PodDisruptionBudgetWrapper, error) {
	var err error
	kokiWrapper := &types.PodDisruptionBudgetWrapper{}
	kokiPDB := &kokiWrapper.PodDisruptionBudget

	kokiPDB.Name = kubePDB.Name
	kokiPDB.Namespace = kubePDB.Namespace
	kokiPDB.Version = kubePDB.APIVersion
	kokiPDB.Cluster = kubePDB.ClusterName
	kokiPDB.Labels = kubePDB.Labels
	kokiPDB.Annotations = kubePDB.Annotations

	err = convertPDBSpec(kubePDB.Spec, kokiPDB)
	if err != nil {
		return nil, err
	}

	err = convertPDBStatus(kubePDB.Status, kokiPDB)
	if err != nil {
		return nil, err
	}

	return kokiWrapper, nil
}

func convertPDBSpec(kubeSpec policyv1beta1.PodDisruptionBudgetSpec, kokiPDB *types.PodDisruptionBudget) error {
	if kokiPDB == nil {
		return fmt.Errorf("Writing to uninitialized Pod Disruption Budget pointer")
	}

	if kubeSpec.MaxUnavailable != nil {
		kokiPDB.MaxEvictionsAllowed = floatstr.Parse(kubeSpec.MaxUnavailable.String())
	}

	if kubeSpec.MinAvailable != nil {
		kokiPDB.MinPodsRequired = floatstr.Parse(kubeSpec.MinAvailable.String())
	}

	// Fill out the Selector
	selector, _, err := convertRSLabelSelector(kubeSpec.Selector, nil)
	if err != nil {
		return err
	}

	if selector != nil && (selector.Labels != nil || selector.Shorthand != "") {
		kokiPDB.Selector = selector
	}

	return nil
}

func convertPDBStatus(kubeStatus policyv1beta1.PodDisruptionBudgetStatus, kokiPDB *types.PodDisruptionBudget) error {
	if kokiPDB == nil {
		return fmt.Errorf("Writing to uninitialized Pod Disruption Budget pointer")
	}

	kokiPDB.ObservedGeneration = kubeStatus.ObservedGeneration
	kokiPDB.DisruptedPods = kubeStatus.DisruptedPods
	kokiPDB.PodDisruptionsAllowed = kubeStatus.PodDisruptionsAllowed
	kokiPDB.CurrentHealthy = kubeStatus.CurrentHealthy
	kokiPDB.DesiredHealthy = kubeStatus.DesiredHealthy
	kokiPDB.ExpectedPods = kubeStatus.ExpectedPods
	return nil
}
