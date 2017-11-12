package imports

type Import struct {
	Name   string
	Path   string
	Params map[string]interface{}

	// IsEvaluated have the Params been applied to the Module?
	IsEvaluated bool

	Module *Module
}

type Module struct {
	Path    string
	Imports []*Import

	// Raw yaml parsed as string or map[string]interface{}
	Raw interface{}

	// IsEvaluated has the Raw yaml been evaluated (template holes filled, etc)?
	IsEvaluated bool

	// TypedResult a koki object (or other object with special meaning)
	TypedResult interface{}
}

type EvalContext struct {
	RawToTyped  func(raw interface{}) (interface{}, error)
	ApplyParams func(params map[string]interface{}, module *Module) error
}
