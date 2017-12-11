package imports

import (
	"strings"

	"github.com/koki/short/template"
	"github.com/koki/short/util"
	"github.com/koki/short/util/objutil"
)

func (c *EvalContext) ResolverForModule(module *Module, params map[string]interface{}) template.Resolver {
	return template.Resolver(func(ident string) (interface{}, error) {
		// Split identifier into segments.
		segments := strings.Split(ident, ".")

		// Check params for the identifier.
		if param, ok := params[segments[0]]; ok {
			val, err := objutil.AtPathIn(param, segments[1:])
			if err != nil {
				return nil, util.InvalidValueContextErrorf(err, param, "resolving %s", ident)
			}
			return val, nil
		}

		// Check imports for the identifier.
		for _, imprt := range module.Imports {
			if imprt.Name == segments[0] {
				// Make sure the Import has been evaluated.
				err := c.EvaluateImport(module, params, imprt)
				if err != nil {
					return nil, util.ContextualizeErrorf(err, "resolving %s", ident)
				}

				_, export, err := objutil.GetOnlyMapEntry(imprt.Module.Export.Raw)
				if err != nil {
					return nil, util.InvalidValueContextErrorf(err, imprt.Module.Export.Raw, "module should export exactly one top-level key")
				}
				val, err := objutil.AtPathIn(export, segments[1:])
				if err != nil {
					return nil, util.InvalidValueContextErrorf(err, imprt, "resolving %s", ident)
				}
				return val, nil
			}
		}

		return nil, util.InvalidValueErrorf(ident, "invalid template param (%s) for (%s)", ident, module.Path)
	})
}
