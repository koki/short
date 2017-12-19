package jsonutil

import (
	"fmt"
	"strings"

	"github.com/koki/json"
	errutil "github.com/koki/structurederrors"
)

// GetOnlyMapEntry get the only entry from a map. Error if the map doesn't contain exactly one entry.
func GetOnlyMapEntry(obj map[string]interface{}) (string, interface{}, error) {
	if len(obj) != 1 {
		return "", nil, errutil.InvalidInstanceErrorf(obj, "expected only one entry")
	}

	for key, val := range obj {
		return key, val, nil
	}

	panic("unreachable")
	return "", nil, fmt.Errorf("unreachable")
}

func GetStringEntry(obj map[string]interface{}, key string) (string, error) {
	if val, ok := obj[key]; ok {
		if str, ok := val.(string); ok {
			return str, nil
		}

		return "", errutil.InvalidValueErrorf(obj, "entry for key (%s) is not a string", key)
	}

	return "", errutil.InvalidValueErrorf(obj, "no entry for key (%s)", key)
}

func UnmarshalMap(data map[string]interface{}, target interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		// This should be impossible.
		panic(err)
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

type ExtraneousFieldsError struct {
	Paths [][]string
}

func (e *ExtraneousFieldsError) Error() string {
	paths := make([]string, len(e.Paths))
	for i, path := range e.Paths {
		paths[i] = "$." + strings.Join(path, ".")
	}
	return fmt.Sprintf("extraneous fields (typos?) at paths: %s", strings.Join(paths, ", "))
}

// ExtraneousFieldPaths finds paths in "data" that aren't in "parsed" by serializing
//   "parsed" to a dictionary and recursively comparing its fields to those of "data".
func ExtraneousFieldPaths(data map[string]interface{}, parsed interface{}) ([][]string, error) {
	unparsed, err := MarshalMap(parsed)
	if err != nil {
		return nil, err
	}

	return ExtraneousMapPaths([]string{}, data, unparsed), nil
}

func ExtraneousMapPaths(prefix []string, before map[string]interface{}, after interface{}) [][]string {
	paths := [][]string{}
	if after, ok := after.(map[string]interface{}); ok {
		for key, beforeVal := range before {
			// Don't count empty fields as extraneous.
			if FieldValIsEmpty(beforeVal) {
				continue
			}
			newPrefix := ExtendPrefix(prefix, key)

			// Is the field also present in "after"?
			if afterVal, ok := after[key]; !ok {
				// Nope!
				paths = append(paths, newPrefix)
			} else {
				paths = append(paths, ExtraneousAnyPaths(newPrefix, beforeVal, afterVal)...)
			}
		}
	}

	return paths
}

func ExtraneousSlicePaths(prefix []string, before []interface{}, after interface{}) [][]string {
	paths := [][]string{}
	if after, ok := after.([]interface{}); ok {
		for i, beforeVal := range before {
			// Don't count empty fields as extraneous.
			if FieldValIsEmpty(beforeVal) {
				continue
			}
			newPrefix := ExtendPrefix(prefix, fmt.Sprintf("%d", i))

			// Is the field also present in "after"
			if i >= len(after) {
				// Nope!
				paths = append(paths, newPrefix)
			} else {
				afterVal := after[i]
				paths = append(paths, ExtraneousAnyPaths(newPrefix, beforeVal, afterVal)...)
			}
		}
	}

	return paths
}

func ExtendPrefix(prefix []string, segment string) []string {
	return append(append([]string{}, prefix...), segment)
}

func ExtraneousAnyPaths(prefix []string, before interface{}, after interface{}) [][]string {
	paths := [][]string{}
	if !FieldValIsEmpty(before) {
		switch before := before.(type) {
		case []interface{}:
			paths = append(paths, ExtraneousSlicePaths(prefix, before, after)...)
		case map[string]interface{}:
			paths = append(paths, ExtraneousMapPaths(prefix, before, after)...)
		}
	}

	return paths
}

func FieldValIsEmpty(fieldVal interface{}) bool {
	if fieldVal == nil {
		return true
	}

	switch fieldVal := fieldVal.(type) {
	case string:
		if len(fieldVal) == 0 {
			return true
		}
	case []interface{}:
		if len(fieldVal) == 0 {
			return true
		}
	case float64:
		if fieldVal == 0 {
			return true
		}
	case bool:
		if !fieldVal {
			return true
		}
	case map[string]interface{}:
		if len(fieldVal) == 0 {
			return true
		}
	}
	return false
}
