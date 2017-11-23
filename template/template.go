package template

import (
	"regexp"

	"github.com/kr/pretty"

	"github.com/koki/short/util"
)

/*

Template "holes" are represented as the string "${NAME}".

If a "hole" is part of (but not all of) a string, then only string/number values are supported.
This behavior is defined in `generic.fillString`.

Parameter values may also have document structure and will retain
this structure when inserted into the template.

*/

// Resolver gets the value to substitute into the template.
type Resolver func(ident string) (interface{}, error)

func FillTemplate(template interface{}, resolver Resolver) (interface{}, error) {
	return ReplaceAny(template, resolver)
}

func ReplaceAny(template interface{}, resolver Resolver) (interface{}, error) {
	switch template := template.(type) {
	case string:
		return ReplaceString(template, resolver)
	case []interface{}:
		return ReplaceSlice(template, resolver)
	case map[string]interface{}:
		return ReplaceMap(template, resolver)
	default:
		// No template parameters in other data types.
	}

	return template, nil
}

func ReplaceMap(template map[string]interface{}, resolver Resolver) (map[string]interface{}, error) {
	var err error
	newTemplate := map[string]interface{}{}
	for key, val := range template {
		newTemplate[key], err = ReplaceAny(val, resolver)
		if err != nil {
			return nil, err
		}
	}

	return newTemplate, nil
}

func ReplaceSlice(template []interface{}, resolver Resolver) ([]interface{}, error) {
	var err error
	newTemplate := make([]interface{}, len(template))
	for ix, val := range template {
		newTemplate[ix], err = ReplaceAny(val, resolver)
		if err != nil {
			return nil, err
		}
	}

	return newTemplate, nil
}

func ReplaceString(template string, resolver Resolver) (interface{}, error) {
	// Find all template holes and Replace them with param values.
	expanded, modified, err := expandString(template, resolver)
	if err != nil {
		return nil, err
	}

	if modified {
		return expanded, nil
	}

	return fillString(template, resolver)
}

// Returns true if it expanded the template.
func expandString(template string, resolver Resolver) (interface{}, bool, error) {
	re := regexp.MustCompile("^\\$\\{([^\\{\\}]*)\\}$")
	matches := re.FindStringSubmatch(template)
	if len(matches) == 0 {
		return template, false, nil
	}

	key := matches[1]
	val, err := resolver(key)
	if err != nil {
		return nil, false, err
	}

	return val, true, nil
}

func fillString(template string, resolver Resolver) (string, error) {
	re := regexp.MustCompile("\\$\\{[^\\{\\}]*\\}")
	errors := []error{}
	result := re.ReplaceAllFunc([]byte(template), func(match []byte) []byte {
		key := match[2 : len(match)-1]
		val, err := resolver(string(key))
		if err != nil {
			errors = append(errors, err)
		}

		switch val := val.(type) {
		case string:
			return []byte(val)
		case float64:
			return []byte(pretty.Sprintf("%v", val))
		case int:
			return []byte(pretty.Sprintf("%v", val))
		default:
			errors = append(errors, util.InvalidValueErrorf(val, "not a string or number"))
			return match
		}
	})

	if len(errors) > 0 {
		return "", errors[0]
	}

	return string(result), nil
}
