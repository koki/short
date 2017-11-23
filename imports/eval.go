package imports

import (
	"github.com/koki/short/template"
	"github.com/koki/short/util"
)

/*

How evaluation works:

Every module starts out Raw and unevaluated.

An import is evaluated by:
  1. Evaluate its Module.
  2. Apply its Params to its Module using the other Imports.

A module is evaluated by:
  1. Build its Result by filling its Raw template from the Module.Raw of its Imports.
  2. Parse its TypedResult

*/

func (c *EvalContext) EvaluateImport(inModule *Module, inModuleParams map[string]interface{}, imprt *Import) error {
	if imprt.IsEvaluated {
		return nil
	}

	// Figure out the values for our parameters.
	params, err := template.ReplaceMap(imprt.Params, c.ResolverForModule(inModule, inModuleParams))
	if err != nil {
		return err
	}
	imprt.Params = params

	// Evaluate the Module with these parameters.
	err = c.EvaluateModule(imprt.Module, imprt.Params)
	if err != nil {
		return err
	}

	imprt.IsEvaluated = true
	return nil
}

func (c *EvalContext) EvaluateModule(module *Module, params map[string]interface{}) error {
	if module.IsEvaluated {
		return nil
	}

	// Lazy-evaluate all variables and substitute them into any template holes.
	// (ResolverForModule does the lazy evaluation.)
	// Evaluation is lazy so we don't have to determine ahead-of-time which order to evaluate the imports in.
	raw, err := template.ReplaceMap(module.Raw, c.ResolverForModule(module, params))
	if err != nil {
		return err
	}
	module.Raw = raw

	// All template substitutions should be complete. Evaluate all exports.
	module.TypedResult, err = c.RawToTyped(module.Raw)
	if err != nil {
		module.TypedResult = err
	}

	module.IsEvaluated = true
	return nil
}

func (c *EvalContext) ResolverForModule(module *Module, params map[string]interface{}) template.Resolver {
	return template.Resolver(func(ident string) (interface{}, error) {
		if val, ok := params[ident]; ok {
			return val, nil
		}

		for _, imprt := range module.Imports {
			if imprt.Name == ident {
				// Make sure the Import has been evaluated.
				err := c.EvaluateImport(module, params, imprt)
				if err != nil {
					return nil, err
				}

				_, val, err := util.GetOnlyMapEntry(imprt.Module.Raw)
				return val, err
			}
		}

		return nil, util.InvalidValueErrorf(ident, "invalid template param (%s) for (%s)", ident, module.Path)
	})
}
