package converters

import (
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_StatefulSet_to_Kube_StatefulSet(statefulSet *types.StatefulSetWrapper) (interface{}, error) {
	// Perform version-agnostic conversion into apps/v1beta2 StatefulSet.
	kubeStatefulSet, err := Convert_Koki_StatefulSet_to_Kube_apps_v1beta2_StatefulSet(statefulSet)
	if err != nil {
		return nil, err
	}

	// Serialize the "generic" kube StatefulSet.
	b, err := yaml.Marshal(kubeStatefulSet)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, kubeStatefulSet, "couldn't serialize 'generic' kube StatefulSet")
	}

	// Deserialize a versioned kube StatefulSet using its apiVersion.
	versionedStatefulSet, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, err
	}

	switch versionedStatefulSet := versionedStatefulSet.(type) {
	case *appsv1beta1.StatefulSet:
		// Perform apps/v1beta1-specific initialization here.
	case *appsv1beta2.StatefulSet:
		// Perform apps/v1beta2-specific initialization here.
	default:
		return nil, serrors.TypeErrorf(versionedStatefulSet, "deserialized the manifest, but not as a supported kube StatefulSet")
	}

	return versionedStatefulSet, nil
}

func Convert_Koki_StatefulSet_to_Kube_apps_v1beta2_StatefulSet(statefulSet *types.StatefulSetWrapper) (*appsv1beta2.StatefulSet, error) {
	var err error
	kubeStatefulSet := &appsv1beta2.StatefulSet{}
	kokiStatefulSet := &statefulSet.StatefulSet

	kubeStatefulSet.Name = kokiStatefulSet.Name
	kubeStatefulSet.Namespace = kokiStatefulSet.Namespace
	kubeStatefulSet.APIVersion = kokiStatefulSet.Version
	kubeStatefulSet.Kind = "StatefulSet"
	kubeStatefulSet.ClusterName = kokiStatefulSet.Cluster
	kubeStatefulSet.Labels = kokiStatefulSet.Labels
	kubeStatefulSet.Annotations = kokiStatefulSet.Annotations

	kubeSpec := &kubeStatefulSet.Spec
	kubeSpec.Replicas = kokiStatefulSet.Replicas

	// Setting the Selector and Template is identical to ReplicaSet
	// Get the right Selector and Template Labels.
	var templateLabelsOverride map[string]string
	var kokiTemplateLabels map[string]string
	if kokiStatefulSet.TemplateMetadata != nil {
		kokiTemplateLabels = kokiStatefulSet.TemplateMetadata.Labels
	}
	kubeSpec.Selector, templateLabelsOverride, err = revertRSSelector(kokiStatefulSet.Name, kokiStatefulSet.Selector, kokiTemplateLabels)
	if err != nil {
		return nil, err
	}
	// Set the right Labels before we fill in the Pod template with this metadata.
	kokiStatefulSet.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, kokiStatefulSet.TemplateMetadata)

	// Fill in the rest of the Pod template.
	kubeTemplate, err := revertTemplate(kokiStatefulSet.TemplateMetadata, kokiStatefulSet.PodTemplate)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}
	if kubeTemplate == nil {
		return nil, serrors.InvalidInstanceErrorf(kokiStatefulSet, "missing pod template")
	}
	kubeSpec.Template = *kubeTemplate

	// End Selector/Template section.

	kubeSpec.UpdateStrategy = revertStatefulSetStrategy(*kokiStatefulSet)
	kubeSpec.PodManagementPolicy = revertPodManagementPolicy(kokiStatefulSet.PodManagementPolicy)

	pvcs, err := revertPVCs(kokiStatefulSet.PVCs)
	if err != nil {
		return nil, err
	}
	kubeSpec.VolumeClaimTemplates = pvcs
	kubeSpec.ServiceName = kokiStatefulSet.Service
	kubeSpec.RevisionHistoryLimit = kokiStatefulSet.RevisionHistoryLimit

	kubeStatefulSet.Status = revertStatefulSetStatus(kokiStatefulSet.StatefulSetStatus)

	return kubeStatefulSet, nil
}

func revertPVCs(kokiPVCs []types.PersistentVolumeClaim) ([]v1.PersistentVolumeClaim, error) {
	var pvcs []v1.PersistentVolumeClaim

	for i := range kokiPVCs {
		kokiPVC := kokiPVCs[i]
		pvcWrapper := &types.PersistentVolumeClaimWrapper{
			PersistentVolumeClaim: kokiPVC,
		}
		pvc, err := Convert_Koki_PVC_to_Kube_PVC(pvcWrapper)
		if err != nil {
			return nil, err
		}
		pvcs = append(pvcs, *pvc.(*v1.PersistentVolumeClaim))
	}

	return pvcs, nil
}

func revertStatefulSetStrategy(statefulSet types.StatefulSet) appsv1beta2.StatefulSetUpdateStrategy {
	if statefulSet.OnDelete == true {
		return appsv1beta2.StatefulSetUpdateStrategy{
			Type: appsv1beta2.OnDeleteStatefulSetStrategyType,
		}
	}

	return appsv1beta2.StatefulSetUpdateStrategy{
		Type: appsv1beta2.RollingUpdateStatefulSetStrategyType,
		RollingUpdate: &appsv1beta2.RollingUpdateStatefulSetStrategy{
			Partition: statefulSet.Partition,
		},
	}
}

func revertPodManagementPolicy(podManagementPolicy types.PodManagementPolicyType) appsv1beta2.PodManagementPolicyType {
	if podManagementPolicy == "" {
		return ""
	}
	switch podManagementPolicy {
	case types.OrderedReadyPodManagement:
		return appsv1beta2.OrderedReadyPodManagement
	case types.ParallelPodManagement:
		return appsv1beta2.ParallelPodManagement
	default:
		return ""
	}
}

func revertStatefulSetStatus(kokiStatus types.StatefulSetStatus) appsv1beta2.StatefulSetStatus {
	return appsv1beta2.StatefulSetStatus{
		ObservedGeneration: kokiStatus.ObservedGeneration,
		Replicas:           kokiStatus.Replicas,
		ReadyReplicas:      kokiStatus.ReadyReplicas,
		CurrentReplicas:    kokiStatus.CurrentReplicas,
		UpdatedReplicas:    kokiStatus.UpdatedReplicas,
		CurrentRevision:    kokiStatus.Revision,
		UpdateRevision:     kokiStatus.UpdateRevision,
		CollisionCount:     kokiStatus.CollisionCount,
	}
}
