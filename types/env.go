package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/koki/short/util"
)

type EnvFrom struct {
	Key      string `json:"key,omitempty"`
	From     string `json:"from,omitempty"`
	Required *bool  `json:"required,omitempty"`
}

type EnvVal struct {
	Key string
	Val string
}

type Env struct {
	Type EnvType
	From *EnvFrom
	Val  *EnvVal
}

type EnvType int

const (
	EnvFromType EnvType = iota
	EnvValType
)

func (e *Env) SetVal(val EnvVal) {
	e.Type = EnvValType
	e.Val = &val
}

func (e *Env) SetFrom(from EnvFrom) {
	e.Type = EnvFromType
	e.From = &from
}

func EnvWithVal(val EnvVal) Env {
	return Env{
		Type: EnvValType,
		Val:  &val,
	}
}

func EnvWithFrom(from EnvFrom) Env {
	return Env{
		Type: EnvFromType,
		From: &from,
	}
}

func ParseEnvVal(s string) (*EnvVal, error) {
	segments := strings.Split(s, "=")
	if len(segments) == 2 {
		return &EnvVal{
			Key: segments[0],
			Val: segments[1],
		}, nil
	}

	return nil, fmt.Errorf("unrecognized EnvVal (%s)", s)
}

func UnparseEnvVal(val EnvVal) string {
	return fmt.Sprintf("%s=%s", val.Key, val.Val)
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (e *Env) UnmarshalJSON(value []byte) error {
	var s string
	err := json.Unmarshal(value, &s)
	if err == nil {
		envVal, err := ParseEnvVal(s)
		if err != nil {
			return fmt.Errorf("unrecognized Env (%s)", string(value))
		}
		e.SetVal(*envVal)
		return nil
	}

	from := EnvFrom{}
	err = json.Unmarshal(value, &from)
	if err == nil {
		e.SetFrom(from)
		return nil
	}

	return util.PrettyTypeError(e, string(value))
}

// MarshalJSON implements the json.Marshaller interface.
func (e Env) MarshalJSON() ([]byte, error) {
	switch e.Type {
	case EnvValType:
		return json.Marshal(UnparseEnvVal(*e.Val))
	case EnvFromType:
		return json.Marshal(e.From)
	default:
		return []byte{}, fmt.Errorf("impossible Env.Type")
	}
}
