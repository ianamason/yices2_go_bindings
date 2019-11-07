package tests

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"os"
	"testing"
)

func TestErrors(t *testing.T) {
	yapi.Init()

	//First with no error
	errcode := yapi.Error_code()
	AssertEqual(t, errcode, yapi.NO_ERROR, "errcode == yices.NO_ERROR")
	yapi.Clear_error()
	errstr := yapi.Error_string()
	AssertEqual(t, errstr, "no error", "errstr == 'no error'")
	yapi.Print_error(os.Stderr)

	// Illegal - only scalar or uninterpreted types allowed
	bool_t := yapi.Bool_type()
	AssertTrue(t, yapi.Type_is_bool(bool_t), "yapi.Type_is_bool(bool_t)")
	const1 := yapi.Constant(bool_t, 0)
	error_string := yapi.Error_string()
	AssertEqual(t, const1, yapi.NULL_TERM, "const1 == yapi.NULL_TERM")
	AssertEqual(t, error_string, "invalid type in constant creation", "error_string == 'invalid type in constant creation'")

	yerror := yapi.YicesError()

	println(yerror.String())

	yapi.Clear_error()
	AssertEqual(t, yapi.Error_code(), yapi.NO_ERROR, "yapi.Error_code() == yapi.NO_ERROR")
	errstr = yapi.Error_string()
	AssertEqual(t, errstr, "no error", "errstr == 'no error'")
	yapi.Print_error(os.Stderr)
	yapi.Clear_error()
	AssertEqual(t, yapi.Error_code(), yapi.NO_ERROR, "yapi.Error_code() == yapi.NO_ERROR")

	yapi.Exit()
}
