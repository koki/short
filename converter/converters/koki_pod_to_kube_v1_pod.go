package converters

import (
	"fmt"
	"net/url"
	"strings"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_Pod_to_Kube_v1_Pod(pod *types.PodWrapper) (*v1.Pod, error) {
	kubePod := &v1.Pod{}
	kokiPod := pod.Pod

	kubePod.Name = kokiPod.Name
	kubePod.Namespace = kokiPod.Namespace
	kubePod.APIVersion = kokiPod.Version
	kubePod.ClusterName = kokiPod.Cluster
	kubePod.Labels = kokiPod.Labels
	kubePod.Annotation = kokiPod.Annotation

	return kubePod, nil
}
