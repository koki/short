package converters

import (
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	exts "k8s.io/api/extensions/v1beta1"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_Deployment_to_Kube_Deployment(deployment *types.DeploymentWrapper) (interface{}, error) {
	// Perform version-agnostic conversion into apps/v1beta2 Deployment.
	kubeDeployment, err := Convert_Koki_Deployment_to_Kube_apps_v1beta2_Deployment(deployment)
	if err != nil {
		return nil, err
	}

	// Serialize the "generic" kube Deployment.
	b, err := yaml.Marshal(kubeDeployment)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, kubeDeployment, "couldn't serialize 'generic' kube Deployment")
	}

	// Deserialize a versioned kube Deployment using its apiVersion.
	versionedDeployment, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, err
	}

	switch versionedDeployment := versionedDeployment.(type) {
	case *appsv1beta1.Deployment:
		// Perform apps/v1beta1-specific initialization here.
	case *appsv1beta2.Deployment:
		// Perform apps/v1beta2-specific initialization here.
	case *exts.Deployment:
		// Perform exts/v1beta1-specific initialization here.
	default:
		return nil, serrors.TypeErrorf(versionedDeployment, "deserialized the manifest, but not as a supported kube Deployment")
	}

	return versionedDeployment, nil
}

func Convert_Koki_Deployment_to_Kube_apps_v1beta2_Deployment(deployment *types.DeploymentWrapper) (*appsv1beta2.Deployment, error) {
	var err error
	kubeDeployment := &appsv1beta2.Deployment{}
	kokiDeployment := &deployment.Deployment

	kubeDeployment.Name = kokiDeployment.Name
	kubeDeployment.Namespace = kokiDeployment.Namespace
	if len(kokiDeployment.Version) == 0 {
		kubeDeployment.APIVersion = "extensions/v1beta1"
	} else {
		kubeDeployment.APIVersion = kokiDeployment.Version
	}
	kubeDeployment.Kind = "Deployment"
	kubeDeployment.ClusterName = kokiDeployment.Cluster
	kubeDeployment.Labels = kokiDeployment.Labels
	kubeDeployment.Annotations = kokiDeployment.Annotations

	kubeSpec := &kubeDeployment.Spec
	kubeSpec.Replicas = kokiDeployment.Replicas

	// Setting the Selector and Template is identical to ReplicaSet
	// Get the right Selector and Template Labels.
	var templateLabelsOverride map[string]string
	var kokiTemplateLabels map[string]string
	if kokiDeployment.TemplateMetadata != nil {
		kokiTemplateLabels = kokiDeployment.TemplateMetadata.Labels
	}
	kubeSpec.Selector, templateLabelsOverride, err = revertRSSelector(kokiDeployment.Name, kokiDeployment.Selector, kokiTemplateLabels)
	if err != nil {
		return nil, err
	}
	// Set the right Labels before we fill in the Pod template with this metadata.
	kokiDeployment.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, kokiDeployment.TemplateMetadata)

	// Fill in the rest of the Pod template.
	kubeTemplate, err := revertTemplate(kokiDeployment.TemplateMetadata, kokiDeployment.PodTemplate)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}
	if kubeTemplate == nil {
		return nil, serrors.InvalidInstanceErrorf(kokiDeployment, "missing pod template")
	}
	kubeSpec.Template = *kubeTemplate

	// End Selector/Template section.

	kubeSpec.Strategy = revertDeploymentStrategy(kokiDeployment)

	kubeSpec.MinReadySeconds = kokiDeployment.MinReadySeconds
	kubeSpec.RevisionHistoryLimit = kokiDeployment.RevisionHistoryLimit
	kubeSpec.Paused = kokiDeployment.Paused
	kubeSpec.ProgressDeadlineSeconds = kokiDeployment.ProgressDeadlineSeconds

	kubeDeployment.Status, err = revertDeploymentStatus(kokiDeployment.DeploymentStatus)
	if err != nil {
		return nil, err
	}

	return kubeDeployment, nil
}

func revertDeploymentStatus(kokiStatus types.DeploymentStatus) (appsv1beta2.DeploymentStatus, error) {
	conditions, err := revertDeploymentConditions(kokiStatus.Conditions)
	if err != nil {
		return appsv1beta2.DeploymentStatus{}, err
	}
	return appsv1beta2.DeploymentStatus{
		ObservedGeneration:  kokiStatus.ObservedGeneration,
		Replicas:            kokiStatus.Replicas.Total,
		UpdatedReplicas:     kokiStatus.Replicas.Updated,
		ReadyReplicas:       kokiStatus.Replicas.Ready,
		AvailableReplicas:   kokiStatus.Replicas.Available,
		UnavailableReplicas: kokiStatus.Replicas.Unavailable,
		Conditions:          conditions,
		CollisionCount:      kokiStatus.CollisionCount,
	}, nil
}

func revertDeploymentConditions(kokiConditions []types.DeploymentCondition) ([]appsv1beta2.DeploymentCondition, error) {
	if len(kokiConditions) == 0 {
		return nil, nil
	}

	kubeConditions := make([]appsv1beta2.DeploymentCondition, len(kokiConditions))
	for i, condition := range kokiConditions {
		status, err := revertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "deployment conditions[%d]", i)
		}
		conditionType, err := revertDeploymentConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "deployment conditions[%d]", i)
		}
		kubeConditions[i] = appsv1beta2.DeploymentCondition{
			Type:               conditionType,
			Status:             status,
			LastUpdateTime:     condition.LastUpdateTime,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}

	return kubeConditions, nil
}

func revertDeploymentConditionType(kokiType types.DeploymentConditionType) (appsv1beta2.DeploymentConditionType, error) {
	switch kokiType {
	case types.DeploymentAvailable:
		return appsv1beta2.DeploymentAvailable, nil
	case types.DeploymentProgressing:
		return appsv1beta2.DeploymentProgressing, nil
	case types.DeploymentReplicaFailure:
		return appsv1beta2.DeploymentReplicaFailure, nil
	default:
		return appsv1beta2.DeploymentReplicaFailure, serrors.InvalidValueErrorf(kokiType, "unrecognized deployment condition type")
	}
}

func revertDeploymentStrategy(kokiDeployment *types.Deployment) appsv1beta2.DeploymentStrategy {
	if kokiDeployment.Recreate {
		return appsv1beta2.DeploymentStrategy{
			Type: appsv1beta2.RecreateDeploymentStrategyType,
		}
	}

	var rollingUpdateConfig *appsv1beta2.RollingUpdateDeployment
	if kokiDeployment.MaxUnavailable != nil || kokiDeployment.MaxSurge != nil {
		rollingUpdateConfig = &appsv1beta2.RollingUpdateDeployment{
			MaxUnavailable: kokiDeployment.MaxUnavailable,
			MaxSurge:       kokiDeployment.MaxSurge,
		}
	}

	return appsv1beta2.DeploymentStrategy{
		Type:          appsv1beta2.RollingUpdateDeploymentStrategyType,
		RollingUpdate: rollingUpdateConfig,
	}
}
