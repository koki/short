package imports

import (
	"fmt"

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

func (c *EvalContext) EvaluateImport(inModule *Module, imprt *Import) error {
	if imprt.IsEvaluated {
		return nil
	}

	// Evaluate the Module.
	err := c.EvaluateModule(imprt.Module)
	if err != nil {
		return err
	}

	params, err := template.ReplaceMap(imprt.Params, c.ResolverForModule(inModule))
	if err != nil {
		return err
	}
	imprt.Params = params

	err = c.ApplyParams(imprt.Params, imprt.Module)
	if err != nil {
		return err
	}

	imprt.IsEvaluated = true
	return nil
}

func (c *EvalContext) EvaluateModule(module *Module) error {
	if module.IsEvaluated {
		return nil
	}

	var err error

	raw, err := template.ReplaceMap(module.Raw, c.ResolverForModule(module))
	if err != nil {
		return err
	}
	module.Raw = raw

	module.TypedResult, err = c.RawToTyped(module.Raw)
	if err != nil {
		module.TypedResult = err
	}

	module.IsEvaluated = true
	return nil
}

func (c *EvalContext) ResolverForModule(module *Module) template.Resolver {
	return template.Resolver(func(ident string) (interface{}, error) {
		for _, imprt := range module.Imports {
			if imprt.Name == ident {
				// Make sure the Import has been evaluated.
				err := c.EvaluateImport(module, imprt)
				if err != nil {
					return nil, err
				}

				_, val, err := util.GetOnlyMapEntry(imprt.Module.Raw)
				return val, err
			}
		}

		return nil, fmt.Errorf("no value for template param (%s) in file (%s)", ident, module.Path)
	})
}
