package converters

import (
	"reflect"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/koki/short/parser"
	"github.com/koki/short/parser/expressions"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
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
		return nil, serrors.InvalidInstanceContextErrorf(err, kubeRS, "couldn't serialize kube ReplicaSet after setting apiVersion to apps/v1beta2")
	}

	// Deserialize the "generic" kube ReplicaSet
	genericReplicaSet, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, string(b), "couldn't deserialize 'generic' kube ReplicaSet")
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

	return nil, serrors.InvalidInstanceErrorf(genericReplicaSet, "didn't deserialize 'generic' ReplicaSet as apps/v1beta2.ReplicaSet")
}

func Convert_Kube_v1beta2_ReplicaSet_to_Koki_ReplicaSet(kubeRS *appsv1beta2.ReplicaSet) (*types.ReplicaSetWrapper, error) {
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
	selector, templateLabelsOverride, err := convertRSLabelSelector(kubeSpec.Selector, kubeSpec.Template.Labels)
	if err != nil {
		return nil, err
	}

	if selector != nil && (selector.Labels != nil || selector.Shorthand != "") {
		kokiRS.Selector = selector
	}

	// Build a Pod from the kube Template. Use it to set the koki Template.
	meta, template, err := convertTemplate(kubeSpec.Template)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}
	kokiRS.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, meta)
	kokiRS.PodTemplate = template

	// End Selector/Template section.

	kokiRS.ReplicaSetStatus, err = convertReplicaSetStatus(kubeRS.Status)
	if err != nil {
		return nil, err
	}

	return &types.ReplicaSetWrapper{
		ReplicaSet: *kokiRS,
	}, nil
}

func convertReplicaSetStatus(kubeStatus appsv1beta2.ReplicaSetStatus) (types.ReplicaSetStatus, error) {
	conditions, err := convertReplicaSetConditions(kubeStatus.Conditions)
	if err != nil {
		return types.ReplicaSetStatus{}, err
	}
	return types.ReplicaSetStatus{
		ObservedGeneration: kubeStatus.ObservedGeneration,
		Replicas: types.ReplicaSetReplicasStatus{
			Total:        kubeStatus.Replicas,
			FullyLabeled: kubeStatus.FullyLabeledReplicas,
			Ready:        kubeStatus.ReadyReplicas,
			Available:    kubeStatus.AvailableReplicas,
		},
		Conditions: conditions,
	}, nil
}

func convertReplicaSetConditions(kubeConditions []appsv1beta2.ReplicaSetCondition) ([]types.ReplicaSetCondition, error) {
	if len(kubeConditions) == 0 {
		return nil, nil
	}

	kokiConditions := make([]types.ReplicaSetCondition, len(kubeConditions))
	for i, condition := range kubeConditions {
		status, err := convertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "replica-set conditions[%d]", i)
		}
		conditionType, err := convertReplicaSetConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "replica-set conditions[%d]", i)
		}
		kokiConditions[i] = types.ReplicaSetCondition{
			Type:               conditionType,
			Status:             status,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}

	return kokiConditions, nil
}

func convertReplicaSetConditionType(kubeType appsv1beta2.ReplicaSetConditionType) (types.ReplicaSetConditionType, error) {
	switch kubeType {
	case appsv1beta2.ReplicaSetReplicaFailure:
		return types.ReplicaSetReplicaFailure, nil
	default:
		return types.ReplicaSetReplicaFailure, serrors.InvalidValueErrorf(kubeType, "unrecognized replica-set condition type")
	}
}

func convertRSLabelSelector(kubeSelector *metav1.LabelSelector, kubeTemplateLabels map[string]string) (*types.RSSelector, map[string]string, error) {
	// If the Selector is unspecified, it defaults to the Template's Labels.
	if kubeSelector == nil {
		return &types.RSSelector{
			Labels: kubeTemplateLabels,
		}, nil, nil
	}

	if len(kubeSelector.MatchExpressions) == 0 {
		if reflect.DeepEqual(kubeSelector.MatchLabels, kubeTemplateLabels) {
			// Selector and template labels are identical. Just keep the selector.
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
