package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_LimitRange_to_Koki(kube *v1.LimitRange) (*types.LimitRangeWrapper, error) {
	var err error
	koki := &types.LimitRange{}

	koki.Name = kube.Name
	koki.Namespace = kube.Namespace
	koki.Version = kube.APIVersion
	koki.Cluster = kube.ClusterName
	koki.Labels = kube.Labels
	koki.Annotations = kube.Annotations

	koki.Limits, err = convertLimitRangeItems(kube.Spec.Limits)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "limit_range limits")
	}

	return &types.LimitRangeWrapper{
		LimitRange: *koki,
	}, nil
}

func convertLimitRangeItems(kubeItems []v1.LimitRangeItem) ([]types.LimitRangeItem, error) {
	if len(kubeItems) == 0 {
		return nil, nil
	}

	var err error
	kokiItems := make([]types.LimitRangeItem, len(kubeItems))
	for i, kubeItem := range kubeItems {
		kokiItems[i], err = convertLimitRangeItem(kubeItem)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "[%d]", i)
		}
	}

	return kokiItems, nil
}

func convertLimitRangeItem(kubeItem v1.LimitRangeItem) (types.LimitRangeItem, error) {
	kokiItem := types.LimitRangeItem{
		Max:                  kubeItem.Max,
		Min:                  kubeItem.Min,
		Default:              kubeItem.Default,
		DefaultRequest:       kubeItem.DefaultRequest,
		MaxLimitRequestRatio: kubeItem.MaxLimitRequestRatio,
	}

	var err error
	kokiItem.Type, err = convertLimitType(kubeItem.Type)
	return kokiItem, err
}

func convertLimitType(kubeItem v1.LimitType) (types.LimitType, error) {
	if len(kubeItem) == 0 {
		return "", nil
	}

	switch kubeItem {
	case v1.LimitTypePod:
		return types.LimitTypePod, nil
	case v1.LimitTypeContainer:
		return types.LimitTypeContainer, nil
	case v1.LimitTypePersistentVolumeClaim:
		return types.LimitTypePersistentVolumeClaim, nil
	default:
		return "", serrors.InvalidInstanceError(kubeItem)
	}
}
