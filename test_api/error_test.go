package tests

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"os"
	"testing"
)

func TestErrors(t *testing.T) {
	yapi.Init()

	//First with no error
	errcode := yapi.ErrorCode()
	AssertEqual(t, errcode, yapi.NoError, "errcode == yices.NoError")
	yapi.ClearError()
	errstr := yapi.ErrorString()
	AssertEqual(t, errstr, "no error", "errstr == 'no error'")
	yapi.PrintError(os.Stderr)

	// Illegal - only scalar or uninterpreted types allowed
	boolT := yapi.BoolType()
	AssertTrue(t, yapi.TypeIsBool(boolT), "yapi.TypeIsBool(boolT)")
	const1 := yapi.Constant(boolT, 0)
	errorString := yapi.ErrorString()
	AssertEqual(t, const1, yapi.NullTerm, "const1 == yapi.NullTerm")
	AssertEqual(t, errorString, "invalid type in constant creation", "errorString == 'invalid type in constant creation'")

	yerror := yapi.YicesError()

	println(yerror.String())

	yapi.ClearError()
	AssertEqual(t, yapi.ErrorCode(), yapi.NoError, "yapi.ErrorCode() == yapi.NoError")
	errstr = yapi.ErrorString()
	AssertEqual(t, errstr, "no error", "errstr == 'no error'")
	yapi.PrintError(os.Stderr)
	yapi.ClearError()
	AssertEqual(t, yapi.ErrorCode(), yapi.NoError, "yapi.ErrorCode() == yapi.NoError")

	yapi.Exit()
}
