package converters

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_StatefulSet_to_Koki_StatefulSet(kubeStatefulSet runtime.Object) (*types.StatefulSetWrapper, error) {
	groupVersionKind := kubeStatefulSet.GetObjectKind().GroupVersionKind()
	groupVersionString := groupVersionKind.GroupVersion().String()
	groupVersionKind.Version = "v1beta2"
	groupVersionKind.Group = "apps"
	kubeStatefulSet.GetObjectKind().SetGroupVersionKind(groupVersionKind)

	// Serialize as v1beta2
	b, err := yaml.Marshal(kubeStatefulSet)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, kubeStatefulSet, "couldn't serialize kube StatefulSet after setting apiVersion to apps/v1beta2")
	}

	// Deserialize the "generic" kube StatefulSet
	genericStatefulSet, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, string(b), "couldn't deserialize 'generic' kube StatefulSet")
	}

	if genericStatefulSet, ok := genericStatefulSet.(*appsv1beta2.StatefulSet); ok {
		kokiWrapper, err := Convert_Kube_v1beta2_StatefulSet_to_Koki_StatefulSet(genericStatefulSet)
		if err != nil {
			return nil, err
		}

		kokiStatefulSet := &kokiWrapper.StatefulSet

		kokiStatefulSet.Version = groupVersionString

		return kokiWrapper, nil
	}

	return nil, serrors.InvalidInstanceErrorf(genericStatefulSet, "didn't deserialize 'generic' kube StatefulSet as apps/v1beta2.StatefulSet")
}

func Convert_Kube_v1beta2_StatefulSet_to_Koki_StatefulSet(kubeStatefulSet *appsv1beta2.StatefulSet) (*types.StatefulSetWrapper, error) {
	kokiStatefulSet := &types.StatefulSet{}

	kokiStatefulSet.Name = kubeStatefulSet.Name
	kokiStatefulSet.Namespace = kubeStatefulSet.Namespace
	kokiStatefulSet.Version = kubeStatefulSet.APIVersion
	kokiStatefulSet.Cluster = kubeStatefulSet.ClusterName
	kokiStatefulSet.Labels = kubeStatefulSet.Labels
	kokiStatefulSet.Annotations = kubeStatefulSet.Annotations

	kubeSpec := &kubeStatefulSet.Spec
	kokiStatefulSet.Replicas = kubeSpec.Replicas

	// Setting the Selector and Template is identical to ReplicaSet

	// Fill out the Selector and Template.Labels.
	// If kubeDeployment only has Template.Labels, we pull it up to Selector.
	selector, templateLabelsOverride, err := convertRSLabelSelector(kubeSpec.Selector, kubeSpec.Template.Labels)
	if err != nil {
		return nil, err
	}

	if selector != nil && (selector.Labels != nil || selector.Shorthand != "") {
		kokiStatefulSet.Selector = selector
	}

	// Build a Pod from the kube Template. Use it to set the koki Template.
	meta, template, err := convertTemplate(kubeSpec.Template)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}
	kokiStatefulSet.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, meta)
	kokiStatefulSet.PodTemplate = template

	// End Selector/Template section.

	kokiStatefulSet.OnDelete, kokiStatefulSet.Partition = convertStatefulSetStrategy(kubeSpec.UpdateStrategy)
	kokiStatefulSet.RevisionHistoryLimit = kubeSpec.RevisionHistoryLimit

	kokiStatefulSet.Service = kubeSpec.ServiceName
	pvcs, err := convertStatefulSetPVCs(kubeSpec.VolumeClaimTemplates)
	if err != nil {
		return nil, err
	}
	kokiStatefulSet.PVCs = pvcs

	podManagementPolicy, err := convertPodManagementPolicy(kubeSpec.PodManagementPolicy)
	if err != nil {
		return nil, err
	}
	kokiStatefulSet.PodManagementPolicy = podManagementPolicy

	kokiStatefulSet.StatefulSetStatus, err = convertStatefulSetStatus(kubeStatefulSet.Status)
	if err != nil {
		return nil, err
	}

	return &types.StatefulSetWrapper{
		StatefulSet: *kokiStatefulSet,
	}, nil
}

func convertPodManagementPolicy(podManagementPolicy appsv1beta2.PodManagementPolicyType) (types.PodManagementPolicyType, error) {
	if podManagementPolicy == "" {
		return "", nil
	}
	switch podManagementPolicy {
	case appsv1beta2.OrderedReadyPodManagement:
		return types.OrderedReadyPodManagement, nil
	case appsv1beta2.ParallelPodManagement:
		return types.ParallelPodManagement, nil
	default:
		return "", serrors.InvalidValueErrorf(podManagementPolicy, "unrecognized StatefulSet pod management policy type")
	}
}

func convertStatefulSetStrategy(strategy appsv1beta2.StatefulSetUpdateStrategy) (onDelete bool, partition *int32) {
	if strategy.RollingUpdate == nil && strategy.Type == appsv1beta2.OnDeleteStatefulSetStrategyType {
		return true, nil
	}
	if strategy.RollingUpdate != nil {
		return false, strategy.RollingUpdate.Partition
	}
	return false, nil
}

func convertStatefulSetPVCs(pvcs []v1.PersistentVolumeClaim) ([]types.PersistentVolumeClaim, error) {
	var kokiPVCs []types.PersistentVolumeClaim

	for i := range pvcs {
		pvc := pvcs[i]
		kokiPVC, err := Convert_Kube_PVC_to_Koki_PVC(&pvc)
		if err != nil {
			return nil, err
		}
		kokiPVCs = append(kokiPVCs, kokiPVC.PersistentVolumeClaim)
	}
	return kokiPVCs, nil
}

func convertStatefulSetStatus(kubeStatus appsv1beta2.StatefulSetStatus) (types.StatefulSetStatus, error) {
	return types.StatefulSetStatus{
		ObservedGeneration: kubeStatus.ObservedGeneration,
		Replicas:           kubeStatus.Replicas,
		ReadyReplicas:      kubeStatus.ReadyReplicas,
		CurrentReplicas:    kubeStatus.CurrentReplicas,
		UpdatedReplicas:    kubeStatus.UpdatedReplicas,
		Revision:           kubeStatus.CurrentRevision,
		UpdateRevision:     kubeStatus.UpdateRevision,
		CollisionCount:     kubeStatus.CollisionCount,
	}, nil
}
