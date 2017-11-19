package converters

import (
	"reflect"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/ghodss/yaml"

	"github.com/koki/short/parser"
	"github.com/koki/short/parser/expressions"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Kube_ReplicaSet_to_Koki_ReplicaSet(kubeRS runtime.Object) (*types.ReplicaSetWrapper, error) {
	groupVersionKind := kubeRS.GetObjectKind().GroupVersionKind()
	groupVersionString := groupVersionKind.GroupVersion().String()
	groupVersionKind.Version = "v1beta2"
	groupVersionKind.Group = "apps"
	kubeRS.GetObjectKind().SetGroupVersionKind(groupVersionKind)

	// Serialize as v1beta2
	b, err := yaml.Marshal(kubeRS)
	if err != nil {
		return nil, util.InvalidInstanceErrorf(kubeRS, "couldn't serialize kube ReplicaSet after setting apiVersion to apps/v1beta2: %s", err.Error())
	}

	// Deserialize the "generic" kube ReplicaSet
	genericReplicaSet, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, util.InvalidValueErrorf(string(b), "couldn't deserialize 'generic' kube ReplicaSet: %s", err.Error())
	}

	if genericReplicaSet, ok := genericReplicaSet.(*appsv1beta2.ReplicaSet); ok {
		kokiWrapper, err := Convert_Kube_v1beta2_ReplicaSet_to_Koki_ReplicaSet(genericReplicaSet)
		if err != nil {
			return nil, err
		}

		kokiRS := &kokiWrapper.ReplicaSet

		kokiRS.Version = groupVersionString

		// Perform version-specific initialization here.

		return kokiWrapper, nil
	}

	return nil, util.InvalidInstanceErrorf(genericReplicaSet, "didn't deserialize 'generic' ReplicaSet as apps/v1beta2.ReplicaSet")
}

func Convert_Kube_v1beta2_ReplicaSet_to_Koki_ReplicaSet(kubeRS *appsv1beta2.ReplicaSet) (*types.ReplicaSetWrapper, error) {
	var err error
	kokiRS := &types.ReplicaSet{}

	kokiRS.Name = kubeRS.Name
	kokiRS.Namespace = kubeRS.Namespace
	kokiRS.Version = kubeRS.APIVersion
	kokiRS.Cluster = kubeRS.ClusterName
	kokiRS.Labels = kubeRS.Labels
	kokiRS.Annotations = kubeRS.Annotations

	kubeSpec := &kubeRS.Spec

	kokiRS.Replicas = kubeSpec.Replicas
	kokiRS.MinReadySeconds = kubeSpec.MinReadySeconds

	// Fill out the Selector and Template.Labels.
	// If kubeRS only has Template.Labels, we pull it up to Selector.
	var templateLabelsOverride map[string]string
	kokiRS.Selector, templateLabelsOverride, err = convertRSLabelSelector(kubeSpec.Selector, kubeSpec.Template.Labels)
	if err != nil {
		return nil, err
	}

	// Build a Pod from the kube Template. Use it to set the koki Template.
	kokiPod, err := convertRSTemplate(&kubeSpec.Template)
	if err != nil {
		return nil, err
	}
	kokiPod.Labels = templateLabelsOverride
	kokiRS.SetTemplate(kokiPod)

	if !reflect.DeepEqual(kubeRS.Status, appsv1beta2.ReplicaSetStatus{}) {
		kokiRS.Status = &kubeRS.Status
	}

	return &types.ReplicaSetWrapper{
		ReplicaSet: *kokiRS,
	}, nil
}

func convertRSTemplate(kubeTemplate *v1.PodTemplateSpec) (*types.Pod, error) {
	if kubeTemplate == nil {
		return nil, nil
	}

	kubePod := &v1.Pod{
		Spec: kubeTemplate.Spec,
	}

	kubePod.Name = kubeTemplate.Name
	kubePod.Namespace = kubeTemplate.Namespace
	kubePod.Labels = kubeTemplate.Labels
	kubePod.Annotations = kubeTemplate.Annotations

	kokiPod, err := Convert_Kube_v1_Pod_to_Koki_Pod(kubePod)
	if err != nil {
		return nil, err
	}

	return &kokiPod.Pod, nil
}

func convertRSLabelSelector(kubeSelector *metav1.LabelSelector, kubeTemplateLabels map[string]string) (*types.RSSelector, map[string]string, error) {
	// If the Selector is unspecified, it defaults to the Template's Labels.
	if kubeSelector == nil {
		return &types.RSSelector{
			Labels: kubeTemplateLabels,
		}, nil, nil
	}

	if len(kubeSelector.MatchExpressions) == 0 {
		// We have Labels for both Selector and Template.
		// If they're equal, we only need one.
		if reflect.DeepEqual(kubeSelector.MatchLabels, kubeTemplateLabels) {
			return &types.RSSelector{
				Labels: kubeSelector.MatchLabels,
			}, nil, nil
		}

		return &types.RSSelector{
			Labels: kubeSelector.MatchLabels,
		}, kubeTemplateLabels, nil
	}

	selectorString, err := expressions.UnparseLabelSelector(kubeSelector)
	if err != nil {
		return nil, nil, err
	}

	return &types.RSSelector{
		Shorthand: selectorString,
	}, kubeTemplateLabels, nil
}
