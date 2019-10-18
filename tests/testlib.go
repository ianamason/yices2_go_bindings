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
		t.Errorf("AssertEqual: %v : %v = %v : %v\n", lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
	}
}

func AssertTrue(t *testing.T, cond interface{}) {
	if cond != true {
		t.Errorf("AssertTrue: %v : %v\n", cond, reflect.TypeOf(cond))
	}
}

func AssertFalse(t *testing.T, cond interface{}) {
	if cond != false {
		t.Errorf("AssertFalse: %v : %v\n", cond, reflect.TypeOf(cond))
	}
}
