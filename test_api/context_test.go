package tests

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"testing"
)

func TestContext0(t *testing.T) {

	yapi.Init()

	var cfg yapi.Config_t

	yapi.Init_config(&cfg)

	var ctx yapi.Context_t

	yapi.Init_context(cfg, &ctx)

	yapi.Close_config(&cfg)

	bv_t := yapi.Bv_type(3)
	bvvar1 := yapi.New_uninterpreted_term(bv_t)
	yapi.Set_term_name(bvvar1, "x")
	bvvar2 := yapi.New_uninterpreted_term(bv_t)
	yapi.Set_term_name(bvvar2, "y")
	bvvar3 := yapi.New_uninterpreted_term(bv_t)
	yapi.Set_term_name(bvvar3, "z")
	fmla1 := yapi.Parse_term("(= x (bv-add y z))")
	fmla2 := yapi.Parse_term("(bv-gt y 0b000)")
	fmla3 := yapi.Parse_term("(bv-gt z 0b000)")
	yapi.Assert_formula(ctx, fmla1)
	yapi.Assert_formulas(ctx, []yapi.Term_t{fmla1, fmla2, fmla3})
	var params yapi.Param_t
	smt_stat := yapi.Check_context(ctx, params)
	AssertEqual(t, smt_stat, yapi.STATUS_SAT, "smt_stat == yapi.STATUS_SAT")

	yapi.Init_param_record(&params)
	yapi.Default_params_for_context(ctx, params)

	errcode := yapi.Set_param(params, "dyn-ack", "true")
	AssertEqual(t, errcode, 0, "errcode == 0") //FIXME: is this right?

	yapi.Close_param_record(&params)

	yapi.Close_context(&ctx)

	yapi.Exit()

}

func TestContext1(t *testing.T) {
	yapi.Init()

	var cfg yapi.Config_t

	var ctx yapi.Context_t

	yapi.Init_config(&cfg)

	yapi.Init_context(cfg, &ctx)

	stat := yapi.Context_status(ctx)
	ret := yapi.Push(ctx)
	AssertEqual(t, ret, 0, "ret == 0")
	ret = yapi.Pop(ctx)
	AssertEqual(t, ret, 0, "ret == 0")
	yapi.Reset_context(ctx)
	ret = yapi.Context_enable_option(ctx, "arith-elim")
	AssertEqual(t, ret, 0, "ret == 0")
	ret = yapi.Context_disable_option(ctx, "arith-elim")
	AssertEqual(t, ret, 0, "ret == 0")
	stat = yapi.Context_status(ctx)
	AssertEqual(t, stat, yapi.STATUS_IDLE, "stat == yapi.STATUS_IDLE")
	yapi.Reset_context(ctx)
	bool_t := yapi.Bool_type()
	bvar1 := yapi.New_variable(bool_t)
	errcode := yapi.Assert_formula(ctx, bvar1)
	error_string := yapi.Error_string()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, error_string, "assertion contains a free variable", "error_string == 'assertion contains a free variable'")
	bv_t := yapi.Bv_type(3)
	bvvar1 := yapi.New_uninterpreted_term(bv_t)
	yapi.Set_term_name(bvvar1, "x")
	bvvar2 := yapi.New_uninterpreted_term(bv_t)
	yapi.Set_term_name(bvvar2, "y")
	bvvar3 := yapi.New_uninterpreted_term(bv_t)
	yapi.Set_term_name(bvvar3, "z")
	fmla1 := yapi.Parse_term("(= x (bv-add y z))")
	fmla2 := yapi.Parse_term("(bv-gt y 0b000)")
	fmla3 := yapi.Parse_term("(bv-gt z 0b000)")
	yapi.Assert_formula(ctx, fmla1)
	yapi.Assert_formulas(ctx, []yapi.Term_t{fmla1, fmla2, fmla3})

	var params yapi.Param_t
	smt_stat := yapi.Check_context(ctx, params) //same as passing NULL to the C
	AssertEqual(t, smt_stat, yapi.STATUS_SAT, "smt_stat == yapi.STATUS_SAT")
	yapi.Assert_blocking_clause(ctx)
	yapi.Stop_search(ctx)

	yapi.Init_param_record(&params)
	yapi.Default_params_for_context(ctx, params)
	yapi.Set_param(params, "dyn-ack", "true")
	errcode = yapi.Set_param(params, "foo", "bar")
	error_string = yapi.Error_string()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, error_string, "invalid parameter", "error_string == 'invalid parameter'")
	errcode = yapi.Set_param(params, "dyn-ack", "bar")
	error_string = yapi.Error_string()
	AssertEqual(t, errcode, -1, "errcode == -1")
	AssertEqual(t, error_string, "value not valid for parameter", "error_string == 'value not valid for parameter'")
	yapi.Close_param_record(&params)

	yapi.Close_context(&ctx)

	yapi.Exit()

}
