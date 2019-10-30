package tests

import (
	"testing"
	"reflect"
)

func Signed(thing interface{}) (retval bool){
	if thing == nil { return }

	tau := reflect.TypeOf(thing).String()

	switch tau {
		case "int", "int8", "int16", "int32", "int64":
		retval = true
	}

	return
}

func Unsigned(thing interface{}) (retval bool){
	if thing == nil { return }

	switch reflect.TypeOf(thing).String() {
		case "uint", "uint8", "uint16", "uint32", "uint64":
		retval = true
	}
	return
}


func SignedEqual(x interface{}, y interface{}) bool {

	if x == nil || y == nil {
		return false
	}

	var xint int = 0
	var yint int = 0

	xtyp := reflect.TypeOf(x)
	switch xtyp.Kind() {
	case reflect.Int:
		xint = int(x.(int))
	case reflect.Int8:
		xint = int(x.(int8))
	case reflect.Int16:
		xint = int(x.(int16))
	case reflect.Int32:
		xint = int(x.(int32))
	case reflect.Int64:
		xint = int(x.(int64))
	default:
		return false
	}


	ytyp := reflect.TypeOf(y)
	switch ytyp.Kind() {
	case reflect.Int:
		yint = int(y.(int))
	case reflect.Int8:
		yint = int(y.(int8))
	case reflect.Int16:
		yint = int(y.(int16))
	case reflect.Int32:
		yint = int(y.(int32))
	case reflect.Int64:
		yint = int(y.(int64))
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

	var xuint uint = 0
	var yuint uint = 0

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


func AssertNotEqual(t *testing.T, lhs interface{}, rhs interface{}) {
	if reflect.DeepEqual(lhs, rhs) {
		t.Errorf("AssertNotEqual: %v = %v\n", lhs, rhs)
	}
}

func AssertEqual(t *testing.T, lhs interface{}, rhs interface{}) {
	if Signed(lhs) && Signed(rhs) {
		if !SignedEqual(lhs, rhs) {
			//panic("here")
			t.Errorf("AssertEqual of signed integers: %v : %v = %v : %v\n", lhs, reflect.TypeOf(lhs), rhs, reflect.TypeOf(rhs))
		} else {
			return
		}
	}
	if !reflect.DeepEqual(lhs, rhs) {
		//panic("here")
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
