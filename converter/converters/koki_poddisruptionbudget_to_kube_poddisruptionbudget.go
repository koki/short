package converters

import (
	policyv1beta1 "k8s.io/api/policy/v1beta1"

	"github.com/koki/short/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Convert_Koki_PodDisruptionBudget_to_Kube_PodDisruptionBudget(podDisruptionBudget *types.PodDisruptionBudgetWrapper) (*policyv1beta1.PodDisruptionBudget, error) {
	kubePDB := &policyv1beta1.PodDisruptionBudget{}
	kokiPDB := &podDisruptionBudget.PodDisruptionBudget

	kubePDB.Name = kokiPDB.Name
	kubePDB.Namespace = kokiPDB.Namespace
	if len(kokiPDB.Version) == 0 {
		kubePDB.APIVersion = "policy/v1beta1"
	} else {
		kubePDB.APIVersion = kokiPDB.Version
	}
	kubePDB.Kind = "PodDisruptionBudget"
	kubePDB.ClusterName = kokiPDB.Cluster
	kubePDB.Labels = kokiPDB.Labels
	kubePDB.Annotations = kokiPDB.Annotations

	kubeSpec, err := revertKokiPDBSpec(kokiPDB)
	if err != nil {
		return nil, err
	}
	kubePDB.Spec = kubeSpec

	kubePDB.Status = revertKokiPDBStatus(kokiPDB)

	return kubePDB, nil
}

func revertKokiPDBSpec(kokiPDB *types.PodDisruptionBudget) (policyv1beta1.PodDisruptionBudgetSpec, error) {
	var kubeSpec policyv1beta1.PodDisruptionBudgetSpec

	if kokiPDB.MaxEvictionsAllowed != nil {
		intstrMaxUnavailable := intstr.Parse(kokiPDB.MaxEvictionsAllowed.String())
		kubeSpec.MaxUnavailable = &intstrMaxUnavailable
	}

	if kokiPDB.MinPodsRequired != nil {
		intstrMinAvailable := intstr.Parse(kokiPDB.MinPodsRequired.String())
		kubeSpec.MinAvailable = &intstrMinAvailable
	}

	selector, _, err := revertRSSelector(kokiPDB.Name, kokiPDB.Selector, nil)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.Selector = selector

	return kubeSpec, nil
}

func revertKokiPDBStatus(kokiPDB *types.PodDisruptionBudget) policyv1beta1.PodDisruptionBudgetStatus {
	var kubeStatus policyv1beta1.PodDisruptionBudgetStatus

	kubeStatus.ObservedGeneration = kokiPDB.ObservedGeneration
	kubeStatus.DisruptedPods = kokiPDB.DisruptedPods
	kubeStatus.PodDisruptionsAllowed = kokiPDB.PodDisruptionsAllowed
	kubeStatus.CurrentHealthy = kokiPDB.CurrentHealthy
	kubeStatus.DesiredHealthy = kokiPDB.DesiredHealthy
	kubeStatus.ExpectedPods = kokiPDB.ExpectedPods

	return kubeStatus
}
