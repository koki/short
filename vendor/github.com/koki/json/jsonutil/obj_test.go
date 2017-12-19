package jsonutil

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/koki/json"
)

type TestA struct {
	A        TestNull `json:"a,omitempty"`
	B        TestNull `json:"b,omitempty"`
	StringsA []string `json:"as,omitempty"`
	StringsB []string `json:"bs,omitempty"`
	StringsC []string `json:"cs,omitempty"`
	StringX  string   `json:"x,omitempty"`
	StringY  string   `json:"y,omitempty"`
	NumA     int      `json:"num_a,omitempty"`
	NumB     int      `json:"num_b,omitempty"`

	NestedB  TestB   `json:"nested_b,omitempty"`
	NestedBB TestB   `json:"nested_bb,omitempty"`
	NestedBs []TestB `json:"nested_bs,omitempty"`
}

type TestB struct {
	G int `json:"g,omitempty"`
}

type TestNull struct {
	Value int
}

var testA = TestA{
	A:        TestNull{Value: 0},
	B:        TestNull{Value: 1},
	StringsB: []string{},
	StringsC: []string{"c"},
	StringY:  "y",
	NumB:     1,

	NestedBB: TestB{G: 1},
	NestedBs: []TestB{TestB{}},
}

var stringA = `{
"a": null,
"b": 1,
"c": 2,
"as": null,
"bs": [],
"cs": ["c"],
"ds": ["d"],
"x": "",
"y": "y",
"z": "z",
"num_a": 0,
"num_b": 1,
"num_c": 2,
"nested_b": {},
"nested_bb": {"g": 1, "h": 2},
"nested_bs": [{"a": "asdf"}]
}`

var pathsA = [][]string{
	[]string{"c"},
	[]string{"ds"},
	[]string{"z"},
	[]string{"num_c"},
	[]string{"nested_bb", "h"},
	[]string{"nested_bs", "0", "a"},
}

func TestExtraneousFields(t *testing.T) {
	obj := map[string]interface{}{}
	err := json.Unmarshal([]byte(stringA), &obj)
	if err != nil {
		t.Fatal(err)
	}

	paths, err := ExtraneousFieldPaths(obj, testA)
	if err != nil {
		t.Fatal(err)
	}

	if !equalPaths(pathsA, paths) {
		t.Fatal("expected, actual", pathsA, paths)
	}
}

func sortedJoinedPaths(paths [][]string) []string {
	result := make([]string, len(paths))
	for i, path := range paths {
		result[i] = strings.Join(path, ".")
	}

	sort.Strings(result)
	return result
}

func equalPaths(expected, actual [][]string) bool {
	sortedExpected := sortedJoinedPaths(expected)
	sortedActual := sortedJoinedPaths(actual)
	return reflect.DeepEqual(sortedExpected, sortedActual)
}

func (t TestNull) MarshalJSON() ([]byte, error) {
	if t.Value == 0 {
		return []byte("null"), nil
	}

	return json.Marshal(t.Value)
}

func (t *TestNull) UnmarshalJSON(b []byte) error {
	if len(b) == 4 && string(b) == "null" {
		t.Value = 0
		return nil
	}

	return json.Unmarshal(b, &t.Value)
}
