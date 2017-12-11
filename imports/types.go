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

type Resource struct {
	// Raw yaml parsed as string or map[string]interface{}
	Raw map[string]interface{} `json:"Value"`

	// TypedResult a koki object (or other object with special meaning)
	TypedResult interface{} `json:"-"`
}

type Module struct {
	Path    string              `json:"-"`
	Imports []*Import           `json:"Imports,omitempty"`
	Params  map[string]ParamDef `json:"Params,omitempty"`

	// IsEvaluated has the Raw yaml in Exports been evaluated (template holes filled, etc)?
	IsEvaluated bool `json:"-"`

	Export Resource
}

type EvalContext struct {
	RawToTyped func(raw interface{}) (interface{}, error)

	// Get an "absolute" path for a given import.
	ResolveImportPath func(rootPath string, importPath string) (string, error)

	// Read the contents of a given path.
	ReadFromPath func(path string) ([]map[string]interface{}, error)
}
