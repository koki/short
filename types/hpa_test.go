package types

import (
	"reflect"
	"testing"

	"github.com/koki/json"
)

func TestCrossVersionObjectReference(t *testing.T) {
	testOneCVORef("group.version.kind:name", CrossVersionObjectReference{
		APIVersion: "group.version",
		Kind:       "kind",
		Name:       "name",
	}, t, false)
	testOneCVORef("kind:name", CrossVersionObjectReference{
		Kind: "kind",
		Name: "name",
	}, t, false)
	testOneCVORef("name", CrossVersionObjectReference{}, t, true)
	testOneCVORef("group.group1.kind", CrossVersionObjectReference{}, t, true)
	testOneCVORef("group.:name", CrossVersionObjectReference{
		APIVersion: "group",
		Kind:       "",
		Name:       "name",
	}, t, false)
}

func testOneCVORef(nakedStr string, obj CrossVersionObjectReference, t *testing.T, decodeError bool) {
	str := `"` + nakedStr + `"`
	t.Log(str, obj)
	newObj := CrossVersionObjectReference{}
	err := json.Unmarshal([]byte(str), &newObj)
	if err != nil {
		if decodeError {
			return
		}

		t.Fatal(err)
	} else if decodeError {
		t.Fatal("expected a decode error")
	}

	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b), str, newObj, obj)
	if !reflect.DeepEqual(obj, newObj) {
		t.Fatal("objects don't match")
	}

	if str != string(b) {
		t.Fatal("strings don't match")
	}
}
