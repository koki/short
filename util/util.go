package util

import (
	"fmt"
	"os"
	"reflect"

	"github.com/golang/glog"
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

func errorf(addedMsg, f interface{}, args ...interface{}) error {
	format := ""
	switch f.(type) {
	case string:
		format = f.(string)
	case fmt.Stringer:
		format = f.(fmt.Stringer).String()
	case error:
		format = f.(error).Error()
	default:
		glog.Errorf("unrecognized format type %v", f)
	}

	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\n %s", msg, addedMsg)
}
