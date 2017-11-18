package client

import (
	"encoding/json"
	"io"

	"github.com/ghodss/yaml"

	"github.com/koki/short/converter"
	"github.com/koki/short/parser"
	"github.com/koki/short/util"
)

/*

Converting text streams to objects.

*/

// ConvertKokiStreams to Kube objects.
func ConvertKokiStreams(kokiStreams []io.ReadCloser) ([]interface{}, error) {
	objs, err := parser.ParseStreams(kokiStreams)
	if err != nil {
		return nil, err
	}

	return converter.ConvertToKubeNative(objs)
}

// ConvertKubeStreams to Koki objects.
func ConvertKubeStreams(kubeStreams []io.ReadCloser) ([]interface{}, error) {
	objs, err := parser.ParseStreams(kubeStreams)
	if err != nil {
		return nil, err
	}

	return converter.ConvertToKokiNative(objs)
}

// ConvertEitherStreamsToKube either Koki or Kube to just Kube objects.
func ConvertEitherStreamsToKube(eitherStreams []io.ReadCloser) ([]interface{}, error) {
	objs, err := parser.ParseStreams(eitherStreams)
	if err != nil {
		return nil, err
	}

	kubeObjs := make([]interface{}, len(eitherStreams))
	for i, obj := range objs {
		kubeObjs[i], err = converter.ConvertOneToKubeNative(obj)
		if err == nil {
			continue
		}

		kubeObjs[i], err = parser.ParseSingleKubeNative(obj)
		if err != nil {
			return nil, util.InvalidValueErrorf(obj, "couldn't parse as Kube or Koki resource")
		}
	}

	return kubeObjs, nil
}

func WriteObjsToYamlStream(objs []interface{}, yamlStream io.Writer) error {
	var err error
	for i, obj := range objs {
		if i > 0 {
			_, err = yamlStream.Write([]byte("---\n"))
			if err != nil {
				return err
			}
		}

		b, err := yaml.Marshal(obj)
		if err != nil {
			return util.InvalidValueErrorf(obj, "couldn't serialize as yaml")
		}
		_, err = yamlStream.Write(b)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteObjsToJSONStream(objs []interface{}, jsonStream io.Writer) error {
	var err error
	for i, obj := range objs {
		if i > 0 {
			_, err = jsonStream.Write([]byte("\n"))
			if err != nil {
				return err
			}
		}

		b, err := json.Marshal(obj)
		if err != nil {
			return util.InvalidValueErrorf(obj, "couldn't serialize as json")
		}
		_, err = jsonStream.Write(b)
		if err != nil {
			return err
		}
	}

	return nil
}
