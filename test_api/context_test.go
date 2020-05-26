package tests

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"testing"
)

func TestContext0(t *testing.T) {

	yapi.Init()

	var cfg yapi.ConfigT

	yapi.InitConfig(&cfg)

	var ctx yapi.ContextT

	yapi.InitContext(cfg, &ctx)

	yapi.CloseConfig(&cfg)

	bvT := yapi.BvType(3)
	bvvar1 := yapi.NewUninterpretedTerm(bvT)
	yapi.SetTermName(bvvar1, "x")
	bvvar2 := yapi.NewUninterpretedTerm(bvT)
	yapi.SetTermName(bvvar2, "y")
	bvvar3 := yapi.NewUninterpretedTerm(bvT)
	yapi.SetTermName(bvvar3, "z")
	fmla1 := yapi.ParseTerm("(= x (bv-add y z))")
	fmla2 := yapi.ParseTerm("(bv-gt y 0b000)")
	fmla3 := yapi.ParseTerm("(bv-gt z 0b000)")
	yapi.AssertFormula(ctx, fmla1)
	yapi.AssertFormulas(ctx, []yapi.TermT{fmla1, fmla2, fmla3})
	var params yapi.ParamT
	smtStat := yapi.CheckContext(ctx, params)
	AssertEqual(t, smtStat, yapi.StatusSat, "smtStat == yapi.StatusSat")

	yapi.InitParamRecord(&params)
	yapi.DefaultParamsForContext(ctx, params)

	errcode := yapi.SetParam(params, "dyn-ack", "true")
	AssertEqual(t, errcode, 0, "errcode == 0") //FIXME: is this right?

	yapi.CloseParamRecord(&params)

	yapi.CloseContext(&ctx)

	yapi.Exit()

}

func TestContext1(t *testing.T) {
	yapi.Init()

	var cfg yapi.ConfigT

	var ctx yapi.ContextT

	yapi.InitConfig(&cfg)

	yapi.InitContext(cfg, &ctx)

	yapi.ContextStatus(ctx)
	ret := yapi.Push(ctx)
	AssertEqual(t, ret, 0, "ret == 0")
	ret = yapi.Pop(ctx)
	AssertEqual(t, ret, 0, "ret == 0")
	yapi.ResetContext(ctx)
	ret = yapi.ContextEnableOption(ctx, "arith-elim")
	AssertEqual(t, ret, 0, "ret == 0")
	ret = yapi.ContextDisableOption(ctx, "arith-elim")
	AssertEqual(t, ret, 0, "ret == 0")
	stat := yapi.ContextStatus(ctx)
	AssertEqual(t, stat, yapi.StatusIdle, "stat == yapi.StatusIdle")
	yapi.ResetContext(ctx)
	boolT := yapi.BoolType()
	bvar1 := yapi.NewVariable(boolT)
	errcode := yapi.AssertFormula(ctx, bvar1)
	errorString := yapi.ErrorString()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, errorString, "assertion contains a free variable", "errorString == 'assertion contains a free variable'")
	bvT := yapi.BvType(3)
	bvvar1 := yapi.NewUninterpretedTerm(bvT)
	yapi.SetTermName(bvvar1, "x")
	bvvar2 := yapi.NewUninterpretedTerm(bvT)
	yapi.SetTermName(bvvar2, "y")
	bvvar3 := yapi.NewUninterpretedTerm(bvT)
	yapi.SetTermName(bvvar3, "z")
	fmla1 := yapi.ParseTerm("(= x (bv-add y z))")
	fmla2 := yapi.ParseTerm("(bv-gt y 0b000)")
	fmla3 := yapi.ParseTerm("(bv-gt z 0b000)")
	yapi.AssertFormula(ctx, fmla1)
	yapi.AssertFormulas(ctx, []yapi.TermT{fmla1, fmla2, fmla3})

	var params yapi.ParamT
	smtStat := yapi.CheckContext(ctx, params) //same as passing NULL to the C
	AssertEqual(t, smtStat, yapi.StatusSat, "smtStat == yapi.StatusSat")
	yapi.AssertBlockingClause(ctx)
	yapi.StopSearch(ctx)

	yapi.InitParamRecord(&params)
	yapi.DefaultParamsForContext(ctx, params)
	yapi.SetParam(params, "dyn-ack", "true")
	errcode = yapi.SetParam(params, "foo", "bar")
	errorString = yapi.ErrorString()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, errorString, "invalid parameter", "errorString == 'invalid parameter'")
	errcode = yapi.SetParam(params, "dyn-ack", "bar")
	errorString = yapi.ErrorString()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, errorString, "value not valid for parameter", "errorString == 'value not valid for parameter'")
	yapi.CloseParamRecord(&params)

	yapi.CloseContext(&ctx)

	yapi.Exit()

}
