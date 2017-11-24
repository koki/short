package imports

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/kr/pretty"
)

var evalModules = map[string]string{
	"module0": `
resource0: stuff
`,
	"module1": `
params:
- param0: a param
exports:
- default: a resource
  value: ${param0}
`,
	"module2": `
imports:
- import0: module1
  params:
    param0: 1234
exports:
- default: a resource
  value: ${import0}
`,
	"module3": `
imports:
- import0: module1
exports:
- default: a resource
  value: ${import0}
`,
	"module4": `
imports:
- import0: module0
- import2: module2
exports:
- export0: first export
  value: ${import0}
- export1: second export
  value: ${import2}
- export2: third export
  value: something
`,
}

var evalResults = map[string]interface{}{
	"module0": map[string]interface{}{
		"default": map[string]interface{}{
			"resource0": "stuff",
		},
	},
	"module2": map[string]interface{}{
		"default": float64(1234),
	},
	"module4": map[string]interface{}{
		"export0": map[string]interface{}{
			"resource0": "stuff",
		},
		"export1": float64(1234),
		"export2": "something",
	},
}

func getFullEvalContext(t *testing.T) *EvalContext {
	return &EvalContext{
		RawToTyped: func(raw interface{}) (interface{}, error) {
			return raw, nil
		},
		ResolveImportPath: func(rootPath string, importPath string) (string, error) {
			return importPath, nil
		},
		ReadFromPath: func(path string) ([]map[string]interface{}, error) {
			if contents, ok := evalModules[path]; ok {
				obj := map[string]interface{}{}
				err := yaml.Unmarshal([]byte(contents), &obj)
				if err != nil {
					t.Fatal(err)
				}
				return []map[string]interface{}{obj}, nil
			}
			return nil, fmt.Errorf("no module (%s)", path)
		},
	}
}

func TestEval(t *testing.T) {
	doTestEval("module0", t, false)
	doTestEval("module1", t, true)
	doTestEval("module2", t, false)
	doTestEval("module3", t, true)
	doTestEval("module4", t, false)
}

func doTestEval(modulePath string, t *testing.T, expectEvalError bool) {
	evalContext := getFullEvalContext(t)
	modules, err := evalContext.Parse(modulePath)
	if err != nil {
		t.Fatal(err)
	}

	if len(modules) != 1 {
		t.Fatal(pretty.Sprintf("expected only one module\n(%# v)", modules))
	}

	module := &modules[0]
	err = evalContext.EvaluateModule(module, nil)
	if err != nil {
		if !expectEvalError {
			t.Fatal(err)
		}

		return
	}
	if err == nil && expectEvalError {
		t.Fatal(pretty.Sprintf("unexpected success evaluating module\n(%# v)", module))
		return
	}

	exports := map[string]interface{}{}
	for name, export := range module.Exports {
		exports[name] = export.Raw
	}

	if !reflect.DeepEqual(exports, evalResults[modulePath]) {
		t.Fatal(pretty.Sprintf("evaluated module doesn't match expected\n(%# v)\n(%# v)", exports, evalResults[modulePath]))
	}
}
