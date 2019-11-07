package tests

import (
	"fmt"
	"os"
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"testing"
)

// generic start up
func setup() (cfg yapi.Config_t, ctx yapi.Context_t, params yapi.Param_t) {
	yapi.Init()
	yapi.Init_config(&cfg)
	yapi.Init_context(cfg, &ctx)
	yapi.Init_param_record(&params)
	yapi.Default_params_for_context(ctx, params)
	return
}

// clean up a generic startup
func cleanup(cfg *yapi.Config_t, ctx *yapi.Context_t, params *yapi.Param_t){
	yapi.Close_config(cfg)
	yapi.Close_param_record(params)
	yapi.Close_context(ctx)
	yapi.Exit()
}

// sam's helper functions
func parse_assert(fmla_str string, ctx yapi.Context_t) {
	fmla := yapi.Parse_term(fmla_str)
	if fmla != yapi.NULL_TERM {
		yapi.Assert_formula(ctx, fmla)
	}
}

// sam's helper functions
func define_const(name string, typ yapi.Type_t) (term yapi.Term_t) {
	term = yapi.New_uninterpreted_term(typ)
	yapi.Set_term_name(term, name)
	return
}

func Test_bool_models(t *testing.T) {

	cfg, ctx, params := setup()

	bool_t := yapi.Bool_type()
	b1 := define_const("b1", bool_t)
	b2 := define_const("b2", bool_t)
	b3 := define_const("b3", bool_t)
	b_fml1 := yapi.Parse_term("(or b1 b2 b3)")
	yapi.Assert_formula(ctx, b_fml1)
	stat := yapi.Check_context(ctx, params)
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp := yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	var bval1 int32
	var bval2 int32
	var bval3 int32
	yapi.Get_bool_value(*modelp, b1, &bval1)
	yapi.Get_bool_value(*modelp, b2, &bval2)
	yapi.Get_bool_value(*modelp, b3, &bval3)
	AssertEqual(t, bval1, 0, "bval1 == 0")
	AssertEqual(t, bval2, 0, "bval2 == 0")
	AssertEqual(t, bval3, 1, "bval3 == 1")
	b_fmla2 := yapi.Parse_term("(not b3)")
	yapi.Assert_formula(ctx, b_fmla2)
	stat = yapi.Check_context(ctx, params)
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp = yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	yapi.Get_bool_value(*modelp, b1, &bval1)
	yapi.Get_bool_value(*modelp, b2, &bval2)
	yapi.Get_bool_value(*modelp, b3, &bval3)
	AssertEqual(t, bval1, 0, "bval1 == 0")
	AssertEqual(t, bval2, 1, "bval2 == 1")
	AssertEqual(t, bval3, 0, "bval3 == 0")

	var yval yapi.Yval_t

	yapi.Get_value(*modelp, b1, &yval)
	AssertEqual(t, yapi.Get_tag(yval), yapi.YVAL_BOOL)
	yapi.Val_get_bool(*modelp, &yval, &bval1)
	AssertEqual(t, bval1, 0, "bval1 == 0")

	cleanup(&cfg, &ctx, &params)
}

func Test_int_models(t *testing.T) {

	cfg, ctx, params := setup()

	int_t := yapi.Int_type()
	i1 := define_const("i1", int_t)
	i2 := define_const("i2", int_t)
	parse_assert("(> i1 3)", ctx)
	parse_assert("(< i2 i1)", ctx)
	stat := yapi.Check_context(ctx, params)
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp := yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	var i32v1 int32
	var i32v2 int32
	yapi.Get_int32_value(*modelp, i1, &i32v1)
	yapi.Get_int32_value(*modelp, i2, &i32v2)
	AssertEqual(t, i32v1, 4, "i32v1 == 4")
	AssertEqual(t, i32v2, 3, "i32v2 == 3")
	var i64v1 int64
	var i64v2 int64
	yapi.Get_int64_value(*modelp, i1, &i64v1)
	yapi.Get_int64_value(*modelp, i2, &i64v2)
	AssertEqual(t, i64v1, 4, "i64v1 == 4")
	AssertEqual(t, i64v2, 3, "i64v2 == 3")
	yapi.Print_model(os.Stdout, *modelp)
	yapi.Pp_model(os.Stdout, *modelp, 80, 100, 0)
	mdlstr := yapi.Model_to_string(*modelp, 80, 100, 0)
	AssertEqual(t, mdlstr, "(= i1 4)\n(= i2 3)")

	cleanup(&cfg, &ctx, &params)

}

