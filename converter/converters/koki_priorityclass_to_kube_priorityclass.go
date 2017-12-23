package converters

import (
	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"

	"github.com/koki/short/types"
)

func Convert_Koki_PriorityClass_to_Kube_PriorityClass(priorityClass *types.PriorityClassWrapper) (*schedulingv1alpha1.PriorityClass, error) {
	kubePriorityClass := &schedulingv1alpha1.PriorityClass{}
	kokiPriorityClass := &priorityClass.PriorityClass

	kubePriorityClass.Name = kokiPriorityClass.Name
	kubePriorityClass.Namespace = kokiPriorityClass.Namespace
	if len(kokiPriorityClass.Version) == 0 {
		kubePriorityClass.APIVersion = "scheduling/v1alpha1"
	} else {
		kubePriorityClass.APIVersion = kokiPriorityClass.Version
	}
	kubePriorityClass.Kind = "PriorityClass"
	kubePriorityClass.ClusterName = kokiPriorityClass.Cluster
	kubePriorityClass.Labels = kokiPriorityClass.Labels
	kubePriorityClass.Annotations = kokiPriorityClass.Annotations

	kubePriorityClass.Value = kokiPriorityClass.Value
	kubePriorityClass.GlobalDefault = kokiPriorityClass.GlobalDefault
	kubePriorityClass.Description = kokiPriorityClass.Description

	return kubePriorityClass, nil
}
