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

func Convert_Kube_DaemonSet_to_Koki_DaemonSet(kubeDaemonSet runtime.Object) (*types.DaemonSetWrapper, error) {
	groupVersionKind := kubeDaemonSet.GetObjectKind().GroupVersionKind()
	groupVersionString := groupVersionKind.GroupVersion().String()
	groupVersionKind.Version = "v1beta2"
	groupVersionKind.Group = "apps"
	kubeDaemonSet.GetObjectKind().SetGroupVersionKind(groupVersionKind)

	// Serialize as v1beta2
	b, err := yaml.Marshal(kubeDaemonSet)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, kubeDaemonSet, "couldn't serialize kube DaemonSet after setting apiVersion to apps/v1beta2")
	}

	// Deserialize the "generic" kube DaemonSet
	genericDaemonSet, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, string(b), "couldn't deserialize 'generic' kube DaemonSet")
	}

	if genericDaemonSet, ok := genericDaemonSet.(*appsv1beta2.DaemonSet); ok {
		kokiWrapper, err := Convert_Kube_v1beta2_DaemonSet_to_Koki_DaemonSet(genericDaemonSet)
		if err != nil {
			return nil, err
		}

		kokiDaemonSet := &kokiWrapper.DaemonSet

		kokiDaemonSet.Version = groupVersionString

		// Perform version-specific initialization here.

		return kokiWrapper, nil
	}

	return nil, serrors.InvalidInstanceErrorf(genericDaemonSet, "didn't deserialize 'generic' kube DaemonSet as apps/v1beta2.DaemonSet")
}

func Convert_Kube_v1beta2_DaemonSet_to_Koki_DaemonSet(kubeDaemonSet *appsv1beta2.DaemonSet) (*types.DaemonSetWrapper, error) {
	kokiDaemonSet := &types.DaemonSet{}

	kokiDaemonSet.Name = kubeDaemonSet.Name
	kokiDaemonSet.Namespace = kubeDaemonSet.Namespace
	kokiDaemonSet.Version = kubeDaemonSet.APIVersion
	kokiDaemonSet.Cluster = kubeDaemonSet.ClusterName
	kokiDaemonSet.Labels = kubeDaemonSet.Labels
	kokiDaemonSet.Annotations = kubeDaemonSet.Annotations

	kubeSpec := &kubeDaemonSet.Spec

	// Setting the Selector and Template is identical to ReplicaSet

	// Fill out the Selector and Template.Labels.
	// If kubeDaemonSet only has Template.Labels, we pull it up to Selector.
	selector, templateLabelsOverride, err := convertRSLabelSelector(kubeSpec.Selector, kubeSpec.Template.Labels)
	if err != nil {
		return nil, err
	}

	if selector != nil && (selector.Labels != nil || selector.Shorthand != "") {
		kokiDaemonSet.Selector = selector
	}

	// Build a Pod from the kube Template. Use it to set the koki Template.
	meta, template, err := convertTemplate(kubeSpec.Template)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}
	kokiDaemonSet.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, meta)
	kokiDaemonSet.PodTemplate = template

	// End Selector/Template section.

	kokiDaemonSet.OnDelete, kokiDaemonSet.MaxUnavailable = convertDaemonSetStrategy(kubeSpec.UpdateStrategy)

	kokiDaemonSet.MinReadySeconds = kubeSpec.MinReadySeconds
	kokiDaemonSet.RevisionHistoryLimit = kubeSpec.RevisionHistoryLimit

	kokiDaemonSet.DaemonSetStatus, err = convertDaemonSetStatus(kubeDaemonSet.Status)
	if err != nil {
		return nil, err
	}

	return &types.DaemonSetWrapper{
		DaemonSet: *kokiDaemonSet,
	}, nil
}

func convertDaemonSetStatus(kubeStatus appsv1beta2.DaemonSetStatus) (types.DaemonSetStatus, error) {
	return types.DaemonSetStatus{
		ObservedGeneration:   kubeStatus.ObservedGeneration,
		NumNodesScheduled:    kubeStatus.CurrentNumberScheduled,
		NumNodesMisscheduled: kubeStatus.NumberMisscheduled,
		NumNodesDesired:      kubeStatus.DesiredNumberScheduled,
		NumReady:             kubeStatus.NumberReady,
		NumUpdated:           kubeStatus.UpdatedNumberScheduled,
		NumAvailable:         kubeStatus.NumberAvailable,
		NumUnavailable:       kubeStatus.NumberUnavailable,
		CollisionCount:       kubeStatus.CollisionCount,
	}, nil
}

func convertDaemonSetStrategy(kubeStrategy appsv1beta2.DaemonSetUpdateStrategy) (onDelete bool, maxUnavailable *intstr.IntOrString) {
	if kubeStrategy.Type == appsv1beta2.OnDeleteDaemonSetStrategyType {
		return true, nil
	}

	if rollingUpdate := kubeStrategy.RollingUpdate; rollingUpdate != nil {
		return false, rollingUpdate.MaxUnavailable
	}

	return false, nil
}
