package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
)

func Convert_Koki_Pod_to_Kube_v1_Pod(pod *types.PodWrapper) (*v1.Pod, error) {
	kubePod := &v1.Pod{}
	kokiPod := pod.Pod

	kubePod.Name = kokiPod.Name
	kubePod.Namespace = kokiPod.Namespace
	kubePod.APIVersion = kokiPod.Version
	kubePod.ClusterName = kokiPod.Cluster
	kubePod.Labels = kokiPod.Labels
	kubePod.Annotations = kokiPod.Annotations

	return kubePod, nil
}
