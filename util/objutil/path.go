package objutil

import (
	"strconv"

	"github.com/koki/short/util"
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
			return nil, util.ContextualizeErrorf(err, key)
		}
		if len(obj) <= int(i) {
			return nil, util.InvalidValueErrorf(i, "index not found")
		}

		val, err := AtPathIn(obj[i], remainder)
		if err != nil {
			return nil, util.ContextualizeErrorf(err, key)
		}
		return val, nil
	case map[string]interface{}:
		if val, ok := obj[key]; ok {
			val, err := AtPathIn(val, remainder)
			if err != nil {
				return nil, util.ContextualizeErrorf(err, key)
			}
			return val, nil
		}
		return nil, util.InvalidValueErrorf(key, "key not found")
	default:
		return nil, util.InvalidValueErrorf(key, "can only index into slice or map")
	}
}
