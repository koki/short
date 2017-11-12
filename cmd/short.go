package cmd

import (
	"github.com/ghodss/yaml"
	"github.com/kr/pretty"

	"github.com/koki/short/imports"
	"github.com/koki/short/param"
	"github.com/koki/short/parser"
)

func doFilesWithImports(filenames []string) error {
	for _, filename := range filenames {
		module, err := imports.Parse(filename)
		if err != nil {
			return err
		}

		evalContext := imports.EvalContext{
			RawToTyped:  parser.ParseKokiNativeObject,
			ApplyParams: param.ApplyParams,
		}

		err = evalContext.EvaluateModule(module)
		if err != nil {
			return err
		}

		bytes, err := yaml.Marshal(module.TypedResult)
		if err != nil {
			return err
		}
		pretty.Println(string(bytes))
	}

	return nil
}
