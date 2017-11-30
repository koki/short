package imports

import (
	"strings"

	"github.com/koki/short/template"
	"github.com/koki/short/util"
)

func (c *EvalContext) ResolverForModule(module *Module, params map[string]interface{}) template.Resolver {
	return template.Resolver(func(ident string) (interface{}, error) {
		// Check params for the identifier.
		if val, ok := params[ident]; ok {
			return val, nil
		}

		// Check imports for the identifier.
		identSegments := strings.Split(ident, ".")
		importName := identSegments[0]
		if len(identSegments) > 2 {
			return nil, util.InvalidValueErrorf(ident, "cannot index into an imported resource. (%s) can have at most two segments. in module (%s)", ident, module.Path)
		}
		exportName := "default"
		if len(identSegments) > 1 {
			exportName = identSegments[1]
		}

		for _, imprt := range module.Imports {
			if imprt.Name == importName {
				// Make sure the Import has been evaluated.
				err := c.EvaluateImport(module, params, imprt)
				if err != nil {
					return nil, err
				}

				if val, ok := imprt.Module.Exports[exportName]; ok {
					return val.Raw, nil
				}
			}
		}

		return nil, util.InvalidValueErrorf(ident, "invalid template param (%s) for (%s)", ident, module.Path)
	})
}
