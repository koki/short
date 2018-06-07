package types

import (
	"fmt"
	"strings"

	"github.com/koki/json"
	util "github.com/koki/structurederrors"
)

type EnvFromType string

const (
	EnvFromTypeSecret EnvFromType = "secret"
	EnvFromTypeConfig EnvFromType = "config"

	EnvFromTypeCPULimits              EnvFromType = "limits.cpu"
	EnvFromTypeMemLimits              EnvFromType = "limits.memory"
	EnvFromTypeEphemeralStorageLimits EnvFromType = "limits.ephemeral-storage"

	EnvFromTypeCPURequests              EnvFromType = "requests.cpu"
	EnvFromTypeMemRequests              EnvFromType = "requests.memory"
	EnvFromTypeEphemeralStorageRequests EnvFromType = "requests.ephemeral-storage"

	EnvFromTypeMetadataName       EnvFromType = "metadata.name"
	EnvFromTypeMetadataNamespace  EnvFromType = "metadata.namespace"
	EnvFromTypeMetadataLabels     EnvFromType = "metadata.labels"
	EnvFromTypeMetadataAnnotation EnvFromType = "metadata.annotations"

	EnvFromTypeSpecNodename           EnvFromType = "spec.nodeName"
	EnvFromTypeSpecServiceAccountName EnvFromType = "spec.serviceAccountName"

	EnvFromTypeStatusHostIP EnvFromType = "status.hostIP"
	EnvFromTypeStatusPodIP  EnvFromType = "status.podIP"
)

func NewEnv(key, val string) (Env, error) {
	if key == "" {
		return Env{}, fmt.Errorf("Env key cannnot be empty")
	}
	return Env{
		Type: EnvValEnvType,
		Val: &EnvVal{
			Key: key,
			Val: val,
		},
	}, nil
}

func NewEnvFrom(key string, from EnvFromType) (Env, error) {
	if key == "" {
		return Env{}, fmt.Errorf("Env key cannot be empty")
	}
	if from == EnvFromTypeConfig || from == EnvFromTypeSecret {
		return Env{}, fmt.Errorf("%s not supported. Use NewEnvFromSecret() or NewEnvFromConfig() for building new envs from Secret or ConfigMap resources", from)
	}
	required := false
	return Env{
		Type: EnvFromEnvType,
		From: &EnvFrom{
			Key:      key,
			From:     string(from),
			Required: &required,
		},
	}, nil
}

func NewEnvFromSecretOrConfig(resType EnvFromType, prefix, resName, resKey string) (Env, error) {
	if resType != EnvFromTypeSecret && resType != EnvFromTypeConfig {
		return Env{}, fmt.Errorf("%s not supported. Use NewEnvFrom() for building new envs from resources other than Secret or ConfigMap resources", resType)
	}
	format := fmt.Sprintf("%s", resType)
	fromVal := ""
	if resKey != "" {
		format = format + ":%s:%s"
		fromVal = fmt.Sprintf(format, resName, resKey)
	} else {
		format = format + ":%s"
		fromVal = fmt.Sprintf(format, resName)
	}
	required := true
	return Env{
		Type: EnvFromEnvType,
		From: &EnvFrom{
			Key:      prefix,
			From:     fromVal,
			Required: &required,
		},
	}, nil
}

func NewEnvFromSecret(key, secretName, secretKey string) (Env, error) {
	return NewEnvFromSecretOrConfig(EnvFromTypeSecret, key, secretName, secretKey)
}

func NewEnvFromConfig(key, configName, configKey string) (Env, error) {
	return NewEnvFromSecretOrConfig(EnvFromTypeConfig, key, configName, configKey)
}

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
	EnvFromEnvType EnvType = iota
	EnvValEnvType
)

func (e EnvFrom) Optional() *bool {
	if e.Required == nil {
		return nil
	}

	optional := !(*e.Required)
	return &optional
}

func (e *Env) SetVal(val EnvVal) {
	e.Type = EnvValEnvType
	e.Val = &val
}

func (e *Env) SetFrom(from EnvFrom) {
	e.Type = EnvFromEnvType
	e.From = &from
}

func EnvWithVal(val EnvVal) Env {
	return Env{
		Type: EnvValEnvType,
		Val:  &val,
	}
}

func EnvWithFrom(from EnvFrom) Env {
	return Env{
		Type: EnvFromEnvType,
		From: &from,
	}
}

func ParseEnvVal(s string) *EnvVal {
	segments := strings.SplitN(s, "=", 2)
	if len(segments) == 2 {
		return &EnvVal{
			Key: segments[0],
			Val: segments[1],
		}
	}

	// Interpret the entire string as the variable name.
	return &EnvVal{
		Key: s,
	}
}

func UnparseEnvVal(val EnvVal) string {
	if len(val.Val) == 0 {
		return val.Key
	}

	return fmt.Sprintf("%s=%s", val.Key, val.Val)
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (e *Env) UnmarshalJSON(value []byte) error {
	var s string
	err := json.Unmarshal(value, &s)
	if err == nil {
		envVal := ParseEnvVal(s)
		e.SetVal(*envVal)
		return nil
	}

	from := EnvFrom{}
	err = json.Unmarshal(value, &from)
	if err == nil {
		e.SetFrom(from)
		return nil
	}

	return util.InvalidInstanceErrorf(e, "couldn't parse from (%s)", string(value))
}

// MarshalJSON implements the json.Marshaller interface.
func (e Env) MarshalJSON() ([]byte, error) {
	var b []byte
	var err error
	switch e.Type {
	case EnvValEnvType:
		b, err = json.Marshal(UnparseEnvVal(*e.Val))
	case EnvFromEnvType:
		b, err = json.Marshal(e.From)
	default:
		return []byte{}, util.InvalidInstanceError(e.Type)
	}

	if err != nil {
		return nil, util.InvalidInstanceContextErrorf(err, e, "marshalling to JSON")
	}

	return b, nil
}
