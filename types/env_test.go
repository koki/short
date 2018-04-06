package types

import (
	"testing"
)

func TestNewEnvWithValidInputs(t *testing.T) {
	e, err := NewEnv("key", "val")
	if err != nil {
		t.Errorf("error thrown on valid inputs by env builder fn NewEnv()")
	}

	if e.Type != EnvValEnvType {
		t.Errorf("Env type is wrongly assigned by env builder fn NewEnv()")
	}

	if e.From != nil {
		t.Errorf("env builder fn NewEnv() wrongly assigning From parameter")
	}

	if e.Val == nil {
		t.Errorf("env builder fn NewEnv() wrongly not assigning Val parameter")
	}

	if e.Val.Key != "key" || e.Val.Val != "val" {
		t.Errorf("Env key or value is wrongly assigned by env builder fn NewEnv()")
	}
}

func TestNewEnvWithInvalidInputs(t *testing.T) {
	_, err := NewEnv("", "val")
	if err == nil {
		t.Errorf("env builder fn NewEnv() does not detect invalid empty key input")
	}
}

func TestNewEnvFromWithValidInputs(t *testing.T) {
	typeList := []EnvFromType{
		EnvFromTypeCPULimits,
		EnvFromTypeMemLimits,
		EnvFromTypeEphemeralStorageLimits,
		EnvFromTypeCPURequests,
		EnvFromTypeMemRequests,
		EnvFromTypeEphemeralStorageRequests,
		EnvFromTypeMetadataName,
		EnvFromTypeMetadataNamespace,
		EnvFromTypeMetadataLabels,
		EnvFromTypeMetadataAnnotation,
		EnvFromTypeSpecNodename,
		EnvFromTypeSpecServiceAccountName,
		EnvFromTypeStatusHostIP,
		EnvFromTypeStatusPodIP,
	}

	for _, typ := range typeList {
		e, err := NewEnvFrom("key", typ)
		if err != nil {
			t.Errorf("env builder fn NewEnvFrom() wrongly detects error for from-type %s", typ)
		}

		if e.From.Key != "key" {
			t.Errorf("env builder fn NewEnvFrom() does not assign env key with from-type %s", typ)
		}

		if e.Type != EnvFromEnvType {
			t.Errorf("env builder fn NewEnvFrom() does not assign the right env type")
		}

		if e.Val != nil {
			t.Errorf("env builder fn NewEnvFrom() wrongly assigning Val parameter")
		}

		if e.From == nil {
			t.Errorf("env builder fn NewEnvFrom() wrongly not assigning From parameter")
		}

		if e.From.Key != "key" || *(e.From.Required) == true || e.From.From != string(typ) {
			t.Errorf("env builder fn NewEnvFrom() wrongly assigning Key, Required or From value")
		}
	}
}

func TestNewEnvFromWithInvalidInputs(t *testing.T) {
	_, err := NewEnvFrom("", EnvFromTypeStatusPodIP)
	if err == nil {
		t.Errorf("env builder fn NewEnvFrom() does not detect invalid empty key input")
	}

	_, err = NewEnvFrom("key", EnvFromTypeSecret)
	if err == nil {
		t.Errorf("env builder fn NewEnvFrom() does not detect invalid From type: %s", EnvFromTypeSecret)
	}

	_, err = NewEnvFrom("key", EnvFromTypeConfig)
	if err == nil {
		t.Errorf("env builder fn NewEnvFrom() does not detect invalid From type: %s", EnvFromTypeConfig)
	}
}

func TestNewEnvFromSecret(t *testing.T) {
	e, err := NewEnvFromSecret("key", "secretName", "secretKey")
	if err != nil {
		t.Errorf("env builder fn NewEnvFromSecret() wrongly detects error with valid key, secretName, and secretKey")
	}
	if e.From.From != "secret:secretName:secretKey" {
		t.Errorf("env builder fn NewEnvFromSecret() wrongly constructs (%s) environment string with valid key, secretName and secretKey", e.From.From)
	}

	e, err = NewEnvFromSecret("key", "secretName", "")
	if err != nil {
		t.Errorf("env builder fn NewEnvFromSecret() wrongly detects error with valid key, secretName, and empty secretKey")
	}

	if e.From.From != "secret:secretName" {
		t.Errorf("env builder fn NewEnvFromSecret() wrongly construct (%s) environment string with valid key, secretName and empty secretKey", e.From.From)
	}
}

