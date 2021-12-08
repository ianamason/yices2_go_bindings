package tests

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"os"
	"testing"
)

// bad trivial use
func TestInterpolationContext0(t *testing.T) {
	yapi.Init()

	var cfg yapi.ConfigT

	yapi.InitConfig(&cfg)

	var ctxA yapi.ContextT
	var ctxB yapi.ContextT

	yapi.InitContext(cfg, &ctxA)
	yapi.InitContext(cfg, &ctxB)

	yapi.CloseConfig(&cfg)

	var params yapi.ParamT
	yapi.InitParamRecord(&params)
	yapi.DefaultParamsForContext(ctxA, params)

	smtStat, _, _ := yapi.CheckContextWithInterpolation(ctxA, ctxB, params, true)

	println(smtStat)

	errstr := yapi.ErrorString()
	AssertEqual(t, errstr, "operation not supported by the context", "errstr == 'operation not supported by the context'")

	yapi.PrintError(os.Stderr)

	yapi.CloseContext(&ctxA)
	yapi.CloseContext(&ctxB)

	yapi.Exit()
}

// good trival use
func TestInterpolationContext1(t *testing.T) {
	yapi.Init()

	var cfg yapi.ConfigT

	yapi.InitConfig(&cfg)

	yapi.SetConfig(cfg, "solver-type", "mcsat")
	yapi.SetConfig(cfg, "model-interpolation", "true")

	var ctxA yapi.ContextT
	var ctxB yapi.ContextT

	yapi.InitContext(cfg, &ctxA)
	yapi.InitContext(cfg, &ctxB)

	yapi.CloseConfig(&cfg)

	var params yapi.ParamT
	yapi.InitParamRecord(&params)
	yapi.DefaultParamsForContext(ctxA, params)

	smtStat, _, _ := yapi.CheckContextWithInterpolation(ctxA, ctxB, params, true)

	println(smtStat)

	errstr := yapi.ErrorString()
	AssertEqual(t, errstr, "no error", "errstr == 'no error'")

	yapi.PrintError(os.Stderr)

	yapi.CloseContext(&ctxA)
	yapi.CloseContext(&ctxB)

	yapi.Exit()
}

// sat
func TestInterpolationContext2(t *testing.T) {
	yapi.Init()

	var cfg yapi.ConfigT

	yapi.InitConfig(&cfg)

	yapi.SetConfig(cfg, "solver-type", "mcsat")
	yapi.SetConfig(cfg, "model-interpolation", "true")

	var ctxA yapi.ContextT
	var ctxB yapi.ContextT

	yapi.InitContext(cfg, &ctxA)
	yapi.InitContext(cfg, &ctxB)

	yapi.CloseConfig(&cfg)

	var params yapi.ParamT
	yapi.InitParamRecord(&params)
	yapi.DefaultParamsForContext(ctxA, params)

	realT := yapi.RealType()
	r1 := yapi.NewUninterpretedTerm(realT)
	yapi.SetTermName(r1, "r1")
	r2 := yapi.NewUninterpretedTerm(realT)
	yapi.SetTermName(r2, "r2")

	fmla1 := yapi.ParseTerm("(> r1 3)")
	fmla2 := yapi.ParseTerm("(< r1 4)")
	fmla3 := yapi.ParseTerm("(< (- r1 r2) 0)")

	yapi.AssertFormulas(ctxA, []yapi.TermT{fmla1, fmla2, fmla3})

	smtStat, modelp, interpolantp := yapi.CheckContextWithInterpolation(ctxA, ctxB, params, true)

	AssertEqual(t, smtStat, yapi.StatusSat, "smtStat == yapi.StatusSat")

	AssertNotEqual(t, modelp, nil, "modelp != nil")
	AssertEqual(t, interpolantp, nil, "interpolantp == nil")

	var r32v1num int32
	var r32v1den uint32
	var r32v2num int32
	var r32v2den uint32

	yapi.GetRational32Value(*modelp, r1, &r32v1num, &r32v1den)
	yapi.GetRational32Value(*modelp, r2, &r32v2num, &r32v2den)

	AssertEqual(t, r32v1num, 7, "r32v1num == 7")
	AssertEqual(t, r32v1den, 2, "r32v1den == 2")
	AssertEqual(t, r32v2num, 5, "r32v2num == 5")
	AssertEqual(t, r32v2den, 1, "r32v2den == 1")

	errstr := yapi.ErrorString()
	AssertEqual(t, errstr, "no error", "errstr == 'no error'")

	yapi.PrintError(os.Stderr)

	yapi.CloseContext(&ctxA)
	yapi.CloseContext(&ctxB)

	yapi.Exit()
}

// unsat
func TestInterpolationContext3(t *testing.T) {
	yapi.Init()

	var cfg yapi.ConfigT

	yapi.InitConfig(&cfg)

	yapi.SetConfig(cfg, "solver-type", "mcsat")
	yapi.SetConfig(cfg, "model-interpolation", "true")

	var ctxA yapi.ContextT
	var ctxB yapi.ContextT

	yapi.InitContext(cfg, &ctxA)
	yapi.InitContext(cfg, &ctxB)

	yapi.CloseConfig(&cfg)

	var params yapi.ParamT
	yapi.InitParamRecord(&params)
	yapi.DefaultParamsForContext(ctxA, params)

	realT := yapi.RealType()
	r1 := yapi.NewUninterpretedTerm(realT)
	yapi.SetTermName(r1, "r1")
	r2 := yapi.NewUninterpretedTerm(realT)
	yapi.SetTermName(r2, "r2")

	fmla1 := yapi.ParseTerm("(> r1 3)")
	fmla2 := yapi.ParseTerm("(< r1 4)")
	fmla3 := yapi.ParseTerm("(< (- r1 r2) 0)")

	yapi.AssertFormulas(ctxA, []yapi.TermT{fmla1, fmla2, fmla3})

	fmla4 := yapi.ParseTerm("(< r2 3)")

	yapi.AssertFormula(ctxB, fmla4)

	smtStat, modelp, interpolantp := yapi.CheckContextWithInterpolation(ctxA, ctxB, params, true)

	println(smtStat)

	AssertEqual(t, smtStat, yapi.StatusUnsat, "smtStat == yapi.StatusUnsat")

	AssertEqual(t, modelp, nil, "modelp == nil")

	AssertNotEqual(t, interpolantp, nil, "modelp != nil")

	println(*interpolantp)
	println(yapi.TermToString(*interpolantp, 80, 0, 0))

	AssertEqual(t, yapi.TermToString(*interpolantp, 80, 0, 0), "(>= (+ -3 r2) 0)", "*interpolantp <-> (>= (+ -3 r2) 0)")

	errstr := yapi.ErrorString()
	AssertEqual(t, errstr, "no error", "errstr == 'no error'")

	yapi.PrintError(os.Stderr)

	yapi.CloseContext(&ctxA)
	yapi.CloseContext(&ctxB)

	yapi.Exit()
}
