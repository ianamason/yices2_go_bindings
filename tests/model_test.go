package tests

import (
	"fmt"
	"os"
	"github.com/ianamason/yices2_go_bindings/yices2"
	"testing"
)

// generic start up
func setup() (cfg yices2.Config_t, ctx yices2.Context_t, params yices2.Param_t) {
	yices2.Init()
	yices2.Init_config(&cfg)
	yices2.Init_context(cfg, &ctx)
	yices2.Init_param_record(&params)
	yices2.Default_params_for_context(ctx, params)
	return
}

// clean up a generic startup
func cleanup(cfg *yices2.Config_t, ctx *yices2.Context_t, params *yices2.Param_t){
	yices2.Close_config(cfg)
	yices2.Close_param_record(params)
	yices2.Close_context(ctx)
	yices2.Exit()
}

// sam's helper functions
func parse_assert(fmla_str string, ctx yices2.Context_t) {
	fmla := yices2.Parse_term(fmla_str)
	if fmla != yices2.NULL_TERM {
		yices2.Assert_formula(ctx, fmla)
	}
}

// sam's helper functions
func define_const(name string, typ yices2.Type_t) (term yices2.Term_t) {
	term = yices2.New_uninterpreted_term(typ)
	yices2.Set_term_name(term, name)
	return
}

func Test_bool_models(t *testing.T) {

	cfg, ctx, params := setup()

	bool_t := yices2.Bool_type()
	b1 := define_const("b1", bool_t)
	b2 := define_const("b2", bool_t)
	b3 := define_const("b3", bool_t)
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

	cleanup(&cfg, &ctx, &params)
}

func Test_int_models(t *testing.T) {

	cfg, ctx, params := setup()

	int_t := yices2.Int_type()
	i1 := define_const("i1", int_t)
	i2 := define_const("i2", int_t)
	parse_assert("(> i1 3)", ctx)
	parse_assert("(< i2 i1)", ctx)
	stat := yices2.Check_context(ctx, params)
	AssertEqual(t, stat, yices2.STATUS_SAT, "stat == yices2.STATUS_SAT")
	modelp := yices2.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	var i32v1 int32
	var i32v2 int32
	yices2.Get_int32_value(*modelp, i1, &i32v1)
	yices2.Get_int32_value(*modelp, i2, &i32v2)
	AssertEqual(t, i32v1, 4, "i32v1 == 4")
	AssertEqual(t, i32v2, 3, "i32v2 == 3")
	var i64v1 int64
	var i64v2 int64
	yices2.Get_int64_value(*modelp, i1, &i64v1)
	yices2.Get_int64_value(*modelp, i2, &i64v2)
	AssertEqual(t, i64v1, 4, "i64v1 == 4")
	AssertEqual(t, i64v2, 3, "i64v2 == 3")
	yices2.Print_model(os.Stdout, *modelp)
	yices2.Pp_model(os.Stdout, *modelp, 80, 100, 0)
	mdlstr := yices2.Model_to_string(*modelp, 80, 100, 0)
	AssertEqual(t, mdlstr, "(= i1 4)\n(= i2 3)")

	cleanup(&cfg, &ctx, &params)

}

