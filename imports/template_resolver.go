package imports

import (
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
		for _, imprt := range module.Imports {
			if imprt.Name == ident {
				// Make sure the Import has been evaluated.
				err := c.EvaluateImport(module, params, imprt)
				if err != nil {
					return nil, err
				}

				return imprt.Module.Export.Raw, nil
			}
		}

		return nil, util.InvalidValueErrorf(ident, "invalid template param (%s) for (%s)", ident, module.Path)
	})
}
