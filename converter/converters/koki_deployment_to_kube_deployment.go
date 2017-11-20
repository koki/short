package converters

import (
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	exts "k8s.io/api/extensions/v1beta1"

	"github.com/ghodss/yaml"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
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
		return nil, util.InvalidValueErrorf(kubeDeployment, "couldn't serialize 'generic' kube Deployment: %s", err.Error())
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
		return nil, util.TypeErrorf(versionedDeployment, "deserialized the manifest, but not as a supported kube Deployment")
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
		return nil, util.InvalidInstanceErrorf(kokiDeployment, "missing pod template")
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
