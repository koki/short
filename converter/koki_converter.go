package converter

import (
	"reflect"

	"github.com/koki/short/converter/converters"
	"github.com/koki/short/types"
	"github.com/koki/short/util"

	"k8s.io/api/core/v1"
	exts "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func detectAndConvertFromKokiObj(kokiObj interface{}) (interface{}, error) {
	switch kokiObj := kokiObj.(type) {
	case *types.DeploymentWrapper:
		return converters.Convert_Koki_Deployment_to_Kube_v1beta1_Deployment(kokiObj)
	case *types.PersistentVolumeWrapper:
		return converters.Convert_Koki_PersistentVolume_to_Kube_v1_PersistentVolume(kokiObj)
	case *types.PodWrapper:
		return converters.Convert_Koki_Pod_to_Kube_v1_Pod(kokiObj)
	case *types.ReplicationControllerWrapper:
		return converters.Convert_Koki_ReplicationController_to_Kube_v1_ReplicationController(kokiObj)
	case *types.ReplicaSetWrapper:
		return converters.Convert_Koki_ReplicaSet_to_Kube_v1beta2_ReplicaSet(kokiObj)
	case *types.ServiceWrapper:
		return converters.Convert_Koki_Service_To_Kube_v1_Service(kokiObj)
	default:
		return nil, util.TypeErrorf(reflect.TypeOf(kokiObj), "Unsupported Type")
	}
}

func detectAndConvertFromKubeObj(kubeObj runtime.Object) (interface{}, error) {
	switch kubeObj := kubeObj.(type) {
	case *exts.Deployment:
		return converters.Convert_Kube_v1beta2_Deployment_to_Koki_Deployment(kubeObj)
	case *v1.PersistentVolume:
		return converters.Convert_Kube_v1_PersistentVolume_to_Koki_PersistentVolume(kubeObj)
	case *v1.Pod:
		return converters.Convert_Kube_v1_Pod_to_Koki_Pod(kubeObj)
	case *v1.ReplicationController:
		return converters.Convert_Kube_v1_ReplicationController_to_Koki_ReplicationController(kubeObj)
	case *exts.ReplicaSet:
		return converters.Convert_Kube_v1beta2_ReplicaSet_to_Koki_ReplicaSet(kubeObj)
	case *v1.Service:
		return converters.Convert_Kube_v1_Service_to_Koki_Service(kubeObj)

	default:
		return nil, util.TypeErrorf(reflect.TypeOf(kubeObj), "Unsupported Type")
	}
}
