package parser

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	unstructuredconversion "k8s.io/apimachinery/pkg/conversion/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/ghodss/yaml"

	"github.com/koki/short/util"
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
		return nil, util.InvalidValueErrorf(u, "unsupported apiVersion/kind (is the manifest kube-native format?): %s", err.Error())
	}

	if err := unstructuredconversion.DefaultConverter.FromUnstructured(obj, typedObj); err != nil {
		return nil, util.InvalidValueForTypeErrorf(obj, typedObj, "couldn't convert to typed kube obj: %s", err.Error())
	}

	return typedObj, nil
}
