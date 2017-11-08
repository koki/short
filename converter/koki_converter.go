package converter

import (
	"reflect"

	"github.com/koki/short/converter/converters"
	"github.com/koki/short/types"
	"github.com/koki/short/util"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func detectAndConvertFromKokiObj(kokiObj interface{}) (interface{}, error) {
	switch kokiObj.(type) {
	case *types.PodWrapper:
		return converters.Convert_Koki_Pod_to_Kube_v1_Pod(kokiObj.(*types.PodWrapper))
	default:
		return nil, util.TypeErrorf(reflect.TypeOf(kokiObj), "Unsupported Type")
	}
}

func detectAndConvertFromKubeObj(kubeObj runtime.Object) (interface{}, error) {
	switch kubeObj.(type) {
	case *v1.Pod:
		kokiObj, err := converters.Convert_Kube_v1_Pod_to_Koki_Pod(kubeObj.(*v1.Pod))
		return kokiObj, err
	case *v1.Service:
		kokiObj, err := converters.Convert_Kube_v1_Service_to_Koki_Service(kubeObj.(*v1.Service))
		return kokiObj, err
	default:
		return nil, util.TypeErrorf(reflect.TypeOf(kubeObj), "Unsupported Type")
	}
}
