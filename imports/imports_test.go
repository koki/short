package imports

import (
	"fmt"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/kr/pretty"
)

var modules = map[string]string{
	"module0": `
resource0: stuff
`,
	"module1": `
exports:
- default: stuff
  value: thing
`,
	"module2": `
imports:
- import1: module1
exports:
- default: stuff
  value: thing
`,
	"module3": `
params:
- param0: a param
exports:
- default: stuff
  value: thing
`,
	"module4": `
imports:
- import2: module2
- import3: module3
  params:
    param0: stuff
params:
- param0: a param
- param1: another param
  default: 123
exports:
- export0: stuff
  value: thing
- export1: another
  value: things
`,
	"module5": `
imports:
  import1: module1
exports:
- default: stuff
  value: thing
`,
	"module6": `
params:
  param0: a param
exports:
- default: stuff
  value: thing
`,
	"module7": `
params:
- param0: a param
export0: ${param0}
`,
}

func getEvalContext(t *testing.T) *EvalContext {
	return &EvalContext{
		ResolveImportPath: func(rootPath string, importPath string) (string, error) {
			return importPath, nil
		},
		ReadFromPath: func(path string) ([]map[string]interface{}, error) {
			if contents, ok := modules[path]; ok {
				obj := map[string]interface{}{}
				err := yaml.Unmarshal([]byte(contents), &obj)
				if err != nil {
					t.Error(err)
				}
				return []map[string]interface{}{obj}, nil
			}
			return nil, fmt.Errorf("no module (%s)", path)
		},
	}
}

func TestImports(t *testing.T) {
	doTestImport("module0", t, false)
	doTestImport("module1", t, false)
	doTestImport("module2", t, false)
	doTestImport("module3", t, false)
	doTestImport("module4", t, false)
	doTestImport("module5", t, true)
	doTestImport("module6", t, true)
	doTestImport("module7", t, true)
}

func doTestImport(modulePath string, t *testing.T, expectParseError bool) {
	evalContext := getEvalContext(t)
	modules, err := evalContext.Parse(modulePath)
	if err != nil {
		if !expectParseError {
			t.Fatal(err)
		}

		return
	} else if expectParseError {
		t.Fatal(pretty.Sprintf("unexpected success evaluating module\n(%# v)", modules))
	}

	if len(modules) != 1 {
		t.Error(pretty.Sprintf("expected only one module\n%# v", modules))
	}
}