func Test_rat_models(t *testing.T) {

	cfg, ctx, params := setup()

	real_t := yapi.Real_type()
	r1 := define_const("r1", real_t)
	r2 := define_const("r2", real_t)
	parse_assert("(> r1 3)", ctx)
	parse_assert("(< r1 4)", ctx)
	parse_assert("(< (- r1 r2) 0)", ctx)

	stat := yapi.Check_context(ctx, params)
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp := yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	var r32v1num int32
	var r32v1den uint32
	var r32v2num int32
	var r32v2den uint32

	yapi.Get_rational32_value(*modelp, r1, &r32v1num, &r32v1den)
	yapi.Get_rational32_value(*modelp, r2, &r32v2num, &r32v2den)

	AssertEqual(t, r32v1num, 7, "r32v1num == 7")
	AssertEqual(t, r32v1den, 2, "r32v1den == 2")
	AssertEqual(t, r32v2num, 4, "r32v2num == 4")
	AssertEqual(t, r32v2den, 1, "r32v2den == 1")

	var r64v1num int64
	var r64v1den uint64
	var r64v2num int64
	var r64v2den uint64

	yapi.Get_rational64_value(*modelp, r1, &r64v1num, &r64v1den)
	yapi.Get_rational64_value(*modelp, r2, &r64v2num, &r64v2den)

	AssertEqual(t, r64v1num, 7, "r64v1num == 7")
	AssertEqual(t, r64v1den, 2, "r64v1den == 2")
	AssertEqual(t, r64v2num, 4, "r64v2num == 4")
	AssertEqual(t, r64v2den, 1, "r64v2den == 1")

	var rdoub1 float64
	var rdoub2 float64

	yapi.Get_double_value(*modelp, r1, &rdoub1)
	yapi.Get_double_value(*modelp, r2, &rdoub2)

	AssertEqual(t, rdoub1, 3.5, "rdoub1 == 3.5")
	AssertEqual(t, rdoub2, 4.0, "rdoub2 == 4.0")

	cleanup(&cfg, &ctx, &params)

}

func Test_mpz_models(t *testing.T) {

	cfg, ctx, params := setup()

	int_t := yapi.Int_type()

	i1 := define_const("i1", int_t)
	i2 := define_const("i2", int_t)

	parse_assert("(> i1 987654321987654321987654321)", ctx)
	parse_assert("(< i2 i1)", ctx)

	stat := yapi.Check_context(ctx, params)
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp := yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	mstr := yapi.Model_to_string(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= i1 987654321987654321987654322)\n(= i2 987654321987654321987654321)")

	var i32val1 int32
	errcode := yapi.Get_int32_value(*modelp, i1, &i32val1)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yapi.Error_string(), "eval error: the term value does not fit the expected type")
	yerror1 := yapi.GetYicesError()

	yapi.Clear_error()

	var i32val2 int32
	errcode = yapi.Get_int32_value(*modelp, i2, &i32val2)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yapi.Error_string(), "eval error: the term value does not fit the expected type")
	yerror2 := yapi.GetYicesError()

	AssertEqual(t, yerror1, yerror2)

	var mpzval1 yapi.Mpz_t
	errcode = yapi.Get_mpz_value(*modelp, i1, &mpzval1)
	AssertEqual(t, errcode, 0)

	mpz1 := yapi.Mpz(&mpzval1)
	AssertEqual(t, yapi.Term_to_string(mpz1, 200, 10, 0), "987654321987654321987654322")

	var mpzval2 yapi.Mpz_t
	errcode = yapi.Get_mpz_value(*modelp, i2, &mpzval2)
	AssertEqual(t, errcode, 0)

	mpz2 := yapi.Mpz(&mpzval2)
	AssertEqual(t, yapi.Term_to_string(mpz2, 200, 10, 0), "987654321987654321987654321")

	cleanup(&cfg, &ctx, &params)

}

