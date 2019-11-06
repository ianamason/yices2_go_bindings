package tests

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"testing"
)

func TestConfig0(t *testing.T) {
	yapi.Init()

	var cfg yapi.Config_t

	yapi.Init_config(&cfg)


	// Valid call
	yapi.Set_config(cfg, "mode", "push-pop")
	// Invalid name
	errcode := yapi.Set_config(cfg, "baz", "bar")
	error_string := yapi.Error_string()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, error_string, "invalid parameter", "error_string == 'invalid parameter'")
	// Invalid value
	errcode = yapi.Set_config(cfg, "mode", "bar")
	error_string = yapi.Error_string()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, error_string, "value not valid for parameter", "error_string == 'value not valid for parameter'")
	yapi.Default_config_for_logic(cfg, "QF_UFNIRA")

	yapi.Close_config(&cfg)

	yapi.Exit()
}
