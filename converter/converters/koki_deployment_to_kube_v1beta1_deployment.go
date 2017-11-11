package converters

import (
	exts "k8s.io/api/extensions/v1beta1"

	"github.com/koki/short/types"
)

func Convert_Koki_Deployment_to_Kube_v1beta1_Deployment(deployment *types.DeploymentWrapper) (*exts.Deployment, error) {
	var err error
	kubeDeployment := &exts.Deployment{}
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

func revertDeploymentStrategy(kokiDeployment *types.Deployment) exts.DeploymentStrategy {
	if kokiDeployment.Recreate {
		return exts.DeploymentStrategy{
			Type: exts.RecreateDeploymentStrategyType,
		}
	}

	var rollingUpdateConfig *exts.RollingUpdateDeployment
	if kokiDeployment.MaxUnavailable != nil || kokiDeployment.MaxSurge != nil {
		rollingUpdateConfig = &exts.RollingUpdateDeployment{
			MaxUnavailable: kokiDeployment.MaxUnavailable,
			MaxSurge:       kokiDeployment.MaxSurge,
		}
	}

	return exts.DeploymentStrategy{
		Type:          exts.RollingUpdateDeploymentStrategyType,
		RollingUpdate: rollingUpdateConfig,
	}
}
