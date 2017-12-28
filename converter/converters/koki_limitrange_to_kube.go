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

func revertLimitRangeItems(kokiItems []types.LimitRangeItem) ([]v1.LimitRangeItem, error) {
	if len(kokiItems) == 0 {
		return nil, nil
	}

	var err error
	kubeItems := make([]v1.LimitRangeItem, len(kokiItems))
	for i, kokiItem := range kokiItems {
		kubeItems[i], err = revertLimitRangeItem(kokiItem)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "[%d]", i)
		}
	}

	return kubeItems, nil
}

func revertLimitRangeItem(kokiItem types.LimitRangeItem) (v1.LimitRangeItem, error) {
	kubeItem := v1.LimitRangeItem{
		Max:                  kokiItem.Max,
		Min:                  kokiItem.Min,
		Default:              kokiItem.Default,
		DefaultRequest:       kokiItem.DefaultRequest,
		MaxLimitRequestRatio: kokiItem.MaxLimitRequestRatio,
	}

	var err error
	kubeItem.Type, err = revertLimitType(kokiItem.Type)
	return kubeItem, err
}

func revertLimitType(kokiType types.LimitType) (v1.LimitType, error) {
	if len(kokiType) == 0 {
		return "", nil
	}

	switch kokiType {
	case types.LimitTypePod:
		return v1.LimitTypePod, nil
	case types.LimitTypeContainer:
		return v1.LimitTypeContainer, nil
	case types.LimitTypePersistentVolumeClaim:
		return v1.LimitTypePersistentVolumeClaim, nil
	default:
		return "", serrors.InvalidInstanceError(kokiType)
	}
}
