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
value: ${param0}
`,
	"module2": `
imports:
- import0: module1
  params:
    param0: 1234
value: ${import0}
`,
	"module3": `
imports:
- import0: module1
value: ${import0}
`,
	"module4": `
imports:
- import0: module0
- import2: module2
value:
- ${import0}
- ${import2}
- something
`,
	"module5": `
params:
- param0: a param with a default value
  default: x
value: ${param0}
`,
	"module6": `
imports:
- import7: module7
value:
- ${import7.blah}
- ${import7.doot.0}
- something
`,
	"module7": `
thing:
  blah: "bleh"
  doot:
  - what: hello
  - not: this
`,
}

var evalResults = map[string]interface{}{
	"module0": map[string]interface{}{
		"resource0": "stuff",
	},
	"module2": map[string]interface{}{
		"value": float64(1234),
	},
	"module4": map[string]interface{}{
		"value": []interface{}{
			"stuff",
			float64(1234),
			"something",
		},
	},
	"module5": map[string]interface{}{
		"value": "x",
	},
	"module6": map[string]interface{}{
		"value": []interface{}{
			"bleh",
			map[string]interface{}{
				"what": "hello",
			},
			"something",
		},
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
	doTestEval("module5", t, false)
	doTestEval("module6", t, false)
}

func doTestEval(modulePath string, t *testing.T, expectEvalError bool) {
	evalContext := getFullEvalContext(t)
	modules, err := evalContext.Parse(modulePath)
	if err != nil {
		t.Fatal(err)
	}

	if len(modules) != 1 {
		t.Fatal(pretty.Sprintf("%s: expected only one module\n(%# v)", modulePath, modules))
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
		t.Fatal(pretty.Sprintf("%s: unexpected success evaluating module\n(%# v)", modulePath, module))
		return
	}

	if !reflect.DeepEqual(module.Export.Raw, evalResults[modulePath]) {
		t.Fatal(pretty.Sprintf("%s: evaluated module doesn't match expected\n(%# v)\n(%# v)", modulePath, module.Export.Raw, evalResults[modulePath]))
	}
}
