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
	VolumeMeta
	VolumeSource
}

type VolumeMeta struct {
	Name string `json:"name"`
}

type VolumeSource struct {
	VolumeSource v1.VolumeSource
}

var volumeSourceLookup = map[string]string{
	"scale_io":          "scaleIO",
	"volume_id":         "volumeID",
	"storage_policy_id": "storagePolicyID",
	"pd_id":             "pdID",
	"disk_uri":          "diskURI",
	"downward_api":      "downwardAPI",
}

func (v *Volume) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &v.VolumeSource)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &v.VolumeMeta)
}

func (v Volume) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(v.VolumeMeta)
	if err != nil {
		return nil, err
	}

	bb, err := v.VolumeSource.MarshalJSON()
	if err != nil {
		return nil, err
	}

	metaObj := map[string]interface{}{}
	err = json.Unmarshal(b, &metaObj)
	if err != nil {
		return nil, err
	}

	sourceObj := map[string]interface{}{}
	err = json.Unmarshal(bb, &sourceObj)
	if err != nil {
		return nil, err
	}

	// Merge metadata with volume-source
	for key, val := range metaObj {
		sourceObj[key] = val
	}

	return json.Marshal(sourceObj)
}

func PreprocessVolumeSourceJSON(v interface{}, data []byte) ([]byte, error) {
	obj := map[string]interface{}{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return nil, util.InvalidValueForTypeErrorf(string(data), v, "couldn't deserialize")
	}

	if len(obj) < 2 {
		return nil, util.InvalidValueForTypeErrorf(string(data), v, "expected at least two fields: name, type")
	}

	var volumeType string
	if vType, ok := obj["type"]; ok {
		switch vType := vType.(type) {
		case string:
			volumeType = vType
			delete(obj, "type")
		default:
			return nil, util.InvalidValueForTypeErrorf(string(data), v, "expected 'type' field to be a string")
		}
	} else {
		return nil, util.InvalidValueForTypeErrorf(string(data), v, "no 'type' field")
	}

	// EXTENSION POINT: Choose to do type-specific deserialization based on volumeType.
	volumeType, err = untweakVolumeSourceFields(volumeType, obj)
	if err != nil {
		return nil, err
	}

	// Generic deserialization:
	volumeSourceContentsObj := util.ConvertMapKeysToCamelCase(volumeSourceLookup, obj)
	volumeType = util.HyphenToCamelCase(volumeSourceLookup, volumeType)
	volumeSourceObj := map[string]interface{}{
		volumeType: volumeSourceContentsObj,
	}

	b, err := json.Marshal(volumeSourceObj)
	if err != nil {
		return nil, util.InvalidValueForTypeErrorf(volumeSourceObj, v, "couldn't reserialize VolumeSource obj")
	}

	return b, nil
}

func (v *VolumeSource) UnmarshalJSON(data []byte) error {
	b, err := PreprocessVolumeSourceJSON(v, data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &v.VolumeSource)
	if err != nil {
		return util.InvalidValueForTypeErrorf(string(b), v, "couldn't deserialize")
	}

	return nil
}

func PostprocessVolumeSourceJSON(v interface{}, data []byte) ([]byte, error) {
	var err error
	obj := map[string]interface{}{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	if len(obj) != 1 {
		return nil, util.InvalidValueForTypeErrorf(obj, v, "should have one field: volume-source")
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
	snakeObj["type"] = kokiType

	return json.Marshal(snakeObj)
}

func (v VolumeSource) MarshalJSON() ([]byte, error) {
	var err error
	// EXTENSION POINT: Choose to do type-specific serialization based on the volume type.

	b, err := json.Marshal(v.VolumeSource)
	if err != nil {
		return nil, err
	}

	return PostprocessVolumeSourceJSON(v, b)
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
