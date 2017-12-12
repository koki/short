package structurederrors

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/golang/glog"
	"github.com/kr/pretty"
	"github.com/kr/text"
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

// TypeContextErrorf is like TypeErrorf, but it adds its message as context to an existing error.
func TypeContextErrorf(baseError error, obj interface{}, msgFormat string, args ...interface{}) *ErrorWithContext {
	msg := pretty.Sprintf(msgFormat, args)
	typeMsg := pretty.Sprintf("unrecognized type (%s)", reflect.TypeOf(obj))
	return ContextualizeErrorf(baseError, "%s: %s", msg, typeMsg)
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

// InvalidInstanceContextErrorf is like InvalidInstanceContextErrorf, except it adds its message as context to an existing error.
func InvalidInstanceContextErrorf(baseError error, obj interface{}, msgFormat string, args ...interface{}) *ErrorWithContext {
	return ContextualizeErrorf(baseError, InvalidInstanceErrorf(obj, msgFormat, args...).Error())
}

// InvalidValueErrorf is used when the type isn't meaningful--just the contents and the
// context matter.
func InvalidValueErrorf(val interface{}, msgFormat string, args ...interface{}) error {
	if verboseErrors {
		return pretty.Errorf("%s\n(%# v)", pretty.Sprintf(msgFormat, args...), val)
	}

	return pretty.Errorf("(%s) value: %s", reflect.TypeOf(val), pretty.Sprintf(msgFormat, args...))
}

func InvalidValueContextErrorf(baseError error, val interface{}, msgFormat string, args ...interface{}) *ErrorWithContext {
	return ContextualizeErrorf(baseError, InvalidValueErrorf(val, msgFormat, args...).Error())
}

// InvalidValueForTypeError is used when the type isn't meaningful--just the contents and the
// context matter.
func InvalidValueForTypeError(val, typedObj interface{}) error {
	if verboseErrors {
		return pretty.Errorf("for type (%s), unrecognized value\n%# v", reflect.TypeOf(typedObj), val)
	}

	return pretty.Errorf("for type (%s), unrecognized (%s) value", reflect.TypeOf(typedObj), reflect.TypeOf(val))
}

func InvalidValueForTypeContextError(baseError error, val, typedObj interface{}) *ErrorWithContext {
	return ContextualizeErrorf(baseError, InvalidValueForTypeError(val, typedObj).Error())
}

// InvalidValueForTypeErrorf is used when the type isn't meaningful--just the contents and the
// context matter.
func InvalidValueForTypeErrorf(val, typedObj interface{}, msgFormat string, args ...interface{}) error {
	if verboseErrors {
		return pretty.Errorf("for type (%s), unrecognized value\n%s\n%# v", reflect.TypeOf(typedObj), pretty.Sprintf(msgFormat, args...), val)
	}

	return pretty.Errorf("for type (%s), unrecognized (%s) value: %s", reflect.TypeOf(typedObj), reflect.TypeOf(val), pretty.Sprintf(msgFormat, args...))
}

func InvalidValueForTypeContextErrorf(baseError error, val, typedObj interface{}, msgFormat string, args ...interface{}) *ErrorWithContext {
	return ContextualizeErrorf(baseError, InvalidValueForTypeErrorf(val, typedObj, msgFormat, args...).Error())
}

type ErrorWithContext struct {
	BaseError error
	Context   []string
}

func (e *ErrorWithContext) Error() string {
	context := ReversedStringsList(e.Context)

	return strings.Join(append(context, e.BaseError.Error()), ": ")
}

func (e *ErrorWithContext) PrettyError() string {
	context := make([]string, len(e.Context))
	indent := ""
	for i, contextItem := range ReversedStringsList(e.Context) {
		context[i] = text.Indent(contextItem, indent)
		indent = indent + "  "
	}

	contextString := strings.Join(context, "\n")
	errString := text.Indent(PrettyError(e.BaseError), indent)
	return fmt.Sprintf("%s\n%s", contextString, errString)
}

// ContextualizeErrorf is for adding an additional message to an existing error.
// This method is only intended for simple messages (contextFormat).
// e.g. If the context includes a printout of a Go struct, use one of the other error generators in this package.
func ContextualizeErrorf(err error, contextFormat string, contextArgs ...interface{}) *ErrorWithContext {
	contextMsg := pretty.Sprintf(contextFormat, contextArgs...)
	switch err := err.(type) {
	case *ErrorWithContext:
		err.Context = append(err.Context, pretty.Sprintf(contextFormat, contextArgs...))
		return err
	default:
		return &ErrorWithContext{
			BaseError: err,
			Context:   []string{contextMsg},
		}
	}
}

func PrettyError(err error) string {
	switch err := err.(type) {
	case *ErrorWithContext:
		return err.PrettyError()
	default:
		return err.Error()
	}
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

func ReversedStringsList(stringsList []string) []string {
	if stringsList == nil {
		return nil
	}

	l := len(stringsList)
	reversed := make([]string, l)
	for i, str := range stringsList {
		j := l - 1 - i
		reversed[j] = str
	}

	return reversed
}
