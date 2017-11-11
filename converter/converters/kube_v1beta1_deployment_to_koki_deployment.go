package converters

import (
	"reflect"

	exts "k8s.io/api/extensions/v1beta1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/types"
)

func Convert_Kube_v1beta2_Deployment_to_Koki_Deployment(kubeDeployment *exts.Deployment) (*types.DeploymentWrapper, error) {
	var err error
	kokiDeployment := &types.Deployment{}

	kokiDeployment.Name = kubeDeployment.Name
	kokiDeployment.Namespace = kubeDeployment.Namespace
	kokiDeployment.Version = kubeDeployment.APIVersion
	kokiDeployment.Cluster = kubeDeployment.ClusterName
	kokiDeployment.Labels = kubeDeployment.Labels
	kokiDeployment.Annotations = kubeDeployment.Annotations

	kubeSpec := &kubeDeployment.Spec
	kokiDeployment.Replicas = kubeSpec.Replicas
	kokiTemplate, err := convertTemplate(&kubeSpec.Template)
	if err != nil {
		return nil, err
	}
	kokiDeployment.Template = *kokiTemplate

	kokiDeployment.Recreate, kokiDeployment.MaxUnavailable, kokiDeployment.MaxSurge = convertDeploymentStrategy(kubeSpec.Strategy)

	kokiDeployment.MinReadySeconds = kubeSpec.MinReadySeconds
	kokiDeployment.RevisionHistoryLimit = kubeSpec.RevisionHistoryLimit
	kokiDeployment.Paused = kubeSpec.Paused
	kokiDeployment.ProgressDeadlineSeconds = kubeSpec.ProgressDeadlineSeconds

	if !reflect.DeepEqual(kubeDeployment.Status, exts.DeploymentStatus{}) {
		kokiDeployment.Status = &kubeDeployment.Status
	}

	return &types.DeploymentWrapper{
		Deployment: *kokiDeployment,
	}, nil
}

func convertDeploymentStrategy(kubeStrategy exts.DeploymentStrategy) (isRecreate bool, maxUnavailable, maxSurge *intstr.IntOrString) {
	if kubeStrategy.Type == exts.RecreateDeploymentStrategyType {
		return true, nil, nil
	}

	if rollingUpdate := kubeStrategy.RollingUpdate; rollingUpdate != nil {
		return false, rollingUpdate.MaxUnavailable, rollingUpdate.MaxSurge
	}

	return false, nil, nil
}
