package converters

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	exts "k8s.io/api/extensions/v1beta1"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_DaemonSet_to_Kube_DaemonSet(daemonSet *types.DaemonSetWrapper) (interface{}, error) {
	// Perform version-agnostic conversion into apps/v1beta2 DaemonSet.
	kubeDaemonSet, err := Convert_Koki_DaemonSet_to_Kube_apps_v1beta2_DaemonSet(daemonSet)
	if err != nil {
		return nil, err
	}

	// Serialize the "generic" kube DaemonSet.
	b, err := yaml.Marshal(kubeDaemonSet)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, kubeDaemonSet, "couldn't serialize 'generic' kube DaemonSet")
	}

	// Deserialize a versioned kube DaemonSet using its apiVersion.
	versionedDaemonSet, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, err
	}

	switch versionedDaemonSet := versionedDaemonSet.(type) {
	case *appsv1beta2.DaemonSet:
		// Perform apps/v1beta2-specific initialization here.
	case *exts.DaemonSet:
		// Perform exts/v1beta1-specific initialization here.
	default:
		return nil, serrors.TypeErrorf(versionedDaemonSet, "deserialized the manifest, but not as a supported kube DaemonSet")
	}

	return versionedDaemonSet, nil
}

func Convert_Koki_DaemonSet_to_Kube_apps_v1beta2_DaemonSet(daemonSet *types.DaemonSetWrapper) (*appsv1beta2.DaemonSet, error) {
	var err error
	kubeDaemonSet := &appsv1beta2.DaemonSet{}
	kokiDaemonSet := &daemonSet.DaemonSet

	kubeDaemonSet.Name = kokiDaemonSet.Name
	kubeDaemonSet.Namespace = kokiDaemonSet.Namespace
	if len(kokiDaemonSet.Version) == 0 {
		kubeDaemonSet.APIVersion = "extensions/v1beta1"
	} else {
		kubeDaemonSet.APIVersion = kokiDaemonSet.Version
	}
	kubeDaemonSet.Kind = "DaemonSet"
	kubeDaemonSet.ClusterName = kokiDaemonSet.Cluster
	kubeDaemonSet.Labels = kokiDaemonSet.Labels
	kubeDaemonSet.Annotations = kokiDaemonSet.Annotations

	kubeSpec := &kubeDaemonSet.Spec

	// Setting the Selector and Template is identical to ReplicaSet
	// Get the right Selector and Template Labels.
	var templateLabelsOverride map[string]string
	var kokiTemplateLabels map[string]string
	if kokiDaemonSet.TemplateMetadata != nil {
		kokiTemplateLabels = kokiDaemonSet.TemplateMetadata.Labels
	}
	kubeSpec.Selector, templateLabelsOverride, err = revertRSSelector(kokiDaemonSet.Name, kokiDaemonSet.Selector, kokiTemplateLabels)
	if err != nil {
		return nil, err
	}
	// Set the right Labels before we fill in the Pod template with this metadata.
	kokiDaemonSet.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, kokiDaemonSet.TemplateMetadata)

	// Fill in the rest of the Pod template.
	kubeTemplate, err := revertTemplate(kokiDaemonSet.TemplateMetadata, kokiDaemonSet.PodTemplate)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}
	if kubeTemplate == nil {
		return nil, serrors.InvalidInstanceErrorf(kokiDaemonSet, "missing pod template")
	}
	kubeSpec.Template = *kubeTemplate

	// End Selector/Template section.

	kubeSpec.UpdateStrategy = revertDaemonSetStrategy(kokiDaemonSet)

	kubeSpec.MinReadySeconds = kokiDaemonSet.MinReadySeconds
	kubeSpec.RevisionHistoryLimit = kokiDaemonSet.RevisionHistoryLimit

	kubeDaemonSet.Status, err = revertDaemonSetStatus(kokiDaemonSet.DaemonSetStatus)
	if err != nil {
		return nil, err
	}

	return kubeDaemonSet, nil
}

func revertDaemonSetStatus(kokiStatus types.DaemonSetStatus) (appsv1beta2.DaemonSetStatus, error) {
	return appsv1beta2.DaemonSetStatus{
		ObservedGeneration:     kokiStatus.ObservedGeneration,
		CurrentNumberScheduled: kokiStatus.NumNodesScheduled,
		NumberMisscheduled:     kokiStatus.NumNodesMisscheduled,
		DesiredNumberScheduled: kokiStatus.NumNodesDesired,
		NumberReady:            kokiStatus.NumReady,
		UpdatedNumberScheduled: kokiStatus.NumUpdated,
		NumberAvailable:        kokiStatus.NumAvailable,
		NumberUnavailable:      kokiStatus.NumUnavailable,
		CollisionCount:         kokiStatus.CollisionCount,
	}, nil
}

func revertDaemonSetStrategy(kokiDaemonSet *types.DaemonSet) appsv1beta2.DaemonSetUpdateStrategy {
	if kokiDaemonSet.OnDelete {
		return appsv1beta2.DaemonSetUpdateStrategy{
			Type: appsv1beta2.OnDeleteDaemonSetStrategyType,
		}
	}

	var rollingUpdateConfig *appsv1beta2.RollingUpdateDaemonSet
	if kokiDaemonSet.MaxUnavailable != nil {
		rollingUpdateConfig = &appsv1beta2.RollingUpdateDaemonSet{
			MaxUnavailable: kokiDaemonSet.MaxUnavailable,
		}
	}

	return appsv1beta2.DaemonSetUpdateStrategy{
		Type:          appsv1beta2.RollingUpdateDaemonSetStrategyType,
		RollingUpdate: rollingUpdateConfig,
	}
}
