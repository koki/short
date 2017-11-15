package parser

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	unstructuredconversion "k8s.io/apimachinery/pkg/conversion/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/ghodss/yaml"
)

func ParseSingleKubeNativeFromBytes(data []byte) (runtime.Object, error) {
	obj := map[string]interface{}{}
	err := yaml.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	return ParseSingleKubeNative(obj)
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
