package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_PVC_to_Koki_PVC(kubePVC *v1.PersistentVolumeClaim) (*types.PersistentVolumeClaimWrapper, error) {
	kokiObj := &types.PersistentVolumeClaimWrapper{}
	kokiPVC := &types.PersistentVolumeClaim{}

	kokiPVC.Version = kubePVC.APIVersion
	kokiPVC.Name = kubePVC.Name
	kokiPVC.Namespace = kubePVC.Namespace
	kokiPVC.Cluster = kubePVC.ClusterName
	kokiPVC.Labels = kubePVC.Labels
	kokiPVC.Annotations = kubePVC.Annotations

	err := convertPVCSpec(kubePVC.Spec, kokiPVC)
	if err != nil {
		return nil, err
	}

	status, err := convertPVCStatus(kubePVC.Status)
	if err != nil {
		return nil, err
	}
	kokiPVC.PersistentVolumeClaimStatus = status

	kokiObj.PersistentVolumeClaim = *kokiPVC

	return kokiObj, nil
}

func convertPVCSpec(kubeSpec v1.PersistentVolumeClaimSpec, kokiPVC *types.PersistentVolumeClaim) error {
	selector, _, err := convertRSLabelSelector(kubeSpec.Selector, nil)
	if err != nil {
		return err
	}

	if selector != nil && (selector.Labels != nil || selector.Shorthand != "") {
		kokiPVC.Selector = selector
	}

	accessModes, err := convertAccessModes(kubeSpec.AccessModes)
	if err != nil {
		return err
	}
	kokiPVC.AccessModes = accessModes

	kokiPVC.StorageClass = kubeSpec.StorageClassName
	kokiPVC.Volume = kubeSpec.VolumeName
	kokiPVC.Storage = convertStorageRequirement(kubeSpec.Resources.Requests)

	return nil
}

func convertAccessModes(accessModes []v1.PersistentVolumeAccessMode) ([]types.PersistentVolumeAccessMode, error) {
	var kokiAccessModes []types.PersistentVolumeAccessMode

	for i := range accessModes {
		mode := accessModes[i]

		kokiMode, err := convertAccessMode(mode)
		if err != nil {
			return nil, err
		}

		kokiAccessModes = append(kokiAccessModes, kokiMode)
	}

	return kokiAccessModes, nil
}

func convertAccessMode(accessMode v1.PersistentVolumeAccessMode) (types.PersistentVolumeAccessMode, error) {
	switch accessMode {
	case v1.ReadWriteOnce:
		return types.ReadWriteOnce, nil
	case v1.ReadOnlyMany:
		return types.ReadOnlyMany, nil
	case v1.ReadWriteMany:
		return types.ReadWriteMany, nil
	default:
		return "", serrors.InvalidValueErrorf(accessMode, "unrecognized access mode")
	}
}

func convertStorageRequirement(resource v1.ResourceList) string {
	if resource == nil {
		return ""
	}

	storage, ok := resource[v1.ResourceStorage]
	if !ok {
		return ""
	}

	return storage.String()
}

func convertPVCStatus(status v1.PersistentVolumeClaimStatus) (types.PersistentVolumeClaimStatus, error) {
	kokiStatus := types.PersistentVolumeClaimStatus{}

	accessModes, err := convertAccessModes(status.AccessModes)
	if err != nil {
		return kokiStatus, err
	}
	kokiStatus.AccessModes = accessModes
	kokiStatus.Storage = convertStorageRequirement(status.Capacity)

	phase, err := convertPVCPhase(status.Phase)
	if err != nil {
		return kokiStatus, err
	}
	kokiStatus.Phase = phase

	conditions, err := convertPVCConditions(status.Conditions)
	if err != nil {
		return kokiStatus, err
	}
	kokiStatus.Conditions = conditions

	return kokiStatus, nil
}

func convertPVCPhase(phase v1.PersistentVolumeClaimPhase) (types.PersistentVolumeClaimPhase, error) {
	if phase == "" {
		return "", nil
	}
	switch phase {
	case v1.ClaimPending:
		return types.ClaimPending, nil
	case v1.ClaimBound:
		return types.ClaimBound, nil
	case v1.ClaimLost:
		return types.ClaimLost, nil
	default:
		return "", serrors.InvalidValueErrorf(phase, "unrecognized volume claim phase")
	}
}

func convertPVCConditions(kubeConditions []v1.PersistentVolumeClaimCondition) ([]types.PersistentVolumeClaimCondition, error) {
	if len(kubeConditions) == 0 {
		return nil, nil
	}

	kokiConditions := make([]types.PersistentVolumeClaimCondition, len(kubeConditions))
	for i, condition := range kubeConditions {
		status, err := convertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "pvc conditions[%d]", i)
		}
		conditionType, err := convertPVCConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "job conditions[%d]", i)
		}
		kokiConditions[i] = types.PersistentVolumeClaimCondition{
			Type:               conditionType,
			Status:             status,
			LastProbeTime:      condition.LastProbeTime,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}

	return kokiConditions, nil
}

func convertPVCConditionType(kubeType v1.PersistentVolumeClaimConditionType) (types.PersistentVolumeClaimConditionType, error) {
	switch kubeType {
	case v1.PersistentVolumeClaimResizing:
		return types.PersistentVolumeClaimResizing, nil
	default:
		return "", serrors.InvalidValueErrorf(kubeType, "unrecognized pvc condition type")
	}
}
