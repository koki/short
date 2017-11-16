package types

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/api/core/v1"

	"github.com/koki/short/util"
)

type Port struct {
	Name          string
	Protocol      v1.Protocol
	IP            string
	HostPort      string
	ContainerPort string
}

func (p *Port) HostPortInt() (int32, error) {
	if len(p.HostPort) > 0 {
		hostPort, err := strconv.ParseInt(p.HostPort, 10, 32)
		if err != nil {
			return 0, err
		}

		return int32(hostPort), nil
	}

	return 0, nil
}

func (p *Port) ContainerPortInt() (int32, error) {
	if len(p.ContainerPort) > 0 {
		containerPort, err := strconv.ParseInt(p.ContainerPort, 10, 32)
		if err != nil {
			return 0, err
		}

		return int32(containerPort), nil
	}

	return 0, nil
}

/*
$protocol://$ip:$host_port:$container_port

expose:
  - 8080:80
  - UDP://127.0.0.1:8080:80
  - 10.10.0.53:8081:9090
  - port_name: 192.168.1.2:8090:80
*/

var protocolPortRegexp = regexp.MustCompile(`^(UDP|TCP)://([0-9.:]*)$`)

func (p *Port) InitFromString(str string) error {
	matches := protocolPortRegexp.FindStringSubmatch(str)
	if len(matches) > 0 {
		p.Protocol = v1.Protocol(matches[1])
		str = matches[2]
	} else {
		p.Protocol = v1.ProtocolTCP
	}

	segments := strings.Split(str, ":")
	parseIndex := 0

	ip := net.ParseIP(segments[parseIndex])
	if ip != nil {
		p.IP = segments[parseIndex]
		parseIndex++
	}

	remaining := len(segments) - parseIndex
	if remaining == 2 {
		p.HostPort = segments[parseIndex]
		p.ContainerPort = segments[parseIndex+1]
		return nil
	}
	if remaining == 1 {
		p.ContainerPort = segments[parseIndex]
		return nil
	}

	return util.InvalidInstanceErrorf(p, "couldn't parse (%s)", str)
}

func appendColonSegment(str, seg string) string {
	if len(str) == 0 {
		return seg
	}

	return fmt.Sprintf("%s:%s", str, seg)
}

func (p *Port) ToString() (string, error) {
	str := ""
	if len(p.IP) > 0 {
		str = appendColonSegment(str, p.IP)
	}

	if len(p.HostPort) > 0 {
		str = appendColonSegment(str, p.HostPort)
	}

	if len(p.ContainerPort) > 0 {
		str = appendColonSegment(str, p.ContainerPort)
	}

	if len(p.Protocol) == 0 || p.Protocol == v1.ProtocolTCP {
		// No need to specify protocol
		return str, nil
	}

	return fmt.Sprintf("%s://%s", p.Protocol, str), nil
}

func (p *Port) UnmarshalJSON(data []byte) error {
	i := 0
	err := json.Unmarshal(data, &i)
	if err == nil {
		return p.InitFromString(fmt.Sprintf("%d", i))
	}

	str := ""
	err = json.Unmarshal(data, &str)
	if err == nil {
		return p.InitFromString(str)
	}

	obj := map[string]interface{}{}
	err = json.Unmarshal(data, &obj)

	if len(obj) != 1 {
		return util.InvalidValueErrorf(obj, "expected only one entry for Port")
	}

	if err == nil {
		for key, val := range obj {
			p.Name = key
			switch val := val.(type) {
			case string:
				err = p.InitFromString(val)
			case float64:
				err = p.InitFromString(fmt.Sprintf("%d", int(val)))
			default:
				err = util.InvalidValueErrorf(obj, "unrecognized value (not a string or number)")
			}
		}
	}

	return err
}

func (p Port) MarshalJSON() ([]byte, error) {
	str, err := p.ToString()
	if err != nil {
		return nil, err
	}
	i, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		// It's just a ContainerPort
		if len(p.Name) > 0 {
			obj := map[string]int{
				p.Name: int(i),
			}
			return json.Marshal(&obj)
		}

		return json.Marshal(&i)
	}

	if len(p.Name) > 0 {
		obj := map[string]string{
			p.Name: str,
		}
		return json.Marshal(&obj)
	}

	return json.Marshal(&str)
}
