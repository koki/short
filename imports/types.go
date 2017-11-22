package imports

type Import struct {
	Name   string
	Path   string
	Params map[string]interface{} `json:"-"`

	// IsEvaluated have the Params been applied to the Module?
	IsEvaluated bool `json:"-"`

	Module *Module
}

type Module struct {
	Path    string    `json:"-"`
	Imports []*Import `json:"Imports,omitempty"`

	// Raw yaml parsed as string or map[string]interface{}
	Raw map[string]interface{} `json:"Contents"`

	// IsEvaluated has the Raw yaml been evaluated (template holes filled, etc)?
	IsEvaluated bool `json:"-"`

	// TypedResult a koki object (or other object with special meaning)
	TypedResult interface{} `json:"-"`
}

type EvalContext struct {
	RawToTyped  func(raw interface{}) (interface{}, error)
	ApplyParams func(params map[string]interface{}, module *Module) error
}
