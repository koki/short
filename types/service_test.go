package types

import (
	"testing"

	"github.com/kr/pretty"
)

func TestServicePort(t *testing.T) {
	tryServicePort("8080:12345", t)
	tryServicePort("8080:12345:6789", t)
	tryServicePort("8080:12345:6789:TCP", t)
	tryServicePort("8080:12345:TCP", t)
	tryServicePort("8080:containerPortName:6789", t)
}

func tryServicePort(s string, t *testing.T) {
	sp := ServicePort{}
	err := sp.InitFromString(s)
	if err != nil {
		t.Error(err)
		return
	}

	ss := sp.String()
	if s != ss {
		t.Error(pretty.Sprintf("round-trip failed:\n(%s)\n(%# v)\n(%s)",
			s, sp, ss))
	}
}

func TestIngress(t *testing.T) {
	tryIngress("127.0.0.1", t)
	tryIngress("imahostname", t)
}

func tryIngress(s string, t *testing.T) {
	i := Ingress{}
	i.InitFromString(s)

	ss := i.String()
	if s != ss {
		t.Error(pretty.Sprintf("round-trip failed:\n(%s)\n(%# v)\n(%s)",
			s, i, ss))
	}
}
