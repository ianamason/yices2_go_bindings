package tests

import (
	"math"
	"reflect"
	"testing"
)

func isNil(x interface{}) bool {
	return x == nil || (reflect.ValueOf(x).Kind() == reflect.Ptr && reflect.ValueOf(x).IsNil())
}

// Signed returns True if it's argument is a signed integer, False othewise
func Signed(thing interface{}) (retval bool) {
	if thing == nil {
		return
	}
	tau := reflect.TypeOf(thing).String()
	switch tau {
	case "int", "int8", "int16", "int32", "int64":
		retval = true
	}
	return
}

// Unsigned returns True if it's argument is an unsigned integer, False othewise
func Unsigned(thing interface{}) (retval bool) {
	if thing == nil {
		return
	}
	switch reflect.TypeOf(thing).String() {
	case "uint", "uint8", "uint16", "uint32", "uint64":
		retval = true
	}
	return
}

func convertSigned2Int64(x interface{}, xint *int64) bool {
	xtyp := reflect.TypeOf(x)
	switch xtyp.Kind() {
	case reflect.Int:
		*xint = int64(x.(int))
		return true
	case reflect.Int8:
		*xint = int64(x.(int8))
		return true
	case reflect.Int16:
		*xint = int64(x.(int16))
		return true
	case reflect.Int32:
		*xint = int64(x.(int32))
		return true
	case reflect.Int64:
		*xint = x.(int64)
		return true
	default:
		return false
	}
}

func convertUnsigned2Int64(y interface{}, yint *int64) bool {
	ytyp := reflect.TypeOf(y)
	switch ytyp.Kind() {
	case reflect.Uint:
		*yint = int64(y.(uint))
		return true
	case reflect.Uint8:
		*yint = int64(y.(uint8))
		return true
	case reflect.Uint16:
		*yint = int64(y.(uint16))
		return true
	case reflect.Uint32:
		*yint = int64(y.(uint32))
		return true
	case reflect.Uint64:
		yval := y.(uint64)
		if yval > math.MaxInt64 {
			return false
		}
		*yint = int64(y.(uint64))
		return true
	default:
		return false
	}
}

func convertUnsigned2UInt(y interface{}, yuint *uint) bool {
	ytyp := reflect.TypeOf(y)
	switch ytyp.Kind() {
	case reflect.Uint:
		*yuint = uint(y.(uint))
		return true
	case reflect.Uint8:
		*yuint = uint(y.(uint8))
		return true
	case reflect.Uint16:
		*yuint = uint(y.(uint16))
		return true
	case reflect.Uint32:
		*yuint = uint(y.(uint32))
		return true
	case reflect.Uint64:
		*yuint = uint(y.(uint64))
		return true
	default:
		return false
	}
}

// SignedUnsignedEqual compares a signed integer with an unsigned one
func SignedUnsignedEqual(x interface{}, y interface{}) bool {
	if x == nil || y == nil {
		return false
	}
	var xint int64
	var yint int64
	if !convertSigned2Int64(x, &xint) {
		return false
	}
	if !convertUnsigned2Int64(y, &yint) {
		return false
	}
	if xint != yint {
		return false
	}
	return true
}

// SignedEqual encodes equality on signed integers.
func SignedEqual(x interface{}, y interface{}) bool {
	if x == nil || y == nil {
		return false
	}
	var xint int64
	var yint int64
	if !convertSigned2Int64(x, &xint) {
		return false
	}
	if !convertSigned2Int64(y, &yint) {
		return false
	}
	if xint != yint {
		return false
	}
	return true
}

// UnsignedEqual encodes equality on unsigned integers.
func UnsignedEqual(x interface{}, y interface{}) bool {
	if x == nil || y == nil {
		return false
	}
	var xuint uint
	var yuint uint
	if !convertUnsigned2UInt(x, &xuint) {
		return false
	}
	if !convertUnsigned2UInt(y, &yuint) {
		return false
	}
	if xuint != yuint {
		return false
	}
	return true
}

// AssertNotEqual asserts that the lhs is not equal to the rhs.
func AssertNotEqual(t *testing.T, lhs interface{}, rhs interface{}, where ...string) {
	if Signed(lhs) && Signed(rhs) {
		if SignedEqual(lhs, rhs) {
			t.Errorf("%s : AssertNotEqual of signed integers %v : %v = %v : %v\n", where, lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
		} else {
			return
		}
	}
	if Signed(lhs) && Unsigned(rhs) {
		if SignedUnsignedEqual(lhs, rhs) {
			t.Errorf("%s : AssertNotEqual of signed/unsigned integers %v : %v = %v : %v\n", where, lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
		} else {
			return
		}
	}
	if Signed(rhs) && Unsigned(lhs) {
		if SignedUnsignedEqual(rhs, lhs) {
			t.Errorf("%s : AssertNotEqual of signed/unsigned integers %v : %v = %v : %v\n", where, lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
		} else {
			return
		}
	}

	if reflect.DeepEqual(lhs, rhs) {
		t.Errorf("AssertNotEqual: %v = %v\n", lhs, rhs)
	}
}

// AssertEqual asserts that the lhs is equal to the rhs.
func AssertEqual(t *testing.T, lhs interface{}, rhs interface{}, where ...string) {
	if isNil(lhs) && isNil(rhs) {
		return
	}
	if isNil(lhs) || isNil(rhs) {
		t.Errorf("%s : AssertEqual of pointer %v : %v != %v : %v\n", where, lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
	}
	if Signed(lhs) && Signed(rhs) {
		if !SignedEqual(lhs, rhs) {
			t.Errorf("%s : AssertEqual of signed integers %v : %v != %v : %v\n", where, lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
		} else {
			return
		}
	}
	if Signed(lhs) && Unsigned(rhs) {
		if !SignedUnsignedEqual(lhs, rhs) {
			t.Errorf("%s : AssertEqual of signed/unsigned integers %v : %v != %v : %v\n", where, lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
		} else {
			return
		}
	}
	if Signed(rhs) && Unsigned(lhs) {
		if !SignedUnsignedEqual(rhs, lhs) {
			t.Errorf("%s : AssertEqual of signed/unsigned integers %v : %v = %v : %v\n", where, lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
		} else {
			return
		}
	}
	if !reflect.DeepEqual(lhs, rhs) {
		t.Errorf("%s : AssertEqual %v : %v = %v : %v\n", where, lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
	}
}

// AssertTrue asserts that cond is True
func AssertTrue(t *testing.T, cond interface{}, where ...string) {
	if cond != true {
		t.Errorf("%s AssertTrue %v : %v\n", where, cond, reflect.TypeOf(cond))
	}
}

// AssertFalse asserts that cond is True
func AssertFalse(t *testing.T, cond interface{}, where ...string) {
	if cond != false {
		t.Errorf("%s AssertFalse %v : %v\n", where, cond, reflect.TypeOf(cond))
	}
}
