package imports

import (
	"github.com/koki/short/template"
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

	// Fill in default values for missing params.
	if params == nil {
		params = map[string]interface{}{}
	}
	for paramName, paramDef := range module.Params {
		if paramDef.Default != nil {
			if _, ok := params[paramName]; ok {
				continue
			}

			params[paramName] = paramDef.Default
		}
	}

	err := c.EvaluateExport(module, params, &module.Export)
	if err != nil {
		return err
	}

	module.IsEvaluated = true
	return nil
}

func (c *EvalContext) EvaluateExport(module *Module, params map[string]interface{}, export *Resource) error {
	// Lazy-evaluate all variables and substitute them into any template holes.
	// (ResolverForModule does the lazy evaluation.)
	// Evaluation is lazy so we don't have to determine ahead-of-time which order to evaluate the imports in.
	raw, err := template.ReplaceAny(export.Raw, c.ResolverForModule(module, params))
	if err != nil {
		return err
	}
	export.Raw = raw

	// All template substitutions should be complete. Evaluate the result.
	export.TypedResult, err = c.RawToTyped(export.Raw)
	if err != nil {
		export.TypedResult = err
	}

	return nil
}
