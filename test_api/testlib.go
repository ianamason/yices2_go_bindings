package tests

import (
	"math"
	"reflect"
	"testing"
)

func isNil(x interface{}) bool {
	return x == nil || (reflect.ValueOf(x).Kind() == reflect.Ptr && reflect.ValueOf(x).IsNil())
}

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

// x is signed; y is unsigned
func SignedUnsignedEqual(x interface{}, y interface{}) bool {

	if x == nil || y == nil {
		return false
	}

	var xint int64
	var yint int64

	xtyp := reflect.TypeOf(x)
	switch xtyp.Kind() {
	case reflect.Int:
		xint = int64(x.(int))
	case reflect.Int8:
		xint = int64(x.(int8))
	case reflect.Int16:
		xint = int64(x.(int16))
	case reflect.Int32:
		xint = int64(x.(int32))
	case reflect.Int64:
		xint = x.(int64)
	default:
		return false
	}

	if xint < 0 {
		return false
	}

	ytyp := reflect.TypeOf(y)
	switch ytyp.Kind() {
	case reflect.Uint:
		yint = int64(y.(uint))
	case reflect.Uint8:
		yint = int64(y.(uint8))
	case reflect.Uint16:
		yint = int64(y.(uint16))
	case reflect.Uint32:
		yint = int64(y.(uint32))
	case reflect.Uint64:
		yval := y.(uint64)
		if yval > math.MaxInt64 {
			return false
		}
		yint = int64(y.(uint64))
	default:
		return false
	}

	if xint != yint {
		return false
	}

	return true

}

func SignedEqual(x interface{}, y interface{}) bool {

	if x == nil || y == nil {
		return false
	}

	var xint int64
	var yint int64

	xtyp := reflect.TypeOf(x)
	switch xtyp.Kind() {
	case reflect.Int:
		xint = int64(x.(int))
	case reflect.Int8:
		xint = int64(x.(int8))
	case reflect.Int16:
		xint = int64(x.(int16))
	case reflect.Int32:
		xint = int64(x.(int32))
	case reflect.Int64:
		xint = x.(int64)
	default:
		return false
	}

	ytyp := reflect.TypeOf(y)
	switch ytyp.Kind() {
	case reflect.Int:
		yint = int64(y.(int))
	case reflect.Int8:
		yint = int64(y.(int8))
	case reflect.Int16:
		yint = int64(y.(int16))
	case reflect.Int32:
		yint = int64(y.(int32))
	case reflect.Int64:
		yint = y.(int64)
	default:
		return false
	}

	if xint != yint {
		return false
	}

	return true
}

func UnsignedEqual(x interface{}, y interface{}) bool {

	if x == nil || y == nil {
		return false
	}

	var xuint uint
	var yuint uint

	xtyp := reflect.TypeOf(x)
	switch xtyp.Kind() {
	case reflect.Uint:
		xuint = uint(x.(uint))
	case reflect.Uint8:
		xuint = uint(x.(uint8))
	case reflect.Uint16:
		xuint = uint(x.(uint16))
	case reflect.Uint32:
		xuint = uint(x.(uint32))
	case reflect.Uint64:
		xuint = uint(x.(uint64))
	default:
		return false
	}

	ytyp := reflect.TypeOf(y)
	switch ytyp.Kind() {
	case reflect.Uint:
		yuint = uint(y.(uint))
	case reflect.Uint8:
		yuint = uint(y.(uint8))
	case reflect.Uint16:
		yuint = uint(y.(uint16))
	case reflect.Uint32:
		yuint = uint(y.(uint32))
	case reflect.Uint64:
		yuint = uint(y.(uint64))
	default:
		return false
	}

	if xuint != yuint {
		return false
	}

	return true
}

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

func AssertTrue(t *testing.T, cond interface{}, where ...string) {
	if cond != true {
		t.Errorf("%s AssertTrue %v : %v\n", where, cond, reflect.TypeOf(cond))
	}
}

func AssertFalse(t *testing.T, cond interface{}, where ...string) {
	if cond != false {
		t.Errorf("%s AssertFalse %v : %v\n", where, cond, reflect.TypeOf(cond))
	}
}
