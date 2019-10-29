package tests

import (
	"github.com/ianamason/yices2_go_bindings/yices2"
	"testing"
)


func TestContext0(t *testing.T) {

	yices2.Init()

	var cfg yices2.Config_t

	yices2.Init_config(&cfg)

	var ctx yices2.Context_t

	yices2.Init_context(cfg, &ctx)

	yices2.Close_config(&cfg)

	bv_t := yices2.Bv_type(3)
	bvvar1 := yices2.New_uninterpreted_term(bv_t)
	yices2.Set_term_name(bvvar1, "x")
	bvvar2 := yices2.New_uninterpreted_term(bv_t)
	yices2.Set_term_name(bvvar2, "y")
	bvvar3 := yices2.New_uninterpreted_term(bv_t)
	yices2.Set_term_name(bvvar3, "z")
	fmla1 := yices2.Parse_term("(= x (bv-add y z))")
	fmla2 := yices2.Parse_term("(bv-gt y 0b000)")
	fmla3 := yices2.Parse_term("(bv-gt z 0b000)")
	yices2.Assert_formula(ctx, fmla1)
	yices2.Assert_formulas(ctx, []yices2.Term_t{fmla1, fmla2, fmla3})
	var params yices2.Param_t
	smt_stat := yices2.Check_context(ctx,params)
	AssertEqual(t, smt_stat, yices2.STATUS_SAT)

	yices2.Init_param_record(&params)
	yices2.Default_params_for_context(ctx, params)

	errcode := yices2.Set_param(params, "dyn-ack", "true")
	AssertEqual(t, errcode, 0)  //FIXME: is this right?
	
	yices2.Close_param_record(&params)

	yices2.Close_context(&ctx)


	yices2.Exit()

}

func NotATestContext1(t *testing.T) {
	yices2.Init()

	var cfg yices2.Config_t

	var ctx yices2.Context_t

	yices2.Init_config(&cfg)

	yices2.Init_context(cfg, &ctx)

	stat := yices2.Context_status(ctx)
	ret := yices2.Push(ctx)
	AssertEqual(t, ret, 0)
	ret = yices2.Pop(ctx)
	AssertEqual(t, ret, 0)
	yices2.Reset_context(ctx)
	ret = yices2.Context_enable_option(ctx, "arith-elim")
	AssertEqual(t, ret, 0)
	ret = yices2.Context_disable_option(ctx, "arith-elim")
	AssertEqual(t, ret, 0)
	stat = yices2.Context_status(ctx)
	AssertEqual(t, stat, yices2.STATUS_IDLE)
	yices2.Reset_context(ctx)
	bool_t := yices2.Bool_type()
	bvar1 := yices2.New_variable(bool_t)
	errcode := yices2.Assert_formula(ctx, bvar1)
	error_string := yices2.Error_string()
	AssertEqual(t, errcode, -1)
	AssertEqual(t, error_string, "assertion contains a free variable")
	bv_t := yices2.Bv_type(3)
	bvvar1 := yices2.New_uninterpreted_term(bv_t)
	yices2.Set_term_name(bvvar1, "x")
	bvvar2 := yices2.New_uninterpreted_term(bv_t)
	yices2.Set_term_name(bvvar2, "y")
	bvvar3 := yices2.New_uninterpreted_term(bv_t)
	yices2.Set_term_name(bvvar3, "z")
	fmla1 := yices2.Parse_term("(= x (bv-add y z))")
	fmla2 := yices2.Parse_term("(bv-gt y 0b000)")
	fmla3 := yices2.Parse_term("(bv-gt z 0b000)")
	yices2.Assert_formula(ctx, fmla1)
	yices2.Assert_formulas(ctx, []yices2.Term_t{fmla1, fmla2, fmla3})

	var params yices2.Param_t
	smt_stat := yices2.Check_context(ctx, params)  //same as passing NULL to the C
	AssertEqual(t, smt_stat, yices2.STATUS_SAT)
	yices2.Assert_blocking_clause(ctx)
	yices2.Stop_search(ctx)

	yices2.Init_param_record(&params)
	yices2.Default_params_for_context(ctx, params)
	yices2.Set_param(params, "dyn-ack", "true")
	errcode = yices2.Set_param(params, "foo", "bar")
	error_string = yices2.Error_string()
	AssertEqual(t, errcode, -1)
	AssertEqual(t, error_string, "invalid parameter")
	errcode = yices2.Set_param(params, "dyn-ack", "bar")
	error_string = yices2.Error_string()
	AssertEqual(t, errcode, -1)
	AssertEqual(t, error_string, "value not valid for parameter")
	yices2.Close_param_record(&params)

	yices2.Close_context(&ctx)


	yices2.Exit()

}

