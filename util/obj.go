package util

import (
	"fmt"

	"github.com/golang/glog"
)

// GetOnlyMapEntry get the only entry from a map. Error if the map doesn't contain exactly one entry.
func GetOnlyMapEntry(obj map[string]interface{}) (string, interface{}, error) {
	if len(obj) != 1 {
		return "", nil, InvalidInstanceErrorf(obj, "expected only one entry")
	}

	for key, val := range obj {
		return key, val, nil
	}

	glog.Fatal("unreachable")
	return "", nil, fmt.Errorf("unreachable")
}
