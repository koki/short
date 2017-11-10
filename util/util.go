package util

import (
	"fmt"
	"os"
	"reflect"

	"github.com/golang/glog"
	"github.com/kr/pretty"
)

func ExitWithErr(msg interface{}) {
	glog.Error(msg)
	os.Exit(1)
}

func UsageErrorf(commandPath, f interface{}, args ...interface{}) error {
	return errorf(fmt.Sprintf("See '%s -h' for help and examples", commandPath), f, args...)
}

func TypeErrorf(t reflect.Type, f interface{}, args ...interface{}) error {
	return errorf(fmt.Sprintf("Unknown type '%s'", t), f, args...)
}

func TypeValueErrorf(obj, f interface{}, args ...interface{}) error {
	return errorf(fmt.Sprintf("Unknown value for type '%s'", reflect.TypeOf(obj)), f, args...)
}

func PrettyTypeError(obj interface{}, msg string) error {
	return TypeValueErrorf(obj, pretty.Sprintf("%s (%# v)", msg, obj))
}

func errorf(addedMsg, f interface{}, args ...interface{}) error {
	format := ""
	switch f := f.(type) {
	case string:
		format = f
	case fmt.Stringer:
		format = f.String()
	case error:
		format = f.Error()
	default:
		glog.Errorf("unrecognized format type %v", f)
	}

	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\n %s", msg, addedMsg)
}
