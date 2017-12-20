package client

import (
	"io"

	"github.com/koki/json"
	"github.com/koki/json/jsonutil"
	"github.com/koki/short/converter"
	"github.com/koki/short/parser"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

/*

Converting text streams and JSON/YAML dictionaries to objects. (And also the reverse.)

This package is the canonical integration point for using Koki Short functionality.
It's used for the command-line tool and functional tests.

*/

// ConvertKokiStreams to Kube objects.
func ConvertKokiStreams(kokiStreams []io.ReadCloser) ([]interface{}, error) {
	objs, err := parser.ParseStreams(kokiStreams)
	if err != nil {
		return nil, err
	}

	return ConvertKokiMaps(objs)
}

func ConvertKokiMaps(objs []map[string]interface{}) ([]interface{}, error) {
	convertedObjs := make([]interface{}, len(objs))
	for i, obj := range objs {
		// 1. Parse.
		parsedObj, err := parser.ParseKokiNativeObject(obj)
		if err != nil {
			return nil, err
		}

		// 2. Check for unparsed fields--potential typos.
		extraneousPaths, err := jsonutil.ExtraneousFieldPaths(obj, parsedObj)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "checking for extraneous fields in input")
		}
		if len(extraneousPaths) > 0 {
			return nil, &jsonutil.ExtraneousFieldsError{Paths: extraneousPaths}
		}

		// 3. Convert.
		convertedObj, err := converter.DetectAndConvertFromKokiObj(parsedObj)
		if err != nil {
			return nil, err
		}
		convertedObjs[i] = convertedObj
	}

	return convertedObjs, nil
}

// ConvertKubeStreams to Koki objects.
func ConvertKubeStreams(kubeStreams []io.ReadCloser) ([]interface{}, error) {
	objs, err := parser.ParseStreams(kubeStreams)
	if err != nil {
		return nil, err
	}

	return ConvertKubeMaps(objs)
}

func ConvertKubeMaps(objs []map[string]interface{}) ([]interface{}, error) {
	convertedObjs := make([]interface{}, len(objs))
	for i, obj := range objs {
		// 1. Parse.
		parsedObj, err := parser.ParseSingleKubeNative(obj)
		if err != nil {
			return nil, err
		}

		// 2. Check for unparsed fields--potential typos.
		extraneousPaths, err := jsonutil.ExtraneousFieldPaths(obj, parsedObj)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "checking for extraneous fields in input")
		}
		if len(extraneousPaths) > 0 {
			return nil, &jsonutil.ExtraneousFieldsError{Paths: extraneousPaths}
		}

		// 3. Convert.
		convertedObj, err := converter.DetectAndConvertFromKubeObj(parsedObj)
		if err != nil {
			return nil, err
		}
		convertedObjs[i] = convertedObj
	}

	return convertedObjs, nil
}

// ConvertEitherStreamsToKube either Koki or Kube to just Kube objects.
func ConvertEitherStreamsToKube(eitherStreams []io.ReadCloser) ([]interface{}, error) {
	objs, err := parser.ParseStreams(eitherStreams)
	if err != nil {
		return nil, err
	}

	kubeObjs := make([]interface{}, len(objs))
	for i, obj := range objs {
		kubeObjs[i], err = converter.ConvertOneToKubeNative(obj)
		if err == nil {
			continue
		}

		kubeObjs[i], err = parser.ParseSingleKubeNative(obj)
		if err != nil {
			return nil, serrors.InvalidValueErrorf(obj, "couldn't parse as Kube or Koki resource")
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
			return serrors.InvalidValueErrorf(obj, "couldn't serialize as yaml")
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

		b, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return serrors.InvalidValueErrorf(obj, "couldn't serialize as json")
		}
		_, err = jsonStream.Write(b)
		if err != nil {
			return err
		}
	}

	return nil
}