func Test_rat_models(t *testing.T) {

	cfg, ctx, params := setup()

	real_t := yices2.Real_type()
	r1 := define_const("r1", real_t)
	r2 := define_const("r2", real_t)
	parse_assert("(> r1 3)", ctx)
	parse_assert("(< r1 4)", ctx)
	parse_assert("(< (- r1 r2) 0)", ctx)

	stat := yices2.Check_context(ctx, params)
	AssertEqual(t, stat, yices2.STATUS_SAT, "stat == yices2.STATUS_SAT")
	modelp := yices2.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	var r32v1num int32
	var r32v1den uint32
	var r32v2num int32
	var r32v2den uint32

	yices2.Get_rational32_value(*modelp, r1, &r32v1num, &r32v1den)
	yices2.Get_rational32_value(*modelp, r2, &r32v2num, &r32v2den)

	AssertEqual(t, r32v1num, 7, "r32v1num == 7")
	AssertEqual(t, r32v1den, 2, "r32v1den == 2")
	AssertEqual(t, r32v2num, 4, "r32v2num == 4")
	AssertEqual(t, r32v2den, 1, "r32v2den == 1")

	var r64v1num int64
	var r64v1den uint64
	var r64v2num int64
	var r64v2den uint64

	yices2.Get_rational64_value(*modelp, r1, &r64v1num, &r64v1den)
	yices2.Get_rational64_value(*modelp, r2, &r64v2num, &r64v2den)

	AssertEqual(t, r64v1num, 7, "r64v1num == 7")
	AssertEqual(t, r64v1den, 2, "r64v1den == 2")
	AssertEqual(t, r64v2num, 4, "r64v2num == 4")
	AssertEqual(t, r64v2den, 1, "r64v2den == 1")

	var rdoub1 float64
	var rdoub2 float64

	yices2.Get_double_value(*modelp, r1, &rdoub1)
	yices2.Get_double_value(*modelp, r2, &rdoub2)

	AssertEqual(t, rdoub1, 3.5, "rdoub1 == 3.5")
	AssertEqual(t, rdoub2, 4.0, "rdoub2 == 4.0")

	cleanup(&cfg, &ctx, &params)

}

func Test_mpz_models(t *testing.T) {

	cfg, ctx, params := setup()

	int_t := yices2.Int_type()

	i1 := define_const("i1", int_t)
	i2 := define_const("i2", int_t)

	parse_assert("(> i1 987654321987654321987654321)", ctx)
	parse_assert("(< i2 i1)", ctx)

	stat := yices2.Check_context(ctx, params)
	AssertEqual(t, stat, yices2.STATUS_SAT, "stat == yices2.STATUS_SAT")
	modelp := yices2.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	mstr := yices2.Model_to_string(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= i1 987654321987654321987654322)\n(= i2 987654321987654321987654321)")

	var i32val1 int32
	errcode := yices2.Get_int32_value(*modelp, i1, &i32val1)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yices2.Error_string(), "eval error: the term value does not fit the expected type")
	yerror1 := yices2.NewYicesError()

	yices2.Clear_error()

	var i32val2 int32
	errcode = yices2.Get_int32_value(*modelp, i2, &i32val2)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yices2.Error_string(), "eval error: the term value does not fit the expected type")
	yerror2 := yices2.NewYicesError()

	AssertEqual(t, yerror1, yerror2)

	var mpzval1 yices2.Mpz_t
	errcode = yices2.Get_mpz_value(*modelp, i1, &mpzval1)
	AssertEqual(t, errcode, 0)

	mpz1 := yices2.Mpz(&mpzval1)
	AssertEqual(t, yices2.Term_to_string(mpz1, 200, 10, 0), "987654321987654321987654322")

	var mpzval2 yices2.Mpz_t
	errcode = yices2.Get_mpz_value(*modelp, i2, &mpzval2)
	AssertEqual(t, errcode, 0)

	mpz2 := yices2.Mpz(&mpzval2)
	AssertEqual(t, yices2.Term_to_string(mpz2, 200, 10, 0), "987654321987654321987654321")

	cleanup(&cfg, &ctx, &params)

}


