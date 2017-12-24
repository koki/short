package converters

import (
	autoscaling "k8s.io/api/autoscaling/v1"

	"github.com/koki/short/types"
)

func Convert_Koki_HPA_to_Kube(wrapper *types.HorizontalPodAutoscalerWrapper) (*autoscaling.HorizontalPodAutoscaler, error) {
	kube := &autoscaling.HorizontalPodAutoscaler{}
	koki := wrapper.HPA

	kube.Name = koki.Name
	kube.Namespace = koki.Namespace
	if len(koki.Version) == 0 {
		kube.APIVersion = "autoscaling/v1"
	} else {
		kube.APIVersion = koki.Version
	}
	kube.Kind = "HorizontalPodAutoscaler"
	kube.ClusterName = koki.Cluster
	kube.Labels = koki.Labels
	kube.Annotations = koki.Annotations

	kube.Spec = revertHPASpec(koki.HorizontalPodAutoscalerSpec)
	kube.Status = revertHPAStatus(koki.HorizontalPodAutoscalerStatus)

	return kube, nil
}

func revertHPASpec(kokiSpec types.HorizontalPodAutoscalerSpec) autoscaling.HorizontalPodAutoscalerSpec {
	return autoscaling.HorizontalPodAutoscalerSpec{
		ScaleTargetRef:                 revertCrossVersionObjectReference(kokiSpec.ScaleTargetRef),
		MinReplicas:                    kokiSpec.MinReplicas,
		MaxReplicas:                    kokiSpec.MaxReplicas,
		TargetCPUUtilizationPercentage: kokiSpec.TargetCPUUtilizationPercentage,
	}
}

func revertHPAStatus(kokiStatus types.HorizontalPodAutoscalerStatus) autoscaling.HorizontalPodAutoscalerStatus {
	return autoscaling.HorizontalPodAutoscalerStatus{
		ObservedGeneration:              kokiStatus.ObservedGeneration,
		LastScaleTime:                   kokiStatus.LastScaleTime,
		CurrentReplicas:                 kokiStatus.CurrentReplicas,
		DesiredReplicas:                 kokiStatus.DesiredReplicas,
		CurrentCPUUtilizationPercentage: kokiStatus.CurrentCPUUtilizationPercentage,
	}
}

func revertCrossVersionObjectReference(kokiRef types.CrossVersionObjectReference) autoscaling.CrossVersionObjectReference {
	return autoscaling.CrossVersionObjectReference{
		Kind:       kokiRef.Kind,
		Name:       kokiRef.Name,
		APIVersion: kokiRef.APIVersion,
	}
}
