package objutil

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"

	"github.com/koki/short/util"
)

// GetOnlyMapEntry get the only entry from a map. Error if the map doesn't contain exactly one entry.
func GetOnlyMapEntry(obj map[string]interface{}) (string, interface{}, error) {
	if len(obj) != 1 {
		return "", nil, util.InvalidInstanceErrorf(obj, "expected only one entry")
	}

	for key, val := range obj {
		return key, val, nil
	}

	glog.Fatal("unreachable")
	return "", nil, fmt.Errorf("unreachable")
}

func GetStringEntry(obj map[string]interface{}, key string) (string, error) {
	if val, ok := obj[key]; ok {
		if str, ok := val.(string); ok {
			return str, nil
		}

		return "", util.InvalidValueErrorf(obj, "entry for key (%s) is not a string", key)
	}

	return "", util.InvalidValueErrorf(obj, "no entry for key (%s)", key)
}

func UnmarshalMap(data map[string]interface{}, target interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		// This should be impossible.
		glog.Fatal(err)
	}

	// Let the caller of UnmarshalMap fill in the correct error context.
	return json.Unmarshal(b, target)
}

func MarshalMap(source interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(source)
	if err != nil {
		// Let the caller of MarshalMap fill in the correct error context.
		return nil, err
	}

	obj := map[string]interface{}{}
	err = json.Unmarshal(b, &obj)

	// Let the caller of MarshalMap fill in the correct error context.
	return obj, err
}
