package converters

import (
	"reflect"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/runtime"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	"github.com/ghodss/yaml"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
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
		return nil, util.InvalidInstanceErrorf(kubeDeployment, "couldn't serialize kube Deployment after setting apiVersion to apps/v1beta2: %s", err.Error())
	}

	// Deserialize the "generic" kube Deployment
	genericDeployment, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, util.InvalidInstanceErrorf(string(b), "couldn't deserialize 'generic' kube Deployment: %s", err.Error())
	}

	if genericDeployment, ok := genericDeployment.(*appsv1beta2.Deployment); ok {
		kokiWrapper, err := Convert_Kube_v1beta2_Deployment_to_Koki_Deployment(genericDeployment)
		if err != nil {
			return nil, err
		}

		kokiDeployment := &kokiWrapper.Deployment

		kokiDeployment.Version = groupVersionString

		// Perform version-specific initialization here.

		return kokiWrapper, nil
	}

	return nil, util.InvalidInstanceErrorf(genericDeployment, "didn't deserialize 'generic' kube Deployment as apps/v1beta2.Deployment")
}

func Convert_Kube_v1beta2_Deployment_to_Koki_Deployment(kubeDeployment *appsv1beta2.Deployment) (*types.DeploymentWrapper, error) {
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

	// Setting the Selector and Template is identical to ReplicaSet

	// Fill out the Selector and Template.Labels.
	// If kubeDeployment only has Template.Labels, we pull it up to Selector.
	var templateLabelsOverride map[string]string
	kokiDeployment.Selector, templateLabelsOverride, err = convertRSLabelSelector(kubeSpec.Selector, kubeSpec.Template.Labels)
	if err != nil {
		return nil, err
	}

	// Build a Pod from the kube Template. Use it to set the koki Template.
	kokiPod, err := convertRSTemplate(&kubeSpec.Template)
	if err != nil {
		return nil, err
	}
	kokiPod.Labels = templateLabelsOverride
	kokiDeployment.SetTemplate(kokiPod)

	// End Selector/Template section.

	kokiDeployment.Recreate, kokiDeployment.MaxUnavailable, kokiDeployment.MaxSurge = convertDeploymentStrategy(kubeSpec.Strategy)

	kokiDeployment.MinReadySeconds = kubeSpec.MinReadySeconds
	kokiDeployment.RevisionHistoryLimit = kubeSpec.RevisionHistoryLimit
	kokiDeployment.Paused = kubeSpec.Paused
	kokiDeployment.ProgressDeadlineSeconds = kubeSpec.ProgressDeadlineSeconds

	if !reflect.DeepEqual(kubeDeployment.Status, appsv1beta2.DeploymentStatus{}) {
		kokiDeployment.Status = &kubeDeployment.Status
	}

	return &types.DeploymentWrapper{
		Deployment: *kokiDeployment,
	}, nil
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
