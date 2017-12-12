package converters

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_PVC_to_Kube_PVC(pvc *types.PersistentVolumeClaimWrapper) (interface{}, error) {
	kubePVC := &v1.PersistentVolumeClaim{}
	kokiPVC := pvc.PersistentVolumeClaim

	kubePVC.Name = kokiPVC.Name
	kubePVC.Namespace = kokiPVC.Namespace
	kubePVC.APIVersion = kokiPVC.Version
	kubePVC.Kind = "PersistentVolumeClaim"
	kubePVC.ClusterName = kokiPVC.Cluster
	kubePVC.Labels = kokiPVC.Labels
	kubePVC.Annotations = kokiPVC.Annotations

	spec, err := revertPVCSpec(kokiPVC)
	if err != nil {
		return nil, err
	}
	kubePVC.Spec = spec

	status, err := revertPVCStatus(kokiPVC)
	if err != nil {
		return nil, err
	}
	kubePVC.Status = status

	return kubePVC, nil
}

func revertPVCSpec(kokiPVC types.PersistentVolumeClaim) (v1.PersistentVolumeClaimSpec, error) {
	kubeSpec := v1.PersistentVolumeClaimSpec{}

	selector, _, err := revertRSSelector(kokiPVC.Name, kokiPVC.Selector, nil)
	if err != nil {
		return kubeSpec, err
	}

	kubeSpec.Selector = selector

	accessModes, err := revertAccessModes(kokiPVC.AccessModes)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.AccessModes = accessModes

	resources, err := revertPVCResources(kokiPVC.Storage)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.Resources = resources

	kubeSpec.VolumeName = kokiPVC.Volume
	kubeSpec.StorageClassName = kokiPVC.StorageClass

	return kubeSpec, nil
}

func revertAccessModes(accessModes []types.PersistentVolumeAccessMode) ([]v1.PersistentVolumeAccessMode, error) {
	var kubeAccessModes []v1.PersistentVolumeAccessMode
	for i := range accessModes {
		mode := accessModes[i]
		kubeAccessMode := revertAccessMode(mode)
		kubeAccessModes = append(kubeAccessModes, kubeAccessMode)
	}
	return kubeAccessModes, nil
}

func revertAccessMode(accessMode types.PersistentVolumeAccessMode) v1.PersistentVolumeAccessMode {
	switch accessMode {
	case types.ReadWriteOnce:
		return v1.ReadWriteOnce
	case types.ReadOnlyMany:
		return v1.ReadOnlyMany
	case types.ReadWriteMany:
		return v1.ReadWriteMany
	default:
		return ""
	}
}

func revertPVCResources(storage string) (v1.ResourceRequirements, error) {
	if storage == "" {
		return v1.ResourceRequirements{}, nil
	}
	quantity, err := resource.ParseQuantity(storage)
	if err != nil {
		return v1.ResourceRequirements{}, err
	}
	return v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceStorage: quantity,
		},
	}, nil
}

func revertPVCStatus(kokiPVC types.PersistentVolumeClaim) (v1.PersistentVolumeClaimStatus, error) {
	kubePVCStatus := v1.PersistentVolumeClaimStatus{}
	kokiStatus := kokiPVC.PersistentVolumeClaimStatus

	conditions, err := revertPVCConditions(kokiStatus.Conditions)
	if err != nil {
		return kubePVCStatus, err
	}
	kubePVCStatus.Conditions = conditions

	capacity, err := revertPVCResources(kokiStatus.Storage)
	if err != nil {
		return kubePVCStatus, err
	}
	kubePVCStatus.Capacity = capacity.Requests

	accessModes, err := revertAccessModes(kokiStatus.AccessModes)
	if err != nil {
		return kubePVCStatus, err
	}
	kubePVCStatus.AccessModes = accessModes

	kubePVCStatus.Phase = revertPVCPhase(kokiStatus.Phase)

	return kubePVCStatus, nil
}

func revertPVCPhase(phase types.PersistentVolumeClaimPhase) v1.PersistentVolumeClaimPhase {
	switch phase {
	case types.ClaimPending:
		return v1.ClaimPending
	case types.ClaimBound:
		return v1.ClaimBound
	case types.ClaimLost:
		return v1.ClaimLost
	default:
		return ""
	}
}

func revertPVCConditions(kokiConditions []types.PersistentVolumeClaimCondition) ([]v1.PersistentVolumeClaimCondition, error) {
	if len(kokiConditions) == 0 {
		return nil, nil
	}

	kubeConditions := make([]v1.PersistentVolumeClaimCondition, len(kokiConditions))
	for i, condition := range kokiConditions {
		status, err := revertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "pvc conditions[%d]", i)
		}
		conditionType, err := revertPVCConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "pvc conditions[%d]", i)
		}
		kubeConditions[i] = v1.PersistentVolumeClaimCondition{
			Type:               conditionType,
			Status:             status,
			LastProbeTime:      condition.LastProbeTime,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}
	return kubeConditions, nil
}

func revertPVCConditionType(kokiType types.PersistentVolumeClaimConditionType) (v1.PersistentVolumeClaimConditionType, error) {
	switch kokiType {
	case types.PersistentVolumeClaimResizing:
		return v1.PersistentVolumeClaimResizing, nil
	default:
		return "", serrors.InvalidValueErrorf(kokiType, "unrecognized pvc condition type")
	}
}
