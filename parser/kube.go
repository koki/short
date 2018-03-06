package parser

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
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
		return nil, serrors.InvalidValueContextErrorf(err, u, "unsupported apiVersion/kind (is the manifest kube-native format?)")
	}

	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj, typedObj); err != nil {
		return nil, serrors.InvalidValueForTypeContextErrorf(err, obj, typedObj, "couldn't convert to typed kube obj")
	}
	return typedObj, nil
}
