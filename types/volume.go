package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/api/core/v1"

	"github.com/koki/short/util"
)

type VolumeWrapper struct {
	Volume Volume `json:"volume"`
}

type Volume struct {
	Volume v1.Volume
}

var volumeSourceLookup = map[string]string{
	"scale_io":          "scaleIO",
	"volume_id":         "volumeID",
	"storage_policy_id": "storagePolicyID",
	"pd_id":             "pdID",
	"disk_uri":          "diskURI",
	"downward_api":      "downwardAPI",
}

// TODO: HostPath Type (key collision) "type: hostPath.directory-or-create"

func (v *Volume) UnmarshalJSON(data []byte) error {
	obj := map[string]interface{}{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return util.InvalidValueForTypeErrorf(string(data), v, "couldn't deserialize")
	}

	if len(obj) < 2 {
		return util.InvalidValueForTypeErrorf(string(data), v, "expected at least two fields: name, type")
	}

	var volumeName string
	if name, ok := obj["name"]; ok {
		switch name := name.(type) {
		case string:
			volumeName = name
			delete(obj, "name")
		default:
			return util.InvalidValueForTypeErrorf(string(data), v, "expected 'name' field to be a string")
		}
	} else {
		return util.InvalidValueForTypeErrorf(string(data), v, "no 'name' field")
	}

	var volumeType string
	if vType, ok := obj["type"]; ok {
		switch vType := vType.(type) {
		case string:
			volumeType = vType
			delete(obj, "type")
		default:
			return util.InvalidValueForTypeErrorf(string(data), v, "expected 'type' field to be a string")
		}
	} else {
		return util.InvalidValueForTypeErrorf(string(data), v, "no 'type' field")
	}

	// EXTENSION POINT: Choose to do type-specific deserialization based on volumeType.
	volumeType, err = untweakVolumeSourceFields(volumeType, obj)
	if err != nil {
		return err
	}

	// Generic deserialization:
	volumeSourceObj := util.ConvertMapKeysToCamelCase(volumeSourceLookup, obj)
	volumeType = util.HyphenToCamelCase(volumeSourceLookup, volumeType)
	volumeObj := map[string]interface{}{
		"name":     volumeName,
		volumeType: volumeSourceObj,
	}

	b, err := json.Marshal(volumeObj)
	if err != nil {
		return util.InvalidValueForTypeErrorf(volumeObj, v, "couldn't reserialize Volume obj")
	}
	err = json.Unmarshal(b, &v.Volume)
	if err != nil {
		return util.InvalidValueForTypeErrorf(string(b), v, "couldn't deserialize")
	}

	return nil
}

func (v Volume) MarshalJSON() ([]byte, error) {
	var err error
	// EXTENSION POINT: Choose to do type-specific serialization based on the volume type.

	b, err := json.Marshal(v.Volume)
	if err != nil {
		return nil, err
	}

	obj := map[string]interface{}{}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return nil, err
	}

	if len(obj) != 2 {
		return nil, util.InvalidValueForTypeErrorf(obj, v, "should have two fields: 'name', volume-source")
	}

	var kokiName string
	if name, ok := obj["name"]; ok {
		switch name := name.(type) {
		case string:
			kokiName = name
		default:
			return nil, util.InvalidValueForTypeErrorf(name, v, "'name' field should be a string")
		}
		delete(obj, "name")
	} else {
		return nil, util.InvalidValueForTypeErrorf(obj, v, "expected a 'name' field")
	}

	var kokiType string
	var volumeSource map[string]interface{}
	for key, val := range obj {
		kokiType = util.CamelToHyphenCase(key)
		switch val := val.(type) {
		case map[string]interface{}:
			volumeSource = val
		default:
			return nil, util.InvalidValueForTypeErrorf(val, v, "VolumeSource wasn't reserialized as a dictionary")
		}
	}

	snakeObj := util.ConvertMapKeysToSnakeCase(volumeSource)
	kokiType, err = tweakVolumeSourceFields(kokiType, snakeObj)
	if err != nil {
		return nil, err
	}
	snakeObj["name"] = kokiName
	snakeObj["type"] = kokiType

	return json.Marshal(snakeObj)
}

func tweakVolumeSourceFields(kokiType string, obj map[string]interface{}) (string, error) {
	// Do type-specific tweaks here.
	switch kokiType {
	case "host-path":
		if val, ok := obj["type"]; ok {
			if str, ok := obj["type"].(string); ok {
				return fmt.Sprintf("%s.%s", kokiType, util.CamelToHyphenCase(str)), nil
			}
			return "", util.InvalidValueErrorf(val, "HostPath 'type' field should be a string")
		}
	case "config-map":
		if val, ok := obj["name"]; ok {
			if str, ok := obj["name"].(string); ok {
				obj["cm-name"] = str
				return kokiType, nil
			}
			return "", util.InvalidValueErrorf(val, "ConfigMap 'name' field should be a string")
		}
	}

	// Check for unhandled collisions.
	if _, ok := obj["name"]; ok {
		return "", util.InvalidValueErrorf(obj, "VolumeSource contains a 'name' field, which collides with the Volume 'name' field")
	}
	if _, ok := obj["type"]; ok {
		return "", util.InvalidValueErrorf(obj, "VolumeSource contains a 'type' field, which collides with the Volume 'type' field")
	}

	return kokiType, nil
}

func untweakVolumeSourceFields(kokiType string, obj map[string]interface{}) (string, error) {
	// Do type-specific untweaks here.
	if strings.HasPrefix(kokiType, "host-path.") {
		str := util.HyphenToCamelCase(volumeSourceLookup, kokiType[len("host-path."):])
		obj["type"] = strings.Title(str)
		return "host-path", nil
	}

	if kokiType == "config-map" {
		if name, ok := obj["cm-name"]; ok {
			obj["name"] = name
			delete(obj, "cm-name")
		}
		return kokiType, nil
	}

	return kokiType, nil
}
