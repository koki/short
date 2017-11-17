package util

import (
	"strings"
	"unicode"
)

// ToSnakeCase convert a camelCased string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
// From: https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func ToSnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := range runes {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

// CamelToHyphenCase convert camelCase to hyphen-case
func CamelToHyphenCase(in string) string {
	return strings.Replace(ToSnakeCase(in), "_", "-", -1)
}

func ConvertMapKeysToSnakeCase(obj map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for key, val := range obj {
		result[ToSnakeCase(key)] = ConvertKeysToSnakeCase(val)
	}
	return result
}

func ConvertSliceKeystoSnakeCase(obj []interface{}) []interface{} {
	result := make([]interface{}, len(obj))
	for i, val := range obj {
		result[i] = ConvertKeysToSnakeCase(val)
	}
	return result
}

func ConvertKeysToSnakeCase(obj interface{}) interface{} {
	switch obj := obj.(type) {
	case map[string]interface{}:
		return ConvertMapKeysToSnakeCase(obj)
	case []interface{}:
		return ConvertSliceKeystoSnakeCase(obj)
	default:
		return obj
	}
}

// ToCamelCase convert string to camelCase.
// 'lookup' is for strings that can't be reconstructed algorithmically.
// (e.g. string contains acronym that got lowercased in snake_case)
func ToCamelCase(lookup map[string]string, in string) string {
	if s, ok := lookup[in]; ok {
		return s
	}

	runes := []rune(in)

	var out []rune
	capitalizeNext := false
	for _, r := range runes {
		if r == '_' {
			capitalizeNext = true
		} else if capitalizeNext {
			out = append(out, unicode.ToUpper(r))
			capitalizeNext = false
		} else {
			out = append(out, r)
		}
	}

	return string(out)
}

func HyphenToCamelCase(lookup map[string]string, in string) string {
	return ToCamelCase(lookup, strings.Replace(ToSnakeCase(in), "-", "_", -1))
}

func ConvertMapKeysToCamelCase(lookup map[string]string, obj map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for key, val := range obj {
		result[ToCamelCase(lookup, key)] = ConvertKeysToCamelCase(lookup, val)
	}
	return result
}

func ConvertSliceKeysToCamelCase(lookup map[string]string, obj []interface{}) []interface{} {
	result := make([]interface{}, len(obj))
	for i, val := range obj {
		result[i] = ConvertKeysToCamelCase(lookup, val)
	}
	return result
}

func ConvertKeysToCamelCase(lookup map[string]string, obj interface{}) interface{} {
	switch obj := obj.(type) {
	case map[string]interface{}:
		return ConvertMapKeysToCamelCase(lookup, obj)
	case []interface{}:
		return ConvertSliceKeysToCamelCase(lookup, obj)
	default:
		return obj
	}
}
