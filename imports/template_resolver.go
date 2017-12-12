package imports

import (
	"strings"

	"github.com/koki/json/jsonutil"
	"github.com/koki/short/template"
	serrors "github.com/koki/structurederrors"
)

func (c *EvalContext) ResolverForModule(module *Module, params map[string]interface{}) template.Resolver {
	return template.Resolver(func(ident string) (interface{}, error) {
		// Split identifier into segments.
		segments := strings.Split(ident, ".")

		// Check params for the identifier.
		if param, ok := params[segments[0]]; ok {
			val, err := jsonutil.AtPathIn(param, segments[1:])
			if err != nil {
				return nil, serrors.InvalidValueContextErrorf(err, param, "resolving %s", ident)
			}
			return val, nil
		}

		// Check imports for the identifier.
		for _, imprt := range module.Imports {
			if imprt.Name == segments[0] {
				// Make sure the Import has been evaluated.
				err := c.EvaluateImport(module, params, imprt)
				if err != nil {
					return nil, serrors.ContextualizeErrorf(err, "resolving %s", ident)
				}

				_, export, err := jsonutil.GetOnlyMapEntry(imprt.Module.Export.Raw)
				if err != nil {
					return nil, serrors.InvalidValueContextErrorf(err, imprt.Module.Export.Raw, "module should export exactly one top-level key")
				}
				val, err := jsonutil.AtPathIn(export, segments[1:])
				if err != nil {
					return nil, serrors.InvalidValueContextErrorf(err, imprt, "resolving %s", ident)
				}
				return val, nil
			}
		}

		return nil, serrors.InvalidValueErrorf(ident, "invalid template param (%s) for (%s)", ident, module.Path)
	})
}
