package tests

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"testing"
)

func TestConfig0(t *testing.T) {
	yapi.Init()

	var cfg yapi.ConfigT

	yapi.InitConfig(&cfg)

	// Valid call
	yapi.SetConfig(cfg, "mode", "push-pop")
	// Invalid name
	errcode := yapi.SetConfig(cfg, "baz", "bar")
	errorString := yapi.ErrorString()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, errorString, "invalid parameter", "errorString == 'invalid parameter'")
	// Invalid value
	errcode = yapi.SetConfig(cfg, "mode", "bar")
	errorString = yapi.ErrorString()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, errorString, "value not valid for parameter", "errorString == 'value not valid for parameter'")
	yapi.DefaultConfigForLogic(cfg, "QF_UFNIRA")

	yapi.CloseConfig(&cfg)

	yapi.Exit()
}
