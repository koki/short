package util

import (
	"fmt"
	"os"
	"reflect"

	"github.com/golang/glog"
	"github.com/kr/pretty"
)

var verboseErrors = false

func SetVerboseErrors(verbose bool) {
	verboseErrors = verbose
}

func ExitWithErr(msg interface{}) {
	glog.Error(msg)
	os.Exit(1)
}

func UsageErrorf(commandPath, f interface{}, args ...interface{}) error {
	return errorf(fmt.Sprintf("See '%s -h' for help and examples", commandPath), f, args...)
}

// TypeError means that obj has an unexpected type.
func TypeError(obj interface{}) error {
	return fmt.Errorf("unrecognized type (%s)", reflect.TypeOf(obj))
}

// TypeErrorf is like TypeError, except with a custom message.
func TypeErrorf(obj interface{}, msgFormat string, args ...interface{}) error {
	return errorf(pretty.Sprintf(msgFormat, args...), "unrecognized type (%s)", reflect.TypeOf(obj))
}

// InvalidInstanceError means that obj is the correct type, but there's something
// wrong with its contents.
func InvalidInstanceError(obj interface{}) error {
	return instanceError(obj, "unrecognized instance")
}

// InvalidInstanceErrorf is like InvalidInstanceError, except with a custom message.
func InvalidInstanceErrorf(obj interface{}, msgFormat string, args ...interface{}) error {
	return instanceError(obj, pretty.Sprintf(msgFormat, args...))
}

// InvalidValueErrorf is used when the type isn't meaningful--just the contents and the
// context matter.
func InvalidValueErrorf(val interface{}, msgFormat string, args ...interface{}) error {
	if verboseErrors {
		return pretty.Errorf("%s\n(%# v)", pretty.Sprintf(msgFormat, args...), val)
	}

	return pretty.Errorf("(%s) value: %s", reflect.TypeOf(val), pretty.Sprintf(msgFormat, args...))
}

// InvalidValueForTypeErrorf is used when the type isn't meaningful--just the contents and the
// context matter.
func InvalidValueForTypeErrorf(val, typedObj interface{}, msgFormat string, args ...interface{}) error {
	if verboseErrors {
		return pretty.Errorf("for type (%s), unrecognized value\n%s\n%# v", reflect.TypeOf(typedObj), pretty.Sprintf(msgFormat, args...), val)
	}

	return pretty.Errorf("for type (%s), unrecognized (%s) value: %s", reflect.TypeOf(typedObj), reflect.TypeOf(val), pretty.Sprintf(msgFormat, args...))
}

// ContextualizeErrorf is for adding an additional message to an existing error.
// This method is only intended for simple messages (contextFormat).
// e.g. If the context includes a printout of a Go struct, use one of the other error generators in this package.
func ContextualizeErrorf(err error, contextFormat string, contextArgs ...interface{}) error {
	return pretty.Errorf("%s: %s", pretty.Sprintf(contextFormat, contextArgs...), err.Error())
}

func instanceError(obj interface{}, msg string) error {
	if verboseErrors {
		return pretty.Errorf("%s: %s\n(%# v)", reflect.TypeOf(obj), msg, obj)
	}

	return pretty.Errorf("%s: %s", reflect.TypeOf(obj), msg)
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

	msg := pretty.Sprintf(format, args...)
	return fmt.Errorf("%s\n %s", msg, addedMsg)
}
