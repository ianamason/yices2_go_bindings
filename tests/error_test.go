package tests

import (
	"os"
	"github.com/ianamason/yices2_go_bindings/yices2"
	"testing"
)

func TestErrors(t *testing.T) {
	yices2.Init()


	//First with no error
	errcode := yices2.Error_code()
	AssertEqual(t, errcode, 0, "errcode == 0")
	yices2.Clear_error()
	errstr := yices2.Error_string()
	AssertEqual(t, errstr, "no error", "errstr == 'no error'")
	yices2.Print_error(os.Stderr)

	// Illegal - only scalar or uninterpreted types allowed
	bool_t := yices2.Bool_type()
	AssertTrue(t, yices2.Type_is_bool(bool_t), "yices2.Type_is_bool(bool_t)")
	const1 := yices2.Constant(bool_t, 0)
	error_string := yices2.Error_string()
	AssertEqual(t, const1, yices2.NULL_TERM, "const1 == yices2.NULL_TERM")
	AssertEqual(t, error_string, "invalid type in constant creation", "error_string == 'invalid type in constant creation'")
	yices2.Clear_error()
	AssertEqual(t, yices2.Error_code(), 0, "yices2.Error_code() == 0")
	errstr = yices2.Error_string()
	AssertEqual(t, errstr, "no error", "errstr == 'no error'")
	yices2.Print_error(os.Stderr)
	yices2.Clear_error()
	AssertEqual(t, yices2.Error_code(), 0, "yices2.Error_code() == 0")

	yices2.Exit()
}
