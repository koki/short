package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_LimitRange_to_Kube(wrapper *types.LimitRangeWrapper) (*v1.LimitRange, error) {
	var err error
	kube := &v1.LimitRange{}
	koki := wrapper.LimitRange

	kube.Name = koki.Name
	kube.Namespace = koki.Namespace
	if len(koki.Version) == 0 {
		kube.APIVersion = "v1"
	} else {
		kube.APIVersion = koki.Version
	}
	kube.Kind = "LimitRange"
	kube.ClusterName = koki.Cluster
	kube.Labels = koki.Labels
	kube.Annotations = koki.Annotations

	kube.Spec.Limits, err = revertLimitRangeItems(koki.Limits)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "LimitRange.Spec.Limits")
	}

	return kube, nil
}

func revertLimitRangeItems(kokis []types.LimitRangeItem) ([]v1.LimitRangeItem, error) {
	if len(kokis) == 0 {
		return nil, nil
	}

	var err error
	kubes := make([]v1.LimitRangeItem, len(kokis))
	for i, koki := range kokis {
		kubes[i], err = revertLimitRangeItem(koki)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "[%d]", i)
		}
	}

	return kubes, nil
}

func revertLimitRangeItem(koki types.LimitRangeItem) (v1.LimitRangeItem, error) {
	kube := v1.LimitRangeItem{
		Max:                  koki.Max,
		Min:                  koki.Min,
		Default:              koki.Default,
		DefaultRequest:       koki.DefaultRequest,
		MaxLimitRequestRatio: koki.MaxLimitRequestRatio,
	}

	var err error
	kube.Type, err = revertLimitType(koki.Type)
	return kube, err
}

func revertLimitType(koki types.LimitType) (v1.LimitType, error) {
	switch koki {
	case types.LimitTypePod:
		return v1.LimitTypePod, nil
	case types.LimitTypeContainer:
		return v1.LimitTypeContainer, nil
	case types.LimitTypePersistentVolumeClaim:
		return v1.LimitTypePersistentVolumeClaim, nil
	default:
		return "", serrors.InvalidInstanceError(koki)
	}
}
