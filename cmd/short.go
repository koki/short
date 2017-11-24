package cmd

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"

	// Just make sure we also build the client package.
	_ "github.com/koki/short/client"
	"github.com/koki/short/converter"
	"github.com/koki/short/imports"
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
			RawToTyped: parser.ParseKokiNativeObject,
		}

		for _, module := range modules {
			err = evalContext.EvaluateModule(&module, nil)
			if err != nil {
				debugLogModule(module)
				return nil, err
			}

			for _, export := range module.Exports {
				if err, ok := export.TypedResult.(error); ok {
					debugLogModule(module)
					return nil, err
				}
			}

			results = append(results, module)
		}
	}

	return results, nil
}

func convertKokiModules(kokiModules []imports.Module) ([]interface{}, error) {
	kubeObjs := []interface{}{}
	for _, kokiModule := range kokiModules {
		for _, kokiExport := range kokiModule.Exports {
			kubeObj, err := converter.DetectAndConvertFromKokiObj(kokiExport.TypedResult)
			if err != nil {
				debugLogModule(kokiModule)
				return nil, err
			}
			kubeObjs = append(kubeObjs, kubeObj)
		}
	}

	return kubeObjs, nil
}
