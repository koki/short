package types

import (
	"reflect"
	"testing"

	"github.com/koki/json"
)

func TestRoleRef(t *testing.T) {
	testOneRoleRef("group.group1.kind:name", RoleRef{
		APIGroup: "group.group1",
		Kind:     "kind",
		Name:     "name",
	}, t, false)
	testOneRoleRef("kind:name", RoleRef{}, t, true)
	testOneRoleRef("name", RoleRef{}, t, true)
	testOneRoleRef("group.group1.kind", RoleRef{}, t, true)
	testOneRoleRef("group.:name", RoleRef{
		APIGroup: "group",
		Kind:     "",
		Name:     "name",
	}, t, false)
	// double backslash because it's a JSON string
	testOneRoleRef(`rbac.authorization.k8s.io.ClusterRole:kops\\:dns-controller`, RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     "kops:dns-controller",
	}, t, false)
}

func testOneRoleRef(nakedStr string, obj RoleRef, t *testing.T, decodeError bool) {
	str := `"` + nakedStr + `"`
	t.Log(str, obj)
	newObj := RoleRef{}
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

func TestSubject(t *testing.T) {
	testOneSubject("group.group1.kind:name", Subject{
		APIGroup: "group.group1",
		Kind:     "kind",
		Name:     "name",
	}, t, false)
	testOneSubject("group.group1.kind:namespace:name", Subject{
		APIGroup:  "group.group1",
		Kind:      "kind",
		Namespace: "namespace",
		Name:      "name",
	}, t, false)
	testOneSubject("kind:name", Subject{
		Kind: "kind",
		Name: "name",
	}, t, false)
	testOneSubject("kind:namespace:name", Subject{
		Kind:      "kind",
		Namespace: "namespace",
		Name:      "name",
	}, t, false)
	testOneSubject("name", Subject{}, t, true)
	testOneSubject("group.group1.kind", Subject{}, t, true)
	testOneSubject("group.:name", Subject{
		APIGroup: "group",
		Kind:     "",
		Name:     "name",
	}, t, false)
}

func testOneSubject(nakedStr string, obj Subject, t *testing.T, decodeError bool) {
	str := `"` + nakedStr + `"`
	t.Log(str, obj)
	newObj := Subject{}
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