func Test_bv_models(t *testing.T) {

	cfg, ctx, params := setup()

	bv_t := yices2.Bv_type(3)
	bv1 := define_const("bv1", bv_t)
	bv2 := define_const("bv2", bv_t)
	bv3 := define_const("bv3", bv_t)
	parse_assert("(= bv1 (bv-add bv2 bv3))", ctx)
	parse_assert("(bv-gt bv2 0b000)", ctx)
	parse_assert("(bv-gt bv3 0b000)", ctx)

	stat := yices2.Check_context(ctx, params)
	AssertEqual(t, stat, yices2.STATUS_SAT, "stat == yices2.STATUS_SAT")
	modelp := yices2.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	bval1 := []int32{0, 0, 0}
	bval2 := []int32{0, 0, 0}
	bval3 := []int32{0, 0, 0}

	errcode := yices2.Get_bv_value(*modelp, bv1, bval1)
	AssertEqual(t, errcode, 0, "errcode == 0")
	fmt.Printf("bval1 = %v\n", bval1)
	AssertEqual(t, bval1, []int32{0, 0, 0}, "bval1 == []int32{0, 0, 0}")

	errcode = yices2.Get_bv_value(*modelp, bv2, bval2)
	AssertEqual(t, errcode, 0, "errcode == 0")
	fmt.Printf("bval2 = %v\n", bval2)
	AssertEqual(t, bval2, []int32{0, 0, 1}, "bval2 == []int32{0, 0, 1}")

	errcode = yices2.Get_bv_value(*modelp, bv3, bval3)
	AssertEqual(t, errcode, 0, "errcode == 0")
	fmt.Printf("bval3 = %v\n", bval3)
	AssertEqual(t, bval3, []int32{0, 0, 1}, "bval1 == []int32{0, 0, 1}")

	cleanup(&cfg, &ctx, &params)

}

func Test_tuple_models(t *testing.T) {

	cfg, ctx, params := setup()


	bool_t := yices2.Bool_type()
	int_t := yices2.Int_type()
	real_t := yices2.Real_type()
	tup_t := yices2.Tuple_type3(bool_t, real_t, int_t)
	t1 := define_const("t1", tup_t)
	parse_assert("(ite (select t1 1) (< (select t1 2) (select t1 3)) (> (select t1 2) (select t1 3)))", ctx)
	stat := yices2.Check_context(ctx, params)
	AssertEqual(t, stat, yices2.STATUS_SAT, "stat == yices2.STATUS_SAT")
	modelp := yices2.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	mstr := yices2.Model_to_string(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= t1 (mk-tuple false 1 0))")
	var yval yices2.Yval_t
	yices2.Get_value(*modelp, t1, &yval)
	AssertEqual(t, yices2.Get_tag(yval), yices2.YVAL_TUPLE)
	AssertEqual(t, yices2.Val_tuple_arity(*modelp, &yval), 3)

	yvec := make([]yices2.Yval_t, 3)
	yices2.Val_expand_tuple(*modelp, &yval, yvec)
	AssertEqual(t, yices2.Get_tag(yvec[0]), yices2.YVAL_BOOL)
	var bval int32
	var ival int32
	yices2.Val_get_bool(*modelp, &yvec[0], &bval)
	yices2.Val_get_int32(*modelp, &yvec[1], &ival)
	AssertEqual(t, bval, 0)
	AssertEqual(t, ival, 1)

	cleanup(&cfg, &ctx, &params)


}

func TestAlgebraicModels(t *testing.T) {

	yices2.Init()

	if yices2.Has_mcsat() == int32(0) {
		fmt.Println("TestAlgebraicModels skipped because no mcsat.")
		return
	}

	real_t := yices2.Real_type()

	var cfg yices2.Config_t

	var ctx yices2.Context_t

	var params yices2.Param_t

	yices2.Init_config(&cfg)

	yices2.Default_config_for_logic(cfg, "QF_NRA")
    yices2.Set_config(cfg, "mode", "one-shot")

	yices2.Init_context(cfg, &ctx)

	x := define_const("x", real_t)

	parse_assert("(= (* x x) 2)", ctx)

	stat := yices2.Check_context(ctx, params)  //params == NULL in the C

	AssertEqual(t, stat, yices2.STATUS_SAT, "stat == yices2.STATUS_SAT")
	modelp := yices2.Get_model(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	yices2.Print_model(os.Stdout, *modelp)

	var xf float64
	yices2.Get_double_value(*modelp, x, &xf)

	AssertEqual(t, xf, -1.414213562373095, "xf == -1.414213562373095")

	yices2.Close_config(&cfg)

	yices2.Close_context(&ctx)

	yices2.Exit()



}
