package converters

import (
	apps "k8s.io/api/apps/v1beta2"

	"github.com/koki/short/types"
)

func Convert_Koki_Deployment_to_Kube_v1beta2_Deployment(deployment *types.DeploymentWrapper) (*apps.Deployment, error) {
	var err error
	kubeDeployment := &apps.Deployment{}
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

	kubeTemplate, err := revertTemplate(&kokiDeployment.Template)
	if err != nil {
		return nil, err
	}
	kubeSpec.Template = *kubeTemplate

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

func revertDeploymentStrategy(kokiDeployment *types.Deployment) apps.DeploymentStrategy {
	if kokiDeployment.Recreate {
		return apps.DeploymentStrategy{
			Type: apps.RecreateDeploymentStrategyType,
		}
	}

	var rollingUpdateConfig *apps.RollingUpdateDeployment
	if kokiDeployment.MaxUnavailable != nil || kokiDeployment.MaxSurge != nil {
		rollingUpdateConfig = &apps.RollingUpdateDeployment{
			MaxUnavailable: kokiDeployment.MaxUnavailable,
			MaxSurge:       kokiDeployment.MaxSurge,
		}
	}

	return apps.DeploymentStrategy{
		Type:          apps.RollingUpdateDeploymentStrategyType,
		RollingUpdate: rollingUpdateConfig,
	}
}
