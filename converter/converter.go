package converter

import (
	unstructuredconversion "k8s.io/apimachinery/pkg/conversion/unstructured"

	"github.com/koki/short/parser"
)

func ConvertToKubeNative(objs []map[string]interface{}) ([]interface{}, error) {
	convertedTypes := []interface{}{}
	for i := range objs {
		obj := objs[i]

		typedObj, err := parser.ParseKokiNativeObject(obj)
		if err != nil {
			return nil, err
		}

		kubeObj, err := DetectAndConvertFromKokiObj(typedObj)
		if err != nil {
			return nil, err
		}
		convertedTypes = append(convertedTypes, kubeObj)
	}
	return convertedTypes, nil
}

func ConvertToKokiNative(objs []map[string]interface{}) ([]interface{}, error) {
	convertedTypes := []interface{}{}
	for i := range objs {
		obj := objs[i]
		typedObj, err := parser.ParseSingleKubeNative(obj)
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