func Test_mpq_models(t *testing.T) {

	cfg, ctx, params := setup()

	real_t := yapi.Real_type()

	r1 := define_const("r1", real_t)
	r2 := define_const("r2", real_t)

	parse_assert("(> (* r1 3456666334217777794) 987654321987654321987654321)", ctx)
	parse_assert("(< r2 r1)", ctx)

	stat := yapi.Check_context(ctx, params)
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp := yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	mstr := yapi.Model_to_string(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= r1 987654325444320656205432115/3456666334217777794)\n(= r2 987654321987654321987654321/3456666334217777794)")

	var r32num1 int32
	var r32den1 uint32
	errcode := yapi.Get_rational32_value(*modelp, r1, &r32num1, &r32den1)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yapi.Error_string(), "eval error: the term value does not fit the expected type")
	yerror1 := yapi.GetYicesError()

	var r64num2 int64
	var r64den2 uint64
	errcode = yapi.Get_rational64_value(*modelp, r2, &r64num2, &r64den2)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yapi.Error_string(), "eval error: the term value does not fit the expected type")
	yerror2 := yapi.GetYicesError()

	AssertEqual(t, yerror1, yerror2)

	var mpqval1 yapi.Mpq_t
	errcode = yapi.Get_mpq_value(*modelp, r1, &mpqval1)
	AssertEqual(t, errcode, 0)

	mpq1 := yapi.Mpq(&mpqval1)
	AssertEqual(t, yapi.Term_to_string(mpq1, 200, 10, 0), "987654325444320656205432115/3456666334217777794")

	var mpqval2 yapi.Mpq_t
	errcode = yapi.Get_mpq_value(*modelp, r2, &mpqval2)
	AssertEqual(t, errcode, 0)

	mpq2 := yapi.Mpq(&mpqval2)
	AssertEqual(t, yapi.Term_to_string(mpq2, 200, 10, 0), "987654321987654321987654321/3456666334217777794")


	cleanup(&cfg, &ctx, &params)

}


func Test_algebraic_models(t *testing.T) {
	yapi.Init()
	if yapi.Has_mcsat() == int32(0) {
		fmt.Println("TestAlgebraicModels skipped because no mcsat.")
		return
	}
	real_t := yapi.Real_type()
	var cfg yapi.Config_t
	var ctx yapi.Context_t
	var params yapi.Param_t
	yapi.Init_config(&cfg)
	yapi.Default_config_for_logic(cfg, "QF_NRA")
    yapi.Set_config(cfg, "mode", "one-shot")
	yapi.Init_context(cfg, &ctx)
	x := define_const("x", real_t)
	parse_assert("(= (* x x) 2)", ctx)
	stat := yapi.Check_context(ctx, params)  //params == NULL in the C
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp := yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	yapi.Print_model(os.Stdout, *modelp)
	var xf float64
	yapi.Get_double_value(*modelp, x, &xf)
	AssertEqual(t, xf, -1.414213562373095, "xf == -1.414213562373095")
	yapi.Close_config(&cfg)
	yapi.Close_context(&ctx)
	yapi.Exit()
}


func Test_bv_models(t *testing.T) {

	cfg, ctx, params := setup()

	bv_t := yapi.Bv_type(3)
	bv1 := define_const("bv1", bv_t)
	bv2 := define_const("bv2", bv_t)
	bv3 := define_const("bv3", bv_t)
	parse_assert("(= bv1 (bv-add bv2 bv3))", ctx)
	parse_assert("(bv-gt bv2 0b000)", ctx)
	parse_assert("(bv-gt bv3 0b000)", ctx)

	stat := yapi.Check_context(ctx, params)
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp := yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	bval1 := []int32{0, 0, 0}
	bval2 := []int32{0, 0, 0}
	bval3 := []int32{0, 0, 0}

	errcode := yapi.Get_bv_value(*modelp, bv1, bval1)
	AssertEqual(t, errcode, 0, "errcode == 0")
	fmt.Printf("bval1 = %v\n", bval1)
	AssertEqual(t, bval1, []int32{0, 0, 0}, "bval1 == []int32{0, 0, 0}")

	errcode = yapi.Get_bv_value(*modelp, bv2, bval2)
	AssertEqual(t, errcode, 0, "errcode == 0")
	fmt.Printf("bval2 = %v\n", bval2)
	AssertEqual(t, bval2, []int32{0, 0, 1}, "bval2 == []int32{0, 0, 1}")

	errcode = yapi.Get_bv_value(*modelp, bv3, bval3)
	AssertEqual(t, errcode, 0, "errcode == 0")
	fmt.Printf("bval3 = %v\n", bval3)
	AssertEqual(t, bval3, []int32{0, 0, 1}, "bval1 == []int32{0, 0, 1}")

	cleanup(&cfg, &ctx, &params)

}

