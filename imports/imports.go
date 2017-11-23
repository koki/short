package imports

import (
	"path/filepath"

	"github.com/golang/glog"

	"github.com/koki/short/parser"
	"github.com/koki/short/util"
)

func Parse(rootPath string) ([]Module, error) {
	objs, err := parser.Parse([]string{rootPath}, false)
	if err != nil {
		return nil, util.InvalidValueErrorf(rootPath, "error reading module: %s", err.Error())
	}

	if len(objs) > 1 {
		glog.V(1).Infof("(%s) has multiple sections. only the first section can be imported by other modules.", rootPath)
	}

	modules := make([]Module, len(objs))
	for i, obj := range objs {
		module, err := ParseComponent(rootPath, obj)
		if err != nil {
			return nil, err
		}

		modules[i] = *module
	}

	return modules, nil
}

func ParseComponent(rootPath string, obj map[string]interface{}) (*Module, error) {
	if len(obj) == 1 {
		// obj is a koki resource wrapper without imports
		return &Module{
			Path: rootPath,
			Raw:  obj,
		}, nil
	}

	imports, err := parseImports(rootPath, obj)
	if err != nil {
		return nil, err
	}

	params, err := parseParamDefs(rootPath, obj)
	if err != nil {
		return nil, err
	}

	kokiResource := map[string]interface{}{}
	for key, val := range obj {
		if key != "imports" && key != "params" {
			kokiResource[key] = val
		}
	}

	return &Module{
		Path:    rootPath,
		Imports: imports,
		Params:  params,
		Raw:     kokiResource,
	}, nil
}

func parseParamDefs(rootPath string, obj map[string]interface{}) (map[string]ParamDef, error) {
	params := map[string]ParamDef{}
	if paramsObj, ok := obj["params"]; ok {
		if paramObjs, ok := paramsObj.([]interface{}); ok {
			for _, paramObj := range paramObjs {
				paramName, paramDef, err := parseParamDef(paramObj)
				if err != nil {
					return nil, util.ContextualizeErrorf(err, "couldn't parse a Param in (%s)", rootPath)
				}
				params[paramName] = paramDef
			}
		} else {
			return nil, util.InvalidValueForTypeErrorf(paramsObj, params, "expected array of params in %s", rootPath)
		}
	}

	if len(params) == 0 {
		return nil, nil
	}

	return params, nil
}

func parseParamDef(obj interface{}) (string, ParamDef, error) {
	def := ParamDef{}
	switch obj := obj.(type) {
	case string:
		return obj, def, nil
	case map[string]interface{}:
		name := ""
		for key, val := range obj {
			switch key {
			case "default":
				def.Default = val
			default:
				name = key
				if description, ok := val.(string); ok {
					def.Description = description
				} else {
					return name, def, util.InvalidValueForTypeErrorf(val, def, "interpreted key (%s) as param name. expected string value (for param description).", key)
				}
			}
		}

		return name, def, nil
	default:
		return "", def, util.InvalidValueForTypeErrorf(obj, def, "expected string or map")
	}
}

func parseImports(rootPath string, obj map[string]interface{}) ([]*Import, error) {
	imports := []*Import{}
	if imprts, ok := obj["imports"]; ok {
		if imprts, ok := imprts.([]interface{}); ok {
			for _, imprt := range imprts {
				if imprt, ok := imprt.(map[string]interface{}); ok {
					anImport, err := parseImport(rootPath, imprt)
					if err != nil {
						return nil, util.InvalidValueForTypeErrorf(imprt, Import{}, "error processing import in module (%s): %s", rootPath, err.Error())
					}
					imports = append(imports, anImport)
				} else {
					return nil, util.InvalidInstanceErrorf(imprt, "expected an import declaration in (%s)", rootPath)
				}
			}
		} else {
			return nil, util.InvalidInstanceErrorf(imprts, "expected array of imports in %s", rootPath)
		}
	}

	if len(imports) == 0 {
		return nil, nil
	}

	return imports, nil
}

func parseImport(rootPath string, imprt map[string]interface{}) (*Import, error) {
	var err error
	if len(imprt) == 0 {
		return nil, util.InvalidInstanceErrorf(imprt, "empty import declaration")
	}
	if len(imprt) > 2 {
		return nil, util.InvalidInstanceErrorf(imprt, "import declaration should have at most params and name:path")
	}

	imp := &Import{}
	for key, val := range imprt {
		if key == "params" {
			if params, ok := val.(map[string]interface{}); ok {
				imp.Params = params
			} else {
				return nil, util.InvalidInstanceErrorf(imprt, "params should be a dictionary")
			}
		} else {
			imp.Name = key
			if importPath, ok := val.(string); ok {
				imp.Path = convertImportPath(rootPath, importPath)
			} else {
				return nil, util.InvalidInstanceErrorf(imprt, "import path should be a string")
			}
		}
	}

	if len(imp.Name) == 0 {
		return nil, util.InvalidInstanceErrorf(imprt, "expected import name and path")
	}

	importModules, err := Parse(imp.Path)
	if err != nil {
		return nil, err
	}
	imp.Module = &importModules[0]

	return imp, nil
}

func convertImportPath(rootPath string, relativePath string) string {
	dirPath, _ := filepath.Split(rootPath)
	return filepath.Join(dirPath, relativePath)
}
