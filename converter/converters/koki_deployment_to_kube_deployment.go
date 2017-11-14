package converters

import (
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"

	"github.com/ghodss/yaml"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_Deployment_to_Kube_Deployment(deployment *types.DeploymentWrapper) (interface{}, error) {
	kubeDeployment, err := Convert_Koki_Deployment_to_Kube_apps_v1beta2_Deployment(deployment)
	if err != nil {
		return nil, err
	}

	b, err := yaml.Marshal(kubeDeployment)
	if err != nil {
		return nil, err
	}

	switch deployment.Deployment.Version {
	case "apps/v1beta1":
		return nil, util.PrettyTypeError(deployment, "unsupported version")
	case "apps/v1beta2":
		return nil, util.PrettyTypeError(deployment, "unsupported version")
	default:
		return nil, util.PrettyTypeError(deployment, "unsupported version")
	}
}

func Convert_Koki_Deployment_to_Kube_apps_v1beta2_Deployment(deployment *types.DeploymentWrapper) (*appsv1beta2.Deployment, error) {
	var err error
	kubeDeployment := &appsv1beta2.Deployment{}
	kokiDeployment := &deployment.Deployment

	kubeDeployment.Name = kokiDeployment.Name
	kubeDeployment.Namespace = kokiDeployment.Namespace
	kubeDeployment.APIVersion = kokiDeployment.Version
	kubeDeployment.Kind = "Deployment"
	kubeDeployment.ClusterName = kokiDeployment.Cluster
	kubeDeployment.Labels = kokiDeployment.Labels
	kubeDeployment.Annotations = kokiDeployment.Annotations

	kubeSpec := &kubeDeployment.Spec
	kubeSpec.Replicas = kokiDeployment.Replicas

	// Setting the Selector and Template is identical to ReplicaSet

	var templateLabelsOverride map[string]string
	var kokiTemplateLabels map[string]string
	if kokiDeployment.TemplateMetadata != nil {
		kokiTemplateLabels = kokiDeployment.TemplateMetadata.Labels
	}
	kubeSpec.Selector, templateLabelsOverride, err = revertRSSelector(kokiDeployment.Name, kokiDeployment.Selector, kokiTemplateLabels)
	if err != nil {
		return nil, err
	}

	kubeTemplate, err := revertTemplate(kokiDeployment.GetTemplate())
	if err != nil {
		return nil, err
	}
	if kubeTemplate == nil {
		return nil, util.TypeValueErrorf(kokiDeployment, "missing pod template")
	}
	kubeTemplate.Labels = templateLabelsOverride
	kubeSpec.Template = *kubeTemplate

	// End Selector/Template section.

	kubeSpec.Strategy = revertDeploymentStrategy(kokiDeployment)

	kubeSpec.MinReadySeconds = kokiDeployment.MinReadySeconds
	kubeSpec.RevisionHistoryLimit = kokiDeployment.RevisionHistoryLimit
	kubeSpec.Paused = kokiDeployment.Paused
	kubeSpec.ProgressDeadlineSeconds = kokiDeployment.ProgressDeadlineSeconds

	if kokiDeployment.Status != nil {
		kubeDeployment.Status = *kokiDeployment.Status
	}

	return kubeDeployment, nil
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
