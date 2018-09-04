package converters

import (
	"reflect"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	exts "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/parser"
	"github.com/koki/short/parser/expressions"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_ReplicaSet_to_Kube_ReplicaSet(rs *types.ReplicaSetWrapper) (interface{}, error) {
	// Perform version-agnostic conversion into apps/v1beta2 ReplicaSet.
	kubeRS, err := Convert_Koki_ReplicaSet_to_Kube_v1beta2_ReplicaSet(rs)
	if err != nil {
		return nil, err
	}

	// Serialize the "generic" kube ReplicaSet.
	b, err := yaml.Marshal(kubeRS)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, kubeRS, "couldn't serialize 'generic' kube ReplicaSet")
	}

	// Deserialize a versioned kube ReplicaSet using its apiVersion.
	versionedReplicaSet, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, err
	}

	switch versionedReplicaSet := versionedReplicaSet.(type) {
	case *appsv1beta2.ReplicaSet:
		// Perform apps/v1beta2-specific initialization here.
	case *exts.ReplicaSet:
		// Perform exts/v1beta1-specific initialization here.
	default:
		return nil, serrors.TypeErrorf(versionedReplicaSet, "deserialized the manifest, but not as a supported kube ReplicaSet")
	}

	return versionedReplicaSet, nil
}

func Convert_Koki_ReplicaSet_to_Kube_v1beta2_ReplicaSet(rs *types.ReplicaSetWrapper) (*appsv1beta2.ReplicaSet, error) {
	var err error
	kubeRS := &appsv1beta2.ReplicaSet{}
	kokiRS := &rs.ReplicaSet

	kubeRS.Name = kokiRS.Name
	kubeRS.Namespace = kokiRS.Namespace
	if len(kokiRS.Version) == 0 {
		kubeRS.APIVersion = "extensions/v1beta1"
	} else {
		kubeRS.APIVersion = kokiRS.Version
	}
	kubeRS.Kind = "ReplicaSet"
	kubeRS.ClusterName = kokiRS.Cluster
	kubeRS.Labels = kokiRS.Labels
	kubeRS.Annotations = kokiRS.Annotations

	kubeSpec := &kubeRS.Spec
	kubeSpec.Replicas = kokiRS.Replicas
	kubeSpec.MinReadySeconds = kokiRS.MinReadySeconds

	// Get the right Selector and Template Labels.
	var templateLabelsOverride map[string]string
	var kokiTemplateLabels map[string]string
	if kokiRS.TemplateMetadata != nil {
		kokiTemplateLabels = kokiRS.TemplateMetadata.Labels
	}
	kubeSpec.Selector, templateLabelsOverride, err = revertRSSelector(kokiRS.Name, kokiRS.Selector, kokiTemplateLabels)
	if err != nil {
		return nil, err
	}
	// Set the right Labels before we fill in the Pod template with this metadata.
	kokiRS.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, kokiRS.TemplateMetadata)

	//  Fill in the rest of the Pod template.
	kubeTemplate, err := revertTemplate(kokiRS.TemplateMetadata, kokiRS.PodTemplate)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}
	if kubeTemplate == nil {
		return nil, serrors.InvalidInstanceErrorf(kokiRS, "missing pod template")
	}
	kubeSpec.Template = *kubeTemplate

	// End Selector/Template section.

	kubeRS.Status, err = revertReplicaSetStatus(kokiRS.ReplicaSetStatus)
	if err != nil {
		return nil, err
	}

	return kubeRS, nil
}

func revertReplicaSetStatus(kokiStatus types.ReplicaSetStatus) (appsv1beta2.ReplicaSetStatus, error) {
	conditions, err := revertReplicaSetConditions(kokiStatus.Conditions)
	if err != nil {
		return appsv1beta2.ReplicaSetStatus{}, err
	}
	return appsv1beta2.ReplicaSetStatus{
		ObservedGeneration:   kokiStatus.ObservedGeneration,
		Replicas:             kokiStatus.Replicas.Total,
		FullyLabeledReplicas: kokiStatus.Replicas.FullyLabeled,
		ReadyReplicas:        kokiStatus.Replicas.Ready,
		AvailableReplicas:    kokiStatus.Replicas.Available,
		Conditions:           conditions,
	}, nil
}

func revertReplicaSetConditions(kokiConditions []types.ReplicaSetCondition) ([]appsv1beta2.ReplicaSetCondition, error) {
	if len(kokiConditions) == 0 {
		return nil, nil
	}

	kubeConditions := make([]appsv1beta2.ReplicaSetCondition, len(kokiConditions))
	for i, condition := range kokiConditions {
		status, err := revertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "replica-set conditions[%d]", i)
		}
		conditionType, err := revertReplicaSetConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "replica-set conditions[%d]", i)
		}
		kubeConditions[i] = appsv1beta2.ReplicaSetCondition{
			Type:               conditionType,
			Status:             status,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}

	return kubeConditions, nil
}

func revertReplicaSetConditionType(kokiType types.ReplicaSetConditionType) (appsv1beta2.ReplicaSetConditionType, error) {
	switch kokiType {
	case types.ReplicaSetReplicaFailure:
		return appsv1beta2.ReplicaSetReplicaFailure, nil
	default:
		return appsv1beta2.ReplicaSetReplicaFailure, serrors.InvalidValueErrorf(kokiType, "unrecognized replica-set condition type")
	}
}

func applyTemplateLabelsOverride(labelsOverride map[string]string, kokiMeta *types.PodTemplateMeta) *types.PodTemplateMeta {
	if kokiMeta == nil {
		if len(labelsOverride) > 0 {
			return &types.PodTemplateMeta{
				Labels: labelsOverride,
			}
		}
		return nil
	} else {
		if len(labelsOverride) > 0 {
			kokiMeta.Labels = labelsOverride
		} else {
			kokiMeta.Labels = nil
		}

		if reflect.DeepEqual(kokiMeta, &types.PodTemplateMeta{}) {
			return nil
		}

		return kokiMeta
	}
}

func revertRSSelector(name string, selector *types.RSSelector, templateLabels map[string]string) (*metav1.LabelSelector, map[string]string, error) {
	if selector == nil {
		return nil, nil, nil
	}

	if len(selector.Shorthand) > 0 {
		labelSelector, err := expressions.ParseLabelSelector(selector.Shorthand)
		if err != nil {
			return nil, nil, serrors.InvalidInstanceErrorf(selector, "%s", err)
		}
		if len(templateLabels) == 0 && len(labelSelector.MatchExpressions) == 0 {
			// Selector is only Labels, and Template.Labels is empty.
			// Push the Selector Labels down into the Template.
			return labelSelector, labelSelector.MatchLabels, nil
		}

		// Template already has Labels specified OR Selector isn't just MatchLabels.
		// Can't copy the Selector into the Template Labels.
		return labelSelector, templateLabels, nil
	}

	if len(templateLabels) == 0 {
		// Copy the Selector Labels into the Template Labels.
		return &metav1.LabelSelector{
			MatchLabels: selector.Labels,
		}, selector.Labels, nil
	}

	// Template already has Labels specified.
	// Can't copy the Selector into the Template Labels.
	return &metav1.LabelSelector{
		MatchLabels: selector.Labels,
	}, templateLabels, nil
}
