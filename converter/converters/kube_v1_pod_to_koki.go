package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_v1_Pod_to_Koki_Pod(pod *v1.Pod) (*types.PodWrapper, error) {
	return nil, nil
}