func TestNewEnvFromConfig(t *testing.T) {
	e, err := NewEnvFromConfig("key", "configName", "configKey")
	if err != nil {
		t.Errorf("env builder fn NewEnvFromConfig() wrongly detects error with valid key, configName, and configKey")
	}
	if e.From.From != "config:configName:configKey" {
		t.Errorf("env builder fn NewEnvFromSecret() wrongly constructs (%s) environment string with valid key, configName and configKey", e.From.From)
	}

	e, err = NewEnvFromConfig("key", "configName", "")
	if err != nil {
		t.Errorf("env builder fn NewEnvFromConifg() wrongly detects error with valid key, configName, and empty configKey")
	}

	if e.From.From != "config:configName" {
		t.Errorf("env builder fn NewEnvFromConfig() wrongly construct (%s) environment string with valid key, configName and empty configKey", e.From.From)
	}
}

func TestNewEnvFromSecretOrConfigWithValidInputs(t *testing.T) {
	e, err := NewEnvFromSecretOrConfig(EnvFromTypeConfig, "key", "configName", "configKey")
	if err != nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly detects error with valid fromType, key, configName and configKey")
	}

	if e.Type != EnvFromEnvType {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns EnvType")
	}

	if e.Val != nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns Val parameter")
	}

	if e.From == nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly does not assign From parameter")
	}

	if *(e.From.Required) != true {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns Required value to false")
	}

	e, err = NewEnvFromSecretOrConfig(EnvFromTypeSecret, "key", "configName", "")
	if err != nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly detects error with valid fromType, key, configName and empty configKey")
	}

	if e.Type != EnvFromEnvType {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns EnvType")
	}

	if e.Val != nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns Val parameter")
	}

	if e.From == nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly does not assign From parameter")
	}

	if *(e.From.Required) != true {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns Required value to false")
	}
	e, err = NewEnvFromSecretOrConfig(EnvFromTypeSecret, "key", "secretName", "secretKey")
	if err != nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly detects error with valid fromType, key, secretName and secretKey")
	}

	if e.Type != EnvFromEnvType {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns EnvType")
	}

	if e.Val != nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns Val parameter")
	}

	if e.From == nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly does not assign From parameter")
	}

	if *(e.From.Required) != true {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns Required value to false")
	}

	e, err = NewEnvFromSecretOrConfig(EnvFromTypeSecret, "key", "secretName", "")
	if err != nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly detects error with valid fromType, key, secretName and empty secretKey")
	}

	if e.Type != EnvFromEnvType {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns EnvType")
	}

	if e.Val != nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns Val parameter")
	}

	if e.From == nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly does not assign From parameter")
	}

	if *(e.From.Required) != true {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() wrongly assigns Required value to false")
	}
}

func TestNewEnvFromSecretOrConfigWithInvalidInputs(t *testing.T) {
	_, err := NewEnvFromSecretOrConfig(EnvFromTypeSecret, "", "secretName", "secretKey")
	if err == nil {
		t.Errorf("env builder fn NewEnvFromSecretOrConfig() does not detect invalid empty key input")
	}

	typeList := []EnvFromType{
		EnvFromTypeCPULimits,
		EnvFromTypeMemLimits,
		EnvFromTypeEphemeralStorageLimits,
		EnvFromTypeCPURequests,
		EnvFromTypeMemRequests,
		EnvFromTypeEphemeralStorageRequests,
		EnvFromTypeMetadataName,
		EnvFromTypeMetadataNamespace,
		EnvFromTypeMetadataLabels,
		EnvFromTypeMetadataAnnotation,
		EnvFromTypeSpecNodename,
		EnvFromTypeSpecServiceAccountName,
		EnvFromTypeStatusHostIP,
		EnvFromTypeStatusPodIP,
	}

	for _, typ := range typeList {
		_, err := NewEnvFromSecretOrConfig(typ, "key", "resName", "resKey")
		if err == nil {
			t.Errorf("env builder fn NewEnvFromSecretOrConfig() does not detect invalid input from-type %s", typ)
		}
	}
}
