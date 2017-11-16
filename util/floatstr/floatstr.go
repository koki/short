package floatstr

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/koki/short/util"
)

type FloatOrString struct {
	Type      Type
	FloatVal  float64
	StringVal string
}

// Type represents the stored type of IntOrBool.
type Type int

const (
	Float  Type = iota // The FloatOrString holds an Float.
	String             // The FloatOrString holds a string.
)

func FromFloat(val float64) *FloatOrString {
	return &FloatOrString{Type: Float, FloatVal: val}
}

func FromString(val string) *FloatOrString {
	return &FloatOrString{Type: String, StringVal: val}
}

func Parse(val string) *FloatOrString {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return FromString(val)
	}
	return FromFloat(f)
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (fs *FloatOrString) UnmarshalJSON(value []byte) error {
	var x float64
	err := json.Unmarshal(value, &x)
	if err == nil {
		fs.Type = Float
		fs.FloatVal = x
		return nil
	}

	var s string
	err = json.Unmarshal(value, &s)
	if err == nil {
		fs.Type = String
		fs.StringVal = s
		return nil
	}

	return util.InvalidValueForTypeErrorf(string(value), fs, "couldn't deserialize")
}

// MarshalJSON implements the json.Marshaller interface.
func (fs FloatOrString) MarshalJSON() ([]byte, error) {
	switch fs.Type {
	case Float:
		return json.Marshal(fs.FloatVal)
	case String:
		return json.Marshal(fs.StringVal)
	default:
		return []byte{}, util.InvalidInstanceError(fs.Type)
	}
}

// String returns the string value, or the float value formatted as a string.
func (fs *FloatOrString) String() string {
	if fs.Type == String {
		return fs.StringVal
	}
	return fmt.Sprintf("%v", fs.FloatVal)
}
