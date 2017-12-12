package imports

import (
	"path/filepath"

	"github.com/golang/glog"

	"github.com/koki/short/parser"
	serrors "github.com/koki/structurederrors"
)

func (c *EvalContext) Parse(rootPath string) ([]Module, error) {
	objs, err := c.ReadFromPath(rootPath)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, rootPath, "reading module")
	}

	if len(objs) > 1 {
		glog.V(1).Infof("(%s) has multiple sections. only the first section can be imported by other modules.", rootPath)
	}

	modules := make([]Module, len(objs))
	for i, obj := range objs {
		module, err := c.ParseComponent(rootPath, obj)
		if err != nil {
			return nil, err
		}

		modules[i] = *module
	}

	return modules, nil
}

func (c *EvalContext) ParseComponent(rootPath string, obj map[string]interface{}) (*Module, error) {
	imports, _, err := c.parseImports(rootPath, obj)
	if err != nil {
		return nil, err
	}
	delete(obj, "imports")

	params, _, err := parseParamDefs(rootPath, obj)
	if err != nil {
		return nil, err
	}
	delete(obj, "params")

	// Last remaining key must be the exported resource.
	if len(obj) != 1 {
		return nil, serrors.InvalidValueErrorf(rootPath, "koki module must contain exactly one resource")
	}

	export := Resource{
		Raw: obj,
	}

	return &Module{
		Path:    rootPath,
		Imports: imports,
		Params:  params,
		Export:  export,
	}, nil
}

func parseParamDefs(rootPath string, obj map[string]interface{}) (map[string]ParamDef, bool, error) {
	hasParamsKey := false
	params := map[string]ParamDef{}
	if paramsObj, ok := obj["params"]; ok {
		hasParamsKey = true
		if paramObjs, ok := paramsObj.([]interface{}); ok {
			for _, paramObj := range paramObjs {
				paramName, paramDef, err := parseParamDef(paramObj)
				if err != nil {
					return nil, hasParamsKey, serrors.ContextualizeErrorf(err, "couldn't parse a Param in (%s)", rootPath)
				}
				params[paramName] = paramDef
			}
		} else {
			return nil, hasParamsKey, serrors.InvalidValueForTypeErrorf(paramsObj, params, "expected array of params in %s", rootPath)
		}
	}

	if len(params) == 0 {
		return nil, hasParamsKey, nil
	}

	return params, hasParamsKey, nil
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
					return name, def, serrors.InvalidValueForTypeErrorf(val, def, "interpreted key (%s) as param name. expected string value (for param description).", key)
				}
			}
		}

		return name, def, nil
	default:
		return "", def, serrors.InvalidValueForTypeErrorf(obj, def, "expected string or map")
	}
}

func (c *EvalContext) parseImports(rootPath string, obj map[string]interface{}) ([]*Import, bool, error) {
	hasImportsKey := false
	imports := []*Import{}
	if imprts, ok := obj["imports"]; ok {
		hasImportsKey = true
		if imprts, ok := imprts.([]interface{}); ok {
			for _, imprt := range imprts {
				if imprt, ok := imprt.(map[string]interface{}); ok {
					anImport, err := c.parseImport(rootPath, imprt)
					if err != nil {
						return nil, hasImportsKey, serrors.InvalidValueForTypeContextErrorf(err, imprt, Import{}, "processing import in module (%s)", rootPath)
					}
					imports = append(imports, anImport)
				} else {
					return nil, hasImportsKey, serrors.InvalidInstanceErrorf(imprt, "expected an import declaration in (%s)", rootPath)
				}
			}
		} else {
			return nil, hasImportsKey, serrors.InvalidInstanceErrorf(imprts, "expected array of imports in %s", rootPath)
		}
	}

	if len(imports) == 0 {
		return nil, hasImportsKey, nil
	}

	return imports, hasImportsKey, nil
}

func (c *EvalContext) parseImport(rootPath string, imprt map[string]interface{}) (*Import, error) {
	var err error
	if len(imprt) == 0 {
		return nil, serrors.InvalidInstanceErrorf(imprt, "empty import declaration")
	}
	if len(imprt) > 2 {
		return nil, serrors.InvalidInstanceErrorf(imprt, "import declaration should have at most params and name:path")
	}

	imp := &Import{}
	for key, val := range imprt {
		if key == "params" {
			if params, ok := val.(map[string]interface{}); ok {
				imp.Params = params
			} else {
				return nil, serrors.InvalidInstanceErrorf(imprt, "params should be a dictionary")
			}
		} else {
			imp.Name = key
			if importPath, ok := val.(string); ok {
				imp.Path, err = c.ResolveImportPath(rootPath, importPath)
				if err != nil {
					return nil, serrors.InvalidValueErrorf(importPath, "couldn't resolve 'absolute' path for import (%s) in module (%s)", importPath, rootPath)
				}
			} else {
				return nil, serrors.InvalidInstanceErrorf(imprt, "import path should be a string")
			}
		}
	}

	if len(imp.Name) == 0 {
		return nil, serrors.InvalidInstanceErrorf(imprt, "expected import name and path")
	}

	importModules, err := c.Parse(imp.Path)
	if err != nil {
		return nil, err
	}
	imp.Module = &importModules[0]

	return imp, nil
}

func ResolveImportLocalPath(rootPath string, importPath string) (string, error) {
	if len(rootPath) > 0 {
		dirPath, _ := filepath.Split(rootPath)
		return filepath.Join(dirPath, importPath), nil
	}

	return importPath, nil
}

func ReadFromLocalPath(path string) ([]map[string]interface{}, error) {
	return parser.Parse([]string{path}, false)
}