func Test_tuple_models(t *testing.T) {

	cfg, ctx, params := setup()


	bool_t := yapi.Bool_type()
	int_t := yapi.Int_type()
	real_t := yapi.Real_type()
	tup_t := yapi.Tuple_type3(bool_t, real_t, int_t)
	t1 := define_const("t1", tup_t)
	parse_assert("(ite (select t1 1) (< (select t1 2) (select t1 3)) (> (select t1 2) (select t1 3)))", ctx)
	stat := yapi.Check_context(ctx, params)
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp := yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	mstr := yapi.Model_to_string(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= t1 (mk-tuple false 1 0))")
	var yval yapi.Yval_t
	yapi.Get_value(*modelp, t1, &yval)
	AssertEqual(t, yapi.Get_tag(yval), yapi.YVAL_TUPLE)
	AssertEqual(t, yapi.Val_tuple_arity(*modelp, &yval), 3)

	yvec := make([]yapi.Yval_t, 3)
	yapi.Val_expand_tuple(*modelp, &yval, yvec)
	AssertEqual(t, yapi.Get_tag(yvec[0]), yapi.YVAL_BOOL)
	var bval int32
	var ival int32
	yapi.Val_get_bool(*modelp, &yvec[0], &bval)
	yapi.Val_get_int32(*modelp, &yvec[1], &ival)
	AssertEqual(t, bval, 0)
	AssertEqual(t, ival, 1)

	cleanup(&cfg, &ctx, &params)


}

func Test_function_models(t *testing.T) {

	cfg, ctx, params := setup()

	bool_t := yapi.Bool_type()
	int_t := yapi.Int_type()
	real_t := yapi.Real_type()
	fun_t := yapi.Function_type3(int_t, bool_t, real_t, real_t)

	fstr := yapi.Type_to_string(fun_t, 100, 80, 0)

	AssertEqual(t, fstr, "(-> int bool real real)")

	fn := define_const("fn", fun_t)
	//i1 :=
	define_const("i1", int_t)
	//b1 :=
	define_const("b1", bool_t)
	//r1 :=
	define_const("r1", real_t)

	parse_assert("(> (fn i1 b1 r1) (fn (+ i1 1) (not b1) (- r1 i1)))", ctx)

	stat := yapi.Check_context(ctx, params)
	AssertEqual(t, stat, yapi.STATUS_SAT, "stat == yapi.STATUS_SAT")
	modelp := yapi.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	mstr := yapi.Model_to_string(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= b1 false)\n(= i1 1463)\n(= r1 -579)\n(function fn\n (type (-> int bool real real))\n (= (fn 1463 false -579) 1)\n (= (fn 1464 true -2042) 0)\n (default 2))")

	var yval yapi.Yval_t
	yapi.Get_value(*modelp, fn, &yval)
	AssertEqual(t, yapi.Get_tag(yval), yapi.YVAL_FUNCTION)
	AssertEqual(t, yapi.Val_function_arity(*modelp, &yval), 3)

	var ydef yapi.Yval_t

	yvec := yapi.Val_expand_function(*modelp, &yval, &ydef)
	AssertNotEqual(t, yvec, nil)
	AssertEqual(t, yapi.Get_tag(ydef), yapi.YVAL_RATIONAL)

	var def32val int32
	yapi.Val_get_int32(*modelp, &ydef, &def32val)
	AssertEqual(t, def32val, 2)
	AssertEqual(t, len(yvec), 2)
	map1 := yvec[0]
	map2 := yvec[1]
	AssertEqual(t, yapi.Get_tag(map1), yapi.YVAL_MAPPING)
	AssertEqual(t, yapi.Get_tag(map2), yapi.YVAL_MAPPING)
	AssertEqual(t, yapi.Val_mapping_arity(*modelp, &map1), 3)
	AssertEqual(t, yapi.Val_mapping_arity(*modelp, &map2), 3)


	cleanup(&cfg, &ctx, &params)

}

func Test_scalar_models(t *testing.T) {

	cfg, ctx, params := setup()

	cleanup(&cfg, &ctx, &params)

}

func Test_yval_numeric_models(t *testing.T) {

	cfg, ctx, params := setup()

	cleanup(&cfg, &ctx, &params)

}

func Test_model_from_map(t *testing.T) {

	cfg, ctx, params := setup()

	cleanup(&cfg, &ctx, &params)

}

func Test_implicant(t *testing.T) {

	cfg, ctx, params := setup()

	cleanup(&cfg, &ctx, &params)

}
