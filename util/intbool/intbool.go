package intbool

import (
	"encoding/json"
	"math"
	"runtime/debug"

	"github.com/golang/glog"

	"github.com/koki/short/util"
)

type IntOrBool struct {
	Type    Type
	IntVal  int32
	BoolVal bool
}

// Type represents the stored type of IntOrBool.
type Type int

const (
	Int  Type = iota // The IntOrBool holds an int.
	Bool             // The IntOrBool holds a bool.
)

func FromInt(val int) *IntOrBool {
	if val > math.MaxInt32 || val < math.MinInt32 {
		glog.Errorf("value: %d overflows int32\n%s\n", val, debug.Stack())
	}
	return &IntOrBool{Type: Int, IntVal: int32(val)}
}

func FromBool(val bool) *IntOrBool {
	return &IntOrBool{Type: Bool, BoolVal: val}
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (ib *IntOrBool) UnmarshalJSON(value []byte) error {
	var x float64
	err := json.Unmarshal(value, &x)
	if err == nil {
		ib.Type = Int
		ib.IntVal = int32(x)
		return nil
	}

	var b bool
	err = json.Unmarshal(value, &b)
	if err == nil {
		ib.Type = Bool
		ib.BoolVal = b
		return nil
	}

	return util.InvalidValueForTypeErrorf(string(value), ib, "couldn't deserialize")
}

// MarshalJSON implements the json.Marshaller interface.
func (ib IntOrBool) MarshalJSON() ([]byte, error) {
	switch ib.Type {
	case Int:
		return json.Marshal(ib.IntVal)
	case Bool:
		return json.Marshal(ib.BoolVal)
	default:
		return []byte{}, util.InvalidInstanceError(ib.Type)
	}
}
