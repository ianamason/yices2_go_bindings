package tests

import (
	"testing"
	"reflect"
)


func AssertNotEqual(t *testing.T, lhs interface{}, rhs interface{}) {
	if reflect.DeepEqual(lhs, rhs) {
		t.Errorf("AssertNotEqual: %v = %v\n", lhs, rhs)
	}
}

func AssertEqual(t *testing.T, lhs interface{}, rhs interface{}) {
	if !reflect.DeepEqual(lhs, rhs) {
		t.Errorf("AssertNotEqual: %v = %v\n", lhs, rhs)
	}
}
