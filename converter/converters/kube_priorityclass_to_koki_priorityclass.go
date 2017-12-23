package converters

import (
	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"

	"github.com/koki/short/types"
)

func Convert_Kube_PriorityClass_to_Koki_PriorityClass(kubePriorityClass *schedulingv1alpha1.PriorityClass) (*types.PriorityClassWrapper, error) {
	kokiWrapper := &types.PriorityClassWrapper{}
	kokiPriorityClass := &kokiWrapper.PriorityClass

	kokiPriorityClass.Name = kubePriorityClass.Name
	kokiPriorityClass.Namespace = kubePriorityClass.Namespace
	kokiPriorityClass.Version = kubePriorityClass.APIVersion
	kokiPriorityClass.Cluster = kubePriorityClass.ClusterName
	kokiPriorityClass.Labels = kubePriorityClass.Labels
	kokiPriorityClass.Annotations = kubePriorityClass.Annotations

	kokiPriorityClass.Description = kubePriorityClass.Description
	kokiPriorityClass.GlobalDefault = kubePriorityClass.GlobalDefault
	kokiPriorityClass.Value = kubePriorityClass.Value

	return kokiWrapper, nil
}
