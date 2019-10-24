package tests

import (
	"github.com/ianamason/yices2_go_bindings/yices2"
	"testing"
)

func TestConfig0(t *testing.T) {
	yices2.Init()

	cfg := yices2.New_config()
	// Valid call
	yices2.Set_config(cfg, "mode", "push-pop")
	// Invalid name
	errcode := yices2.Set_config(cfg, "baz", "bar")
	error_string := yices2.Error_string()
	AssertEqual(t, errcode, -1)
	AssertEqual(t, error_string, "invalid parameter")
	// Invalid value
	errcode = yices2.Set_config(cfg, "mode", "bar")
	error_string = yices2.Error_string()
	AssertEqual(t, errcode, -1)
	AssertEqual(t, error_string, "value not valid for parameter")
	yices2.Default_config_for_logic(cfg, "QF_UFNIRA")
	yices2.Free_config(cfg)

	yices2.Exit()

}
