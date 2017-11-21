package converters

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	exts "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ghodss/yaml"

	"github.com/koki/short/parser"
	"github.com/koki/short/parser/expressions"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
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
		return nil, util.InvalidValueErrorf(kubeRS, "couldn't serialize 'generic' kube ReplicaSet: %s", err.Error())
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
		return nil, util.TypeErrorf(versionedReplicaSet, "deserialized the manifest, but not as a supported kube ReplicaSet")
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

	var templateLabelsOverride map[string]string
	var kokiTemplateLabels map[string]string
	if kokiRS.TemplateMetadata != nil {
		kokiTemplateLabels = kokiRS.TemplateMetadata.Labels
	}
	kubeSpec.Selector, templateLabelsOverride, err = revertRSSelector(kokiRS.Name, kokiRS.Selector, kokiTemplateLabels)
	if err != nil {
		return nil, err
	}

	kubeTemplate, err := revertTemplate(kokiRS.GetTemplate())
	if err != nil {
		return nil, err
	}
	if kubeTemplate == nil {
		return nil, util.InvalidInstanceErrorf(kokiRS, "missing pod template")
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

func revertRSSelector(name string, selector *types.RSSelector, templateLabels map[string]string) (*metav1.LabelSelector, map[string]string, error) {
	if selector == nil {
		defaultSelector := map[string]string{
			"koki.io/selector.name": name,
		}
		return &metav1.LabelSelector{
			MatchLabels: defaultSelector,
		}, defaultSelector, nil
	}

	if len(selector.Shorthand) > 0 {
		labelSelector, err := expressions.ParseLabelSelector(selector.Shorthand)
		if err != nil {
			return nil, nil, util.InvalidInstanceErrorf(selector, "%s", err)
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
