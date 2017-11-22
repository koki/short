package cmd

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"

	// Just make sure we also build the client package.
	_ "github.com/koki/short/client"
	"github.com/koki/short/converter"
	"github.com/koki/short/imports"
	"github.com/koki/short/param"
	"github.com/koki/short/parser"
)

func debugLogModule(module imports.Module) {
	trimmed := imports.TrimToDepth(&module, debugImportsDepth)

	b, err := yaml.Marshal(module)
	if err != nil {
		glog.V(0).Info("couldn't output loaded module as yaml")
	}

	if debugImportsDepth == defaultDebugImportsDepth && trimmed {
		glog.V(0).Infof("Use flag '--debug-imports-depth <number>' to show more levels of imports.\n%s", string(b))
	} else {
		glog.V(0).Infof("Couldn't load module:\n%s", string(b))
	}
}

func loadKokiFiles(filenames []string) ([]imports.Module, error) {
	results := []imports.Module{}
	for _, filename := range filenames {
		modules, err := imports.Parse(filename)
		if err != nil {
			return nil, err
		}

		evalContext := imports.EvalContext{
			RawToTyped:  parser.ParseKokiNativeObject,
			ApplyParams: param.ApplyParams,
		}

		for _, module := range modules {
			err = evalContext.EvaluateModule(&module)
			if err != nil {
				debugLogModule(module)
				return nil, err
			}

			if err, ok := module.TypedResult.(error); ok {
				debugLogModule(module)
				return nil, err
			}

			results = append(results, module)
		}
	}

	return results, nil
}

func convertKokiModules(kokiModules []imports.Module) ([]interface{}, error) {
	var err error
	kubeObjs := make([]interface{}, len(kokiModules))
	for i, kokiModule := range kokiModules {
		kubeObjs[i], err = converter.DetectAndConvertFromKokiObj(kokiModule.TypedResult)
		if err != nil {
			debugLogModule(kokiModule)
			return nil, err
		}
	}

	return kubeObjs, nil
}
