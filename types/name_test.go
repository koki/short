package types

import (
	"reflect"
	"testing"
)

func TestNameEscaping(t *testing.T) {
	testOneName("", "", t)
	testOneName(`\`, `\\`, t)
	testOneName(`:`, `\:`, t)
	testOneName(`\:`, `\\\:`, t)
	testOneName(`:\`, `\:\\`, t)
	testOneName(`system:serviceaccount:kube-system:dns-controller`,
		`system\:serviceaccount\:kube-system\:dns-controller`, t)
	testOneName(`bob`, `bob`, t)
	testOneName(`bob\`, `bob\\`, t)
}

func testOneName(unescaped Name, escaped string, t *testing.T) {
	actualEscaped := EscapeName(unescaped)
	actualUnescaped := UnescapeName(escaped)
	if actualEscaped != escaped || actualUnescaped != unescaped {
		t.Fatalf("actual:\n  %s\n  %s\nexpected:\n  %s\n  %s\n", actualEscaped, actualUnescaped, escaped, unescaped)
	}
}

func TestSplitAtUnescapedColons(t *testing.T) {
	testOneSplit("", t, "")
	testOneSplit(":", t, "", "")
	testOneSplit(`\:`, t, `\:`)
	testOneSplit(`\\:`, t, `\\`, "")
	testOneSplit(`a:b:kops\:dns-controller`, t, `a`, `b`, `kops\:dns-controller`)
	testOneSplit(`a:b:c`, t, `a`, `b`, `c`)
}

func testOneSplit(s string, t *testing.T, segments ...string) {
	actualSegments := SplitAtUnescapedColons(s)
	if !reflect.DeepEqual(actualSegments, segments) {
		t.Fatalf("actual:\n  %#v\nexpected:\n  %#v", actualSegments, segments)
	}
}
