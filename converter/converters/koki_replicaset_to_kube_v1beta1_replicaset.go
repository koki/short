package converters

import (
	"k8s.io/api/core/v1"
	exts "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/parser/expressions"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_ReplicaSet_to_Kube_v1beta1_ReplicaSet(rs *types.ReplicaSetWrapper) (*exts.ReplicaSet, error) {
	var err error
	kubeRS := &exts.ReplicaSet{}
	kokiRS := &rs.ReplicaSet

	kubeRS.Name = kokiRS.Name
	kubeRS.Namespace = kokiRS.Namespace
	kubeRS.APIVersion = kokiRS.Version
	kubeRS.Kind = "ReplicaSet"
	kubeRS.ClusterName = kokiRS.Cluster
	kubeRS.Labels = kokiRS.Labels
	kubeRS.Annotations = kokiRS.Annotations

	kubeSpec := &kubeRS.Spec
	kubeSpec.Replicas = kokiRS.Replicas
	kubeSpec.MinReadySeconds = kokiRS.MinReadySeconds

	var templateLabelsOverride map[string]string
	var kokiTemplateLabels map[string]string
	if kokiRS.TemplateMetadata != nil {
		kokiTemplateLabels = kokiRS.TemplateMetadata.Labels
	}
	kubeSpec.Selector, templateLabelsOverride, err = revertRSSelector(kokiRS.Selector, kokiTemplateLabels)
	if err != nil {
		return nil, err
	}

	kubeTemplate, err := revertRSTemplate(kokiRS)
	if err != nil {
		return nil, err
	}
	if kubeTemplate == nil {
		return nil, util.TypeValueErrorf(kokiRS, "missing pod template")
	}
	kubeTemplate.Labels = templateLabelsOverride
	kubeSpec.Template = *kubeTemplate

	// Make sure the Selector and the Template.Labels are set correctly.
	if len(kubeSpec.Template.Labels) == 0 {
	}

	if kokiRS.Status != nil {
		kubeRS.Status = *kokiRS.Status
	}

	return kubeRS, nil
}

func revertRSSelector(selector *types.RSSelector, templateLabels map[string]string) (*metav1.LabelSelector, map[string]string, error) {
	if selector == nil {
		return nil, nil, util.PrettyTypeError(selector, "Selector is required for ReplicaSet")
	}

	if len(selector.Shorthand) > 0 {
		labelSelector, err := expressions.ParseLabelSelector(selector.Shorthand)
		if err != nil {
			return nil, nil, err
		}
		if len(templateLabels) == 0 && len(labelSelector.MatchExpressions) == 0 {
			// Selector is only Labels, and Template.Labels is empty.
			// Push the Selector Labels down into the Template.
			return nil, labelSelector.MatchLabels, nil
		}

		// Can't push the Selector down into the Template.
		return labelSelector, templateLabels, nil
	}

	if len(templateLabels) == 0 {
		// Push the Selector Labels down into the Template.
		return nil, selector.Labels, nil
	}

	// Can't push the Selector down into the Template.
	return &metav1.LabelSelector{
		MatchLabels: selector.Labels,
	}, templateLabels, nil
}

func revertRSTemplate(kokiRS *types.ReplicaSet) (*v1.PodTemplateSpec, error) {
	return revertTemplate(kokiRS.GetTemplate())
}
