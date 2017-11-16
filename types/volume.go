package types

import (
	"encoding/json"

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
			volumeType = util.HyphenToCamelCase(volumeSourceLookup, vType)
			delete(obj, "type")
		default:
			return util.InvalidValueForTypeErrorf(string(data), v, "expected 'type' field to be a string")
		}
	} else {
		return util.InvalidValueForTypeErrorf(string(data), v, "no 'type' field")
	}

	// Choose to do type-specific deserialization based on volumeType

	// Generic deserialization:
	volumeSourceObj := util.ConvertKeysToCamelCase(volumeSourceLookup, obj)
	volumeObj := map[string]interface{}{
		"name":     volumeName,
		volumeType: volumeSourceObj,
	}
	b, err := json.Marshal(volumeSourceObj)
	if err != nil {
		return util.InvalidValueForTypeErrorf(volumeObj, v, "couldn't reserialize Volume obj")
	}
	err = json.Unmarshal(b, &v.Volume)
	if err != nil {
		return util.InvalidValueForTypeErrorf(string(b), v, "couldn't deserialize")
	}

	return nil
}
