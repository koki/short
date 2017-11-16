package types

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/kr/pretty"
)

var nsp0 = "name0: 80\n"
var nsp1 = "name1: 80:1234\n"
var nsp2 = "name2: 80:1234\nnode_port: 3789\n"

func TestNamedServicePort(t *testing.T) {
	tryNamedServicePort(nsp0, t)
	tryNamedServicePort(nsp1, t)
	tryNamedServicePort(nsp2, t)
}

func tryNamedServicePort(s string, t *testing.T) {
	nsp := NamedServicePort{}
	err := yaml.Unmarshal([]byte(s), &nsp)
	if err != nil {
		t.Error(pretty.Sprintf("%s:\n(%s)",
			err.Error(), s))
		return
	}

	b, err := yaml.Marshal(nsp)
	if err != nil {
		t.Error(pretty.Sprintf("%s:\n(%s)\n(%# v)",
			err.Error(), s, nsp))
		return
	}

	if s != string(b) {
		t.Error(pretty.Sprintf("round-trip failed:\n(%s)\n(%# v)\n(%s)",
			s, nsp, string(b)))
	}
}

func TestServicePort(t *testing.T) {
	tryServicePort("8080\n", t)
	tryServicePort("8080:12345\n", t)
	tryServicePort("UDP://8080:12345\n", t)
	tryServicePort("8080:containerPortName\n", t)
}

func tryServicePort(s string, t *testing.T) {
	sp := ServicePort{}
	err := yaml.Unmarshal([]byte(s), &sp)
	if err != nil {
		t.Error(pretty.Sprintf("%s:\n(%s)",
			err.Error(), s))
		return
	}

	b, err := yaml.Marshal(sp)
	if err != nil {
		t.Error(pretty.Sprintf("%s:\n(%s)\n(%# v)",
			err.Error(), s, sp))
		return
	}

	if s != string(b) {
		t.Error(pretty.Sprintf("round-trip failed:\n(%s)\n(%# v)\n(%s)",
			s, sp, string(b)))
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
