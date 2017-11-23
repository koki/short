package imports

type Import struct {
	Name   string
	Path   string
	Params map[string]interface{} `json:"Params,omitempty"`

	// IsEvaluated have the Params been applied to the Module?
	IsEvaluated bool `json:"-"`

	Module *Module
}

type ParamDef struct {
	Description string
	Default     interface{}
}

type Module struct {
	Path    string              `json:"-"`
	Imports []*Import           `json:"Imports,omitempty"`
	Params  map[string]ParamDef `json:"Params,omitempty"`

	// Raw yaml parsed as string or map[string]interface{}
	Raw map[string]interface{} `json:"Contents"`

	// IsEvaluated has the Raw yaml been evaluated (template holes filled, etc)?
	IsEvaluated bool `json:"-"`

	// TypedResult a koki object (or other object with special meaning)
	TypedResult interface{} `json:"-"`
}

type EvalContext struct {
	RawToTyped func(raw interface{}) (interface{}, error)
}
