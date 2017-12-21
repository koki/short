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

func convertLimitRangeItems(kubes []v1.LimitRangeItem) ([]types.LimitRangeItem, error) {
	if len(kubes) == 0 {
		return nil, nil
	}

	var err error
	kokis := make([]types.LimitRangeItem, len(kubes))
	for i, kube := range kubes {
		kokis[i], err = convertLimitRangeItem(kube)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "[%d]", i)
		}
	}

	return kokis, nil
}

func convertLimitRangeItem(kube v1.LimitRangeItem) (types.LimitRangeItem, error) {
	koki := types.LimitRangeItem{
		Max:                  kube.Max,
		Min:                  kube.Min,
		Default:              kube.Default,
		DefaultRequest:       kube.DefaultRequest,
		MaxLimitRequestRatio: kube.MaxLimitRequestRatio,
	}

	var err error
	koki.Type, err = convertLimitType(kube.Type)
	return koki, err
}

func convertLimitType(kube v1.LimitType) (types.LimitType, error) {
	switch kube {
	case v1.LimitTypePod:
		return types.LimitTypePod, nil
	case v1.LimitTypeContainer:
		return types.LimitTypeContainer, nil
	case v1.LimitTypePersistentVolumeClaim:
		return types.LimitTypePersistentVolumeClaim, nil
	default:
		return "", serrors.InvalidInstanceError(kube)
	}
}
