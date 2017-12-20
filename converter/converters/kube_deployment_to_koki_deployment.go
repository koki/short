package converters

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/runtime"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_Deployment_to_Koki_Deployment(kubeDeployment runtime.Object) (*types.DeploymentWrapper, error) {
	groupVersionKind := kubeDeployment.GetObjectKind().GroupVersionKind()
	groupVersionString := groupVersionKind.GroupVersion().String()
	groupVersionKind.Version = "v1beta2"
	groupVersionKind.Group = "apps"
	kubeDeployment.GetObjectKind().SetGroupVersionKind(groupVersionKind)

	// Serialize as v1beta2
	b, err := yaml.Marshal(kubeDeployment)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, kubeDeployment, "couldn't serialize kube Deployment after setting apiVersion to apps/v1beta2")
	}

	// Deserialize the "generic" kube Deployment
	genericDeployment, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, string(b), "couldn't deserialize 'generic' kube Deployment")
	}

	if genericDeployment, ok := genericDeployment.(*appsv1beta2.Deployment); ok {
		kokiWrapper, err := Convert_Kube_v1beta2_Deployment_to_Koki_Deployment(genericDeployment)
		if err != nil {
			return nil, err
		}

		kokiDeployment := &kokiWrapper.Deployment

		kokiDeployment.Version = groupVersionString

		return kokiWrapper, nil
	}

	return nil, serrors.InvalidInstanceErrorf(genericDeployment, "didn't deserialize 'generic' kube Deployment as apps/v1beta2.Deployment")
}

func Convert_Kube_v1beta2_Deployment_to_Koki_Deployment(kubeDeployment *appsv1beta2.Deployment) (*types.DeploymentWrapper, error) {
	kokiDeployment := &types.Deployment{}

	kokiDeployment.Name = kubeDeployment.Name
	kokiDeployment.Namespace = kubeDeployment.Namespace
	kokiDeployment.Version = kubeDeployment.APIVersion
	kokiDeployment.Cluster = kubeDeployment.ClusterName
	kokiDeployment.Labels = kubeDeployment.Labels
	kokiDeployment.Annotations = kubeDeployment.Annotations

	kubeSpec := &kubeDeployment.Spec
	kokiDeployment.Replicas = kubeSpec.Replicas

	// Setting the Selector and Template is identical to ReplicaSet

	// Fill out the Selector and Template.Labels.
	// If kubeDeployment only has Template.Labels, we pull it up to Selector.
	selector, templateLabelsOverride, err := convertRSLabelSelector(kubeSpec.Selector, kubeSpec.Template.Labels)
	if err != nil {
		return nil, err
	}

	if selector != nil && (selector.Labels != nil || selector.Shorthand != "") {
		kokiDeployment.Selector = selector
	}

	// Build a Pod from the kube Template. Use it to set the koki Template.
	meta, template, err := convertTemplate(kubeSpec.Template)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}
	kokiDeployment.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, meta)
	kokiDeployment.PodTemplate = template

	// End Selector/Template section.

	kokiDeployment.Recreate, kokiDeployment.MaxUnavailable, kokiDeployment.MaxSurge = convertDeploymentStrategy(kubeSpec.Strategy)

	kokiDeployment.MinReadySeconds = kubeSpec.MinReadySeconds
	kokiDeployment.RevisionHistoryLimit = kubeSpec.RevisionHistoryLimit
	kokiDeployment.Paused = kubeSpec.Paused
	kokiDeployment.ProgressDeadlineSeconds = kubeSpec.ProgressDeadlineSeconds

	kokiDeployment.DeploymentStatus, err = convertDeploymentStatus(kubeDeployment.Status)
	if err != nil {
		return nil, err
	}

	return &types.DeploymentWrapper{
		Deployment: *kokiDeployment,
	}, nil
}

func convertDeploymentStatus(kubeStatus appsv1beta2.DeploymentStatus) (types.DeploymentStatus, error) {
	conditions, err := convertDeploymentConditions(kubeStatus.Conditions)
	if err != nil {
		return types.DeploymentStatus{}, err
	}
	return types.DeploymentStatus{
		ObservedGeneration: kubeStatus.ObservedGeneration,
		Replicas: types.DeploymentReplicasStatus{
			Total:       kubeStatus.Replicas,
			Updated:     kubeStatus.UpdatedReplicas,
			Ready:       kubeStatus.ReadyReplicas,
			Available:   kubeStatus.AvailableReplicas,
			Unavailable: kubeStatus.UnavailableReplicas,
		},
		Conditions:     conditions,
		CollisionCount: kubeStatus.CollisionCount,
	}, nil
}

func convertDeploymentConditions(kubeConditions []appsv1beta2.DeploymentCondition) ([]types.DeploymentCondition, error) {
	if len(kubeConditions) == 0 {
		return nil, nil
	}

	kokiConditions := make([]types.DeploymentCondition, len(kubeConditions))
	for i, condition := range kubeConditions {
		status, err := convertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "deployment conditions[%d]", i)
		}
		conditionType, err := convertDeploymentConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "deployment conditions[%d]", i)
		}
		kokiConditions[i] = types.DeploymentCondition{
			Type:               conditionType,
			Status:             status,
			LastUpdateTime:     condition.LastUpdateTime,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}

	return kokiConditions, nil
}

func convertDeploymentConditionType(kubeType appsv1beta2.DeploymentConditionType) (types.DeploymentConditionType, error) {
	switch kubeType {
	case appsv1beta2.DeploymentAvailable:
		return types.DeploymentAvailable, nil
	case appsv1beta2.DeploymentProgressing:
		return types.DeploymentProgressing, nil
	case appsv1beta2.DeploymentReplicaFailure:
		return types.DeploymentReplicaFailure, nil
	default:
		return types.DeploymentReplicaFailure, serrors.InvalidValueErrorf(kubeType, "unrecognized deployment condition type")
	}
}

func convertDeploymentStrategy(kubeStrategy appsv1beta2.DeploymentStrategy) (isRecreate bool, maxUnavailable, maxSurge *intstr.IntOrString) {
	if kubeStrategy.Type == appsv1beta2.RecreateDeploymentStrategyType {
		return true, nil, nil
	}

	if rollingUpdate := kubeStrategy.RollingUpdate; rollingUpdate != nil {
		return false, rollingUpdate.MaxUnavailable, rollingUpdate.MaxSurge
	}

	return false, nil, nil
}
