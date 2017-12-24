package converters

import (
	autoscaling "k8s.io/api/autoscaling/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_HPA_to_Koki(kube *autoscaling.HorizontalPodAutoscaler) (*types.HorizontalPodAutoscalerWrapper, error) {
	koki := &types.HorizontalPodAutoscaler{}

	koki.Name = kube.Name
	koki.Namespace = kube.Namespace
	koki.Version = kube.APIVersion
	koki.Cluster = kube.ClusterName
	koki.Labels = kube.Labels
	koki.Annotations = kube.Annotations

	koki.HorizontalPodAutoscalerSpec = convertHPASpec(kube.Spec)
	koki.HorizontalPodAutoscalerStatus = convertHPAStatus(kube.Status)

	return &types.HorizontalPodAutoscalerWrapper{
		HPA: *koki,
	}, nil
}

func convertHPASpec(kubeSpec autoscaling.HorizontalPodAutoscalerSpec) types.HorizontalPodAutoscalerSpec {
	return types.HorizontalPodAutoscalerSpec{
		ScaleTargetRef: convertCrossVersionObjectReference(kubeSpec.ScaleTargetRef),

		MinReplicas:                    kubeSpec.MinReplicas,
		MaxReplicas:                    kubeSpec.MaxReplicas,
		TargetCPUUtilizationPercentage: kubeSpec.TargetCPUUtilizationPercentage,
	}
}

func convertHPAStatus(kubeStatus autoscaling.HorizontalPodAutoscalerStatus) types.HorizontalPodAutoscalerStatus {
	return types.HorizontalPodAutoscalerStatus{
		ObservedGeneration:              kubeStatus.ObservedGeneration,
		LastScaleTime:                   kubeStatus.LastScaleTime,
		CurrentReplicas:                 kubeStatus.CurrentReplicas,
		DesiredReplicas:                 kubeStatus.DesiredReplicas,
		CurrentCPUUtilizationPercentage: kubeStatus.CurrentCPUUtilizationPercentage,
	}
}

func convertCrossVersionObjectReference(kubeRef autoscaling.CrossVersionObjectReference) types.CrossVersionObjectReference {
	return types.CrossVersionObjectReference{
		Kind:       kubeRef.Kind,
		Name:       kubeRef.Name,
		APIVersion: kubeRef.APIVersion,
	}
}
