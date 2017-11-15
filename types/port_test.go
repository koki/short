package types

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"

	"github.com/ghodss/yaml"
	"github.com/kr/pretty"
)

var port0 = &Port{
	Name:          "port0",
	Protocol:      v1.ProtocolUDP,
	IP:            "1.2.3.4",
	HostPort:      "8080",
	ContainerPort: "80",
}

var port1 = &Port{
	Name:          "port1",
	Protocol:      v1.ProtocolUDP,
	HostPort:      "8080",
	ContainerPort: "80",
}

var port2 = &Port{
	Name:          "port2",
	Protocol:      v1.ProtocolTCP,
	HostPort:      "8080",
	ContainerPort: "80",
}

var port3 = &Port{
	Name:          "port3",
	Protocol:      v1.ProtocolTCP,
	IP:            "1.2.3.4",
	HostPort:      "8080",
	ContainerPort: "80",
}

var port4 = &Port{
	Protocol:      v1.ProtocolTCP,
	IP:            "1.2.3.4",
	ContainerPort: "80",
}

var port5 = &Port{
	Protocol:      v1.ProtocolTCP,
	ContainerPort: "80",
}

var port6 = &Port{
	Name:          "port6",
	Protocol:      v1.ProtocolTCP,
	ContainerPort: "80",
}

func TestPort(t *testing.T) {
	doTest("UDP://1.2.3.4:8080:80", t)
	doTest("1.2.3.4:8080:80", t)
	doTest("8080:80", t)

	doPortTest(port0, "port0: UDP://1.2.3.4:8080:80\n", t)
	doPortTest(port1, "port1: UDP://8080:80\n", t)
	doPortTest(port2, "port2: 8080:80\n", t)
	doPortTest(port3, "port3: 1.2.3.4:8080:80\n", t)
	doPortTest(port4, "1.2.3.4:80\n", t)
	doPortTest(port5, "80\n", t)
	doPortTest(port6, "port6: 80\n", t)
}

func doPortTest(port *Port, str string, t *testing.T) {
	b, err := yaml.Marshal(port)
	if err != nil {
		t.Error(pretty.Sprint(err.Error(), port))
	}

	if string(b) != str {
		t.Error(pretty.Sprint(str, string(b), port))
	}

	port1 := &Port{}
	err = yaml.Unmarshal(b, &port1)
	if err != nil {
		t.Error(pretty.Sprint(err.Error(), port, string(b)))
	}

	if !reflect.DeepEqual(port, port1) {
		t.Error(port, port1, string(b))
	}
}

func doTest(str string, t *testing.T) {
	p := Port{}
	err := p.InitFromString(str)
	if err != nil {
		t.Error(pretty.Sprint(err.Error(), str))
	}

	str1, err := p.ToString()
	if err != nil {
		t.Error(pretty.Sprint(err.Error(), str, p))
	}

	if str != str1 {
		t.Error(pretty.Sprint(str, str1, p))
	}
}
