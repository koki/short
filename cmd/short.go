package cmd

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"

	// Just make sure we also build the client package.
	_ "github.com/koki/short/client"
	"github.com/koki/short/converter"
	"github.com/koki/short/imports"
	"github.com/koki/short/parser"
	"github.com/koki/short/util"
	"github.com/koki/short/util/objutil"
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
		evalContext := imports.EvalContext{
			RawToTyped:        parser.ParseKokiNativeObject,
			ResolveImportPath: imports.ResolveImportLocalPath,
			ReadFromPath:      imports.ReadFromLocalPath,
		}

		modules, err := evalContext.Parse(filename)
		if err != nil {
			return nil, err
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
			if data, ok := kokiExport.Raw.(map[string]interface{}); ok {
				extraneousPaths, err := objutil.ExtraneousFieldPaths(data, kokiExport.TypedResult)
				if err != nil {
					return nil, util.ContextualizeErrorf(err, "checking for extraneous fields in input")
				}
				if len(extraneousPaths) > 0 {
					return nil, &objutil.ExtraneousFieldsError{
						Paths: extraneousPaths,
					}
				}
			}

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
