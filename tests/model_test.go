package tests

import (
	"github.com/ianamason/yices2_go_bindings/yices2"
	"testing"
)


func parseStringAndAssert(fmla_str string, ctx yices2.Context_t) {
	fmla := yices2.Parse_term(fmla_str)
	if fmla != yices2.NULL_TERM {
		yices2.Assert_formula(ctx, fmla)
	}
}

func defineConstant(name string, typ yices2.Type_t) (term yices2.Term_t) {
	term = yices2.New_uninterpreted_term(typ)
	yices2.Set_term_name(term, name)
	return
}

func test_bool_models(t *testing.T, ctx yices2.Context_t, params yices2.Param_t) {
	bool_t := yices2.Bool_type()
	//int_t := yices2.Int_type()
	//real_t := yices2.Real_type()
	b1 := defineConstant("b1", bool_t)
	b2 := defineConstant("b2", bool_t)
	b3 := defineConstant("b3", bool_t)
	b_fml1 := yices2.Parse_term("(or b1 b2 b3)")
	yices2.Assert_formula(ctx, b_fml1)
	stat := yices2.Check_context(ctx, params)
	AssertEqual(t, stat, yices2.STATUS_SAT, "stat == yices2.STATUS_SAT")
	modelp := yices2.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	var bval1 int32
	var bval2 int32
	var bval3 int32
	yices2.Get_bool_value(*modelp, b1, &bval1)
	yices2.Get_bool_value(*modelp, b2, &bval2)
	yices2.Get_bool_value(*modelp, b3, &bval3)
	AssertEqual(t, bval1, 0, "bval1 == 0")
	AssertEqual(t, bval2, 0, "bval2 == 0")
	AssertEqual(t, bval3, 1, "bval3 == 1")
	b_fmla2 := yices2.Parse_term("(not b3)")
	yices2.Assert_formula(ctx, b_fmla2)
	stat = yices2.Check_context(ctx, params)
	AssertEqual(t, stat, yices2.STATUS_SAT, "stat == yices2.STATUS_SAT")
	modelp = yices2.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	yices2.Get_bool_value(*modelp, b1, &bval1)
	yices2.Get_bool_value(*modelp, b2, &bval2)
	yices2.Get_bool_value(*modelp, b3, &bval3)
	AssertEqual(t, bval1, 0, "bval1 == 0")
	AssertEqual(t, bval2, 1, "bval2 == 1")
	AssertEqual(t, bval3, 0, "bval3 == 0")

	var yval yices2.Yval_t

	yices2.Get_value(*modelp, b1, &yval)
	AssertEqual(t, yices2.Get_tag(yval), yices2.YVAL_BOOL)
	yices2.Val_get_bool(*modelp, &yval, &bval1)
	AssertEqual(t, bval1, 0, "bval1 == 0")

}

func test_int_models(t *testing.T, ctx yices2.Context_t, params yices2.Param_t) {

}

func test_rat_models(t *testing.T, ctx yices2.Context_t, params yices2.Param_t) {

}

func test_bv_models(t *testing.T, ctx yices2.Context_t, params yices2.Param_t) {

}

func TestModel0(t *testing.T) {

	yices2.Init()

	var cfg yices2.Config_t

	var ctx yices2.Context_t

	var params yices2.Param_t

	yices2.Init_config(&cfg)

	yices2.Init_context(cfg, &ctx)

	yices2.Init_param_record(&params)

	yices2.Default_params_for_context(ctx, params)

	test_bool_models(t, ctx, params)

	test_int_models(t, ctx, params)

	test_rat_models(t, ctx, params)

	test_bv_models(t, ctx, params)

	yices2.Close_config(&cfg)

	yices2.Close_param_record(&params)

	yices2.Close_context(&ctx)

	yices2.Exit()



}
