package util

import (
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Int32Ptr(i int32) *int32 {
	return &i
}

func IntOrStringPtr(i intstr.IntOrString) *intstr.IntOrString {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}

func BoolPtrOrNil(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

func StringPtrOrNil(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return &s
}

func FromBoolPtr(b *bool) bool {
	if b == nil {
		return false
	}

	return *b
}

func FromStringPtr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
