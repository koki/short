package converter

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	unstructuredconversion "k8s.io/apimachinery/pkg/conversion/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/koki/short/parser"
)

func ConvertToKubeNative(objs []map[string]interface{}) ([]interface{}, error) {
	convertedTypes := []interface{}{}
	for i := range objs {
		obj := objs[i]

		typedObj, err := parser.ParseKokiNativeObject(obj)
		if err != nil {
			return nil, nil
		}

		kubeObj, err := DetectAndConvertFromKokiObj(typedObj)
		if err != nil {
			return nil, err
		}
		convertedTypes = append(convertedTypes, kubeObj)
	}
	return convertedTypes, nil
}

func ParseSingleKubeNative(obj map[string]interface{}) (runtime.Object, error) {
	u := &unstructured.Unstructured{
		Object: obj,
	}

	typedObj, err := creator.New(u.GetObjectKind().GroupVersionKind())
	if err != nil {
		return nil, err
	}

	if err := unstructuredconversion.DefaultConverter.FromUnstructured(obj, typedObj); err != nil {
		return nil, err
	}

	return typedObj, err
}

func ConvertToKokiNative(objs []map[string]interface{}) ([]interface{}, error) {
	convertedTypes := []interface{}{}
	for i := range objs {
		obj := objs[i]
		typedObj, err := ParseSingleKubeNative(obj)
		if err != nil {
			return nil, err
		}

		if err := unstructuredconversion.DefaultConverter.FromUnstructured(obj, typedObj); err != nil {
			return nil, err
		}

		kokiObj, err := DetectAndConvertFromKubeObj(typedObj)
		if err != nil {
			return nil, err
		}

		convertedTypes = append(convertedTypes, kokiObj)
	}

	return convertedTypes, nil
}
