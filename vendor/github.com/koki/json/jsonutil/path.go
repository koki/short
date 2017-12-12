package jsonutil

import (
	"strconv"

	serrors "github.com/koki/structurederrors"
)

func AtPathIn(obj interface{}, path []string) (interface{}, error) {
	if len(path) == 0 {
		return obj, nil
	}

	key := path[0]
	remainder := path[1:]
	switch obj := obj.(type) {
	case []interface{}:
		i, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, key)
		}
		if len(obj) <= int(i) {
			return nil, serrors.InvalidValueErrorf(i, "index not found")
		}

		val, err := AtPathIn(obj[i], remainder)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, key)
		}
		return val, nil
	case map[string]interface{}:
		if val, ok := obj[key]; ok {
			val, err := AtPathIn(val, remainder)
			if err != nil {
				return nil, serrors.ContextualizeErrorf(err, key)
			}
			return val, nil
		}
		return nil, serrors.InvalidValueErrorf(key, "key not found")
	default:
		return nil, serrors.InvalidValueErrorf(key, "can only index into slice or map")
	}
}
