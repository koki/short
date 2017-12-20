package template

import (
	"strings"
	"testing"

	"github.com/koki/short/yaml"
)

var template0 = `
list:
- a
- ${key...}
- d
`

var params0 = `
key:
- b
- c
`

var result0 = `
list:
- a
- b
- c
- d
`

var template1 = `
list: ${key2}
string: a${key1}d
`

var params1 = `
key1: 12
key2:
- b
- c
`

var result1 = `
list:
- b
- c
string: a12d
`

func TestTemplate(t *testing.T) {
	doTest(template0, params0, result0, t)
	doTest(template1, params1, result1, t)
}

func doTest(template, params, result string, t *testing.T) {
	template = strings.Trim(template, "\n")
	params = strings.Trim(params, "\n")
	result = strings.Trim(result, "\n")

	t.Log(template)
	t.Log(params)

	templateObj := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(template), &templateObj)
	if err != nil {
		t.Error(err)
	}

	paramsObj := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(params), &paramsObj)
	if err != nil {
		t.Error(err)
	}

	resultObj, err := ReplaceMap(templateObj, ResolverForParams(paramsObj))
	if err != nil {
		t.Error(err)
	}

	b, err := yaml.Marshal(resultObj)
	if err != nil {
		t.Error(err)
	}
	resultStr := strings.Trim(string(b), "\n")
	t.Log(resultStr)
	t.Log(result)

	if resultStr != result {
		t.Error("unexpected result")
	}
}
