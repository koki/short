package imports

import (
	"fmt"
	"path/filepath"

	"github.com/koki/short/util"
)

func Parse(rootPath string) (*Module, error) {
	components, err := ReadYamls(rootPath)
	if err != nil {
		return nil, err
	}

	if len(components) > 2 {
		return nil, fmt.Errorf("file (%s) should have at most one imports section and one manifest section", rootPath)
	}

	if len(components) == 1 {
		return &Module{
			Path: rootPath,
			Raw:  components[0],
		}, nil
	}

	imports := []*Import{}
	if imprtsDoc, ok := components[0].(map[string]interface{}); ok {
		if imprts, ok := imprtsDoc["imports"]; ok {
			if imprts, ok := imprts.([]interface{}); ok {
				for _, imprt := range imprts {
					if imprt, ok := imprt.(map[string]interface{}); ok {
						anImport, err := parseImport(rootPath, imprt)
						if err != nil {
							return nil, err
						}
						imports = append(imports, anImport)
					} else {
						return nil, util.InvalidInstanceErrorf(imprt, "expected an import declaration in (%s)", rootPath)
					}
				}
			} else {
				return nil, util.InvalidInstanceErrorf(imprts, "expected array of imports in %s", rootPath)
			}
		} else {
			return nil, fmt.Errorf("file (%s) should have 'imports' as its first section", rootPath)
		}
	} else {
		return nil, util.InvalidInstanceErrorf(components[0], "imports section should be a map in (%s)", rootPath)
	}

	return &Module{
		Path:    rootPath,
		Imports: imports,
		Raw:     components[1],
	}, nil
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

	imp.Module, err = Parse(imp.Path)
	if err != nil {
		return nil, err
	}

	return imp, nil
}

func convertImportPath(rootPath string, relativePath string) string {
	dirPath, _ := filepath.Split(rootPath)
	return filepath.Join(dirPath, relativePath)
}
