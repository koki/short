package imports

import (
	"path/filepath"

	"github.com/golang/glog"

	"github.com/koki/short/parser"
	"github.com/koki/short/util"
)

func (c *EvalContext) Parse(rootPath string) ([]Module, error) {
	objs, err := c.ReadFromPath(rootPath)
	if err != nil {
		return nil, util.InvalidValueContextErrorf(err, rootPath, "reading module")
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
	imports, hasImportsKey, err := c.parseImports(rootPath, obj)
	if err != nil {
		return nil, err
	}
	delete(obj, "imports")

	params, hasParamsKey, err := parseParamDefs(rootPath, obj)
	if err != nil {
		return nil, err
	}
	delete(obj, "params")

	exports, hasExportsKey, err := parseExports(rootPath, obj)
	if err != nil {
		return nil, err
	}
	delete(obj, "exports")

	if hasImportsKey || hasParamsKey || hasExportsKey {
		if len(obj) > 0 {
			return nil, util.InvalidValueErrorf(rootPath, "a koki module file can only have imports/params/exports as its top-level keys")
		}
	} else {
		if len(obj) != 1 {
			return nil, util.InvalidValueErrorf(rootPath, "a simple koki resource file must have exactly one top-level key")
		}

		// obj is a koki resource wrapper without imports
		return &Module{
			Path: rootPath,
			Exports: map[string]*Resource{
				"default": &Resource{
					Raw: obj,
				},
			},
		}, nil
	}

	return &Module{
		Path:    rootPath,
		Imports: imports,
		Params:  params,
		Exports: exports,
	}, nil
}

func parseExports(rootPath string, obj map[string]interface{}) (map[string]*Resource, bool, error) {
	hasExportsKey := false
	exports := map[string]*Resource{}
	if exportsObj, ok := obj["exports"]; ok {
		hasExportsKey = true
		if exportObjs, ok := exportsObj.([]interface{}); ok {
			for _, exportObj := range exportObjs {
				exportName, exportDef, err := parseExport(exportObj)
				if err != nil {
					return nil, hasExportsKey, util.ContextualizeErrorf(err, "couldn't parse an Export in (%s)", rootPath)
				}
				exports[exportName] = exportDef
			}
		} else {
			return nil, hasExportsKey, util.InvalidValueForTypeErrorf(exportsObj, exports, "expected array of exports in %s", rootPath)
		}
	}

	if len(exports) == 0 {
		return nil, hasExportsKey, nil
	}

	return exports, hasExportsKey, nil
}

func parseExport(obj interface{}) (string, *Resource, error) {
	def := &Resource{}
	if dict, ok := obj.(map[string]interface{}); ok {
		if val, ok := dict["value"]; ok {
			def.Raw = val
		} else {
			return "", def, util.InvalidValueForTypeErrorf(dict, def, "exports entry must contain a \"value\" key with the exported value")
		}
		if len(dict) != 2 {
			return "", def, util.InvalidValueForTypeErrorf(obj, def, "exports entry should have two keys, one for the exported name, and one for the exported value")
		}
		for name, description := range dict {
			if name != "value" {
				if descriptionStr, ok := description.(string); ok {
					def.Description = descriptionStr
					return name, def, nil
				}

				return "", def, util.InvalidValueForTypeErrorf(dict, def, "expected the export description to be a string value for name key (%s)", name)
			}
		}
	}

	return "", def, util.InvalidValueForTypeErrorf(obj, def, "expected exports entry to be a dictionary with two keys, one for the exported name, and one for the exported value")
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
					return nil, hasParamsKey, util.ContextualizeErrorf(err, "couldn't parse a Param in (%s)", rootPath)
				}
				params[paramName] = paramDef
			}
		} else {
			return nil, hasParamsKey, util.InvalidValueForTypeErrorf(paramsObj, params, "expected array of params in %s", rootPath)
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
					return name, def, util.InvalidValueForTypeErrorf(val, def, "interpreted key (%s) as param name. expected string value (for param description).", key)
				}
			}
		}

		return name, def, nil
	default:
		return "", def, util.InvalidValueForTypeErrorf(obj, def, "expected string or map")
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
						return nil, hasImportsKey, util.InvalidValueForTypeContextErrorf(err, imprt, Import{}, "processing import in module (%s)", rootPath)
					}
					imports = append(imports, anImport)
				} else {
					return nil, hasImportsKey, util.InvalidInstanceErrorf(imprt, "expected an import declaration in (%s)", rootPath)
				}
			}
		} else {
			return nil, hasImportsKey, util.InvalidInstanceErrorf(imprts, "expected array of imports in %s", rootPath)
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
				imp.Path, err = c.ResolveImportPath(rootPath, importPath)
				if err != nil {
					return nil, util.InvalidValueErrorf(importPath, "couldn't resolve 'absolute' path for import (%s) in module (%s)", importPath, rootPath)
				}
			} else {
				return nil, util.InvalidInstanceErrorf(imprt, "import path should be a string")
			}
		}
	}

	if len(imp.Name) == 0 {
		return nil, util.InvalidInstanceErrorf(imprt, "expected import name and path")
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
