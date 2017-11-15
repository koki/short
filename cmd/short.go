package cmd

import (
	"github.com/koki/short/converter"
	"github.com/koki/short/imports"
	"github.com/koki/short/param"
	"github.com/koki/short/parser"
)

func loadKokiFiles(filenames []string) ([]interface{}, error) {
	results := []interface{}{}
	for _, filename := range filenames {
		module, err := imports.Parse(filename)
		if err != nil {
			return nil, err
		}

		evalContext := imports.EvalContext{
			RawToTyped:  parser.ParseKokiNativeObject,
			ApplyParams: param.ApplyParams,
		}

		err = evalContext.EvaluateModule(module)
		if err != nil {
			return nil, err
		}

		if err, ok := module.TypedResult.(error); ok {
			return nil, err
		}

		if err, ok := module.TypedResult.(error); ok {
			return nil, err
		}

		results = append(results, module.TypedResult)
	}

	return results, nil
}

func convertKokiObjs(kokiObjs []interface{}) ([]interface{}, error) {
	var err error
	kubeObjs := make([]interface{}, len(kokiObjs))
	for i, kokiObj := range kokiObjs {
		kubeObjs[i], err = converter.DetectAndConvertFromKokiObj(kokiObj)
		if err != nil {
			return nil, err
		}
	}

	return kubeObjs, nil
}
