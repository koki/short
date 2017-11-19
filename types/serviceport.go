package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/util"
)

type NamedServicePort struct {
	Name     string
	Port     ServicePort
	NodePort int32
}

type ServicePort struct {
	Expose int32

	// PodPort is a port or the name of a containerPort.
	PodPort *intstr.IntOrString

	// Protocol is optional. "" is empty.
	Protocol v1.Protocol
}

func (p *ServicePort) InitFromInt(i int32) {
	p.Protocol = v1.ProtocolTCP
	p.Expose = i
}

func (p *ServicePort) InitFromString(str string) error {
	matches := protocolPortRegexp.FindStringSubmatch(str)

	// Extract the Protocol first.
	if len(matches) > 0 {
		p.Protocol = v1.Protocol(matches[1])
		str = matches[2]
	} else {
		p.Protocol = v1.ProtocolTCP
	}

	segments := strings.Split(str, ":")
	l := len(segments)
	if l < 1 {
		return util.InvalidValueForTypeErrorf(str, p, "too few sections")
	}
	if l > 2 {
		return util.InvalidValueForTypeErrorf(str, p, "too many sections")
	}

	// Extract the exposed port, which is the only required field.
	expose, err := strconv.ParseInt(segments[0], 10, 32)
	if err != nil {
		return util.InvalidValueForTypeErrorf(str, p, "couldn't parse exposed service port")
	}
	p.Expose = int32(expose)

	// Extract the Pod/Container Port if it exists.
	if l > 1 {
		p.PodPort = util.IntOrStringPtr(intstr.Parse(segments[1]))
	}

	return nil
}

func (p *ServicePort) String() string {
	str := fmt.Sprintf("%d", p.Expose)
	if p.PodPort != nil {
		str = fmt.Sprintf("%s:%s", str, p.PodPort.String())
	}

	if len(p.Protocol) == 0 || p.Protocol == v1.ProtocolTCP {
		// No need to specify protocol
		return str
	}

	return fmt.Sprintf("%s://%s", p.Protocol, str)
}

func (p *ServicePort) ToInt() (int32, error) {
	if len(p.Protocol) == 0 || p.Protocol == v1.ProtocolTCP {
		if p.PodPort == nil {
			return p.Expose, nil
		}
	}
	return -1, util.InvalidInstanceErrorf(p, "can't serialize as int32")
}

func (p *ServicePort) UnmarshalJSON(data []byte) error {
	var i int32
	intErr := json.Unmarshal(data, &i)
	if intErr == nil {
		p.InitFromInt(i)
		return nil
	}

	var s string
	strErr := json.Unmarshal(data, &s)
	if strErr != nil {
		return util.InvalidValueForTypeErrorf(string(data), p, "couldn't unmarshal JSON as int or string: (%s), (%s)", intErr.Error(), strErr.Error())
	}

	return p.InitFromString(s)
}

func (p ServicePort) MarshalJSON() ([]byte, error) {
	i, intErr := p.ToInt()
	if intErr == nil {
		b, err := json.Marshal(i)
		if err != nil {
			return nil, util.InvalidInstanceErrorf(p, "couldn't marshal port number (%d) to JSON: %s", i, err.Error())
		}

		return b, nil
	}

	str := p.String()
	b, strErr := json.Marshal(str)
	if strErr != nil {
		return nil, util.InvalidInstanceErrorf(p, "couldn't marshal to JSON from string (%s): %s", str, strErr.Error())
	}

	return b, nil
}

func (n *NamedServicePort) InitFromMap(obj map[string]interface{}) error {
	if len(obj) > 2 {
		return util.InvalidValueForTypeErrorf(obj, n, "expected at most 2 fields")
	}

	for key, val := range obj {
		if key == "node_port" {
			if val, ok := val.(float64); ok {
				n.NodePort = int32(val)
			}
		} else {
			n.Name = key
			switch val := val.(type) {
			case string:
				n.Port.InitFromString(val)
			case float64:
				n.Port.InitFromInt(int32(val))
			default:
				return util.InvalidValueForTypeErrorf(obj, n, "expected string or int for ServicePort")
			}
		}
	}

	return nil
}

func (n *NamedServicePort) UnmarshalJSON(data []byte) error {
	var obj = map[string]interface{}{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return util.InvalidValueForTypeErrorf(string(data), n, "couldn't deserialize")
	}

	return n.InitFromMap(obj)
}

func (n NamedServicePort) MarshalJSON() ([]byte, error) {
	var obj = map[string]interface{}{}
	i, err := n.Port.ToInt()
	if err == nil {
		obj[n.Name] = i
	} else {
		obj[n.Name] = n.Port.String()
	}

	if n.NodePort > 0 {
		obj["node_port"] = n.NodePort
	}

	return json.Marshal(obj)
}
