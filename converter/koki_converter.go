package converter

import (
	"reflect"

	"github.com/koki/short/converter/converters"
	"github.com/koki/short/types"
	"github.com/koki/short/util"

	apps "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func detectAndConvertFromKokiObj(kokiObj interface{}) (interface{}, error) {
	switch kokiObj := kokiObj.(type) {
	case *types.PodWrapper:
		return converters.Convert_Koki_Pod_to_Kube_v1_Pod(kokiObj)
	case *types.ReplicationControllerWrapper:
		return converters.Convert_Koki_ReplicationController_to_Kube_v1_ReplicationController(kokiObj)
	case *types.ReplicaSetWrapper:
		return converters.Convert_Koki_ReplicaSet_to_Kube_v1_ReplicaSet(kokiObj)
	default:
		return nil, util.TypeErrorf(reflect.TypeOf(kokiObj), "Unsupported Type")
	}
}

func detectAndConvertFromKubeObj(kubeObj runtime.Object) (interface{}, error) {
	switch kubeObj := kubeObj.(type) {
	case *v1.Pod:
		return converters.Convert_Kube_v1_Pod_to_Koki_Pod(kubeObj)
	case *v1.Service:
		return converters.Convert_Kube_v1_Service_to_Koki_Service(kubeObj)
	case *v1.ReplicationController:
		return converters.Convert_Kube_v1_ReplicationController_to_Koki_ReplicationController(kubeObj)
	case *apps.ReplicaSet:
		return converters.Convert_Kube_v1_ReplicaSet_to_Koki_ReplicaSet(kubeObj)
	default:
		return nil, util.TypeErrorf(reflect.TypeOf(kubeObj), "Unsupported Type")
	}
}
