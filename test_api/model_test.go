package tests

import (
	"fmt"
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"os"
	"testing"
)

// generic start up
func setup() (cfg yapi.ConfigT, ctx yapi.ContextT, params yapi.ParamT) {
	yapi.Init()
	yapi.InitConfig(&cfg)
	yapi.InitContext(cfg, &ctx)
	yapi.InitParamRecord(&params)
	yapi.DefaultParamsForContext(ctx, params)
	return
}

// clean up a generic startup
func cleanup(cfg *yapi.ConfigT, ctx *yapi.ContextT, params *yapi.ParamT) {
	yapi.CloseConfig(cfg)
	yapi.CloseParamRecord(params)
	yapi.CloseContext(ctx)
	yapi.Exit()
}

// sam's helper functions
func parseAssert(fmlaStr string, ctx yapi.ContextT) {
	fmla := yapi.ParseTerm(fmlaStr)
	if fmla != yapi.NullTerm {
		yapi.AssertFormula(ctx, fmla)
	}
}

// sam's helper functions
func defineConst(name string, typ yapi.TypeT) (term yapi.TermT) {
	term = yapi.NewUninterpretedTerm(typ)
	yapi.SetTermName(term, name)
	return
}

func TestBoolModels(t *testing.T) {

	cfg, ctx, params := setup()

	boolT := yapi.BoolType()
	b1 := defineConst("b1", boolT)
	b2 := defineConst("b2", boolT)
	b3 := defineConst("b3", boolT)
	bFml1 := yapi.ParseTerm("(or b1 b2 b3)")
	yapi.AssertFormula(ctx, bFml1)
	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	var bval1 int32
	var bval2 int32
	var bval3 int32
	yapi.GetBoolValue(*modelp, b1, &bval1)
	yapi.GetBoolValue(*modelp, b2, &bval2)
	yapi.GetBoolValue(*modelp, b3, &bval3)
	AssertEqual(t, bval1, 0, "bval1 == 0")
	AssertEqual(t, bval2, 0, "bval2 == 0")
	AssertEqual(t, bval3, 1, "bval3 == 1")
	bFmla2 := yapi.ParseTerm("(not b3)")

	yapi.CloseModel(modelp)

	yapi.AssertFormula(ctx, bFmla2)
	stat = yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp = yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	yapi.GetBoolValue(*modelp, b1, &bval1)
	yapi.GetBoolValue(*modelp, b2, &bval2)
	yapi.GetBoolValue(*modelp, b3, &bval3)
	AssertEqual(t, bval1, 0, "bval1 == 0")
	AssertEqual(t, bval2, 1, "bval2 == 1")
	AssertEqual(t, bval3, 0, "bval3 == 0")

	var yval yapi.YvalT

	yapi.GetValue(*modelp, b1, &yval)
	AssertEqual(t, yapi.GetTag(yval), yapi.YvalBool)
	yapi.ValGetBool(*modelp, &yval, &bval1)
	AssertEqual(t, bval1, 0, "bval1 == 0")

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)
}

func TestIntModels(t *testing.T) {

	cfg, ctx, params := setup()

	intT := yapi.IntType()
	i1 := defineConst("i1", intT)
	i2 := defineConst("i2", intT)
	parseAssert("(> i1 3)", ctx)
	parseAssert("(< i2 i1)", ctx)
	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	var i32v1 int32
	var i32v2 int32
	yapi.GetInt32Value(*modelp, i1, &i32v1)
	yapi.GetInt32Value(*modelp, i2, &i32v2)
	AssertEqual(t, i32v1, 4, "i32v1 == 4")
	AssertEqual(t, i32v2, 3, "i32v2 == 3")
	var i64v1 int64
	var i64v2 int64
	yapi.GetInt64Value(*modelp, i1, &i64v1)
	yapi.GetInt64Value(*modelp, i2, &i64v2)
	AssertEqual(t, i64v1, 4, "i64v1 == 4")
	AssertEqual(t, i64v2, 3, "i64v2 == 3")
	yapi.PrintModel(os.Stdout, *modelp)
	yapi.PpModel(os.Stdout, *modelp, 80, 100, 0)
	mdlstr := yapi.ModelToString(*modelp, 80, 100, 0)
	AssertEqual(t, mdlstr, "(= i1 4)\n(= i2 3)")

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func Test_rat_models(t *testing.T) {

	cfg, ctx, params := setup()

	realT := yapi.RealType()
	r1 := defineConst("r1", realT)
	r2 := defineConst("r2", realT)
	parseAssert("(> r1 3)", ctx)
	parseAssert("(< r1 4)", ctx)
	parseAssert("(< (- r1 r2) 0)", ctx)

	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	var r32v1num int32
	var r32v1den uint32
	var r32v2num int32
	var r32v2den uint32

	yapi.GetRational32Value(*modelp, r1, &r32v1num, &r32v1den)
	yapi.GetRational32Value(*modelp, r2, &r32v2num, &r32v2den)

	AssertEqual(t, r32v1num, 7, "r32v1num == 7")
	AssertEqual(t, r32v1den, 2, "r32v1den == 2")
	AssertEqual(t, r32v2num, 4, "r32v2num == 4")
	AssertEqual(t, r32v2den, 1, "r32v2den == 1")

	var r64v1num int64
	var r64v1den uint64
	var r64v2num int64
	var r64v2den uint64

	yapi.GetRational64Value(*modelp, r1, &r64v1num, &r64v1den)
	yapi.GetRational64Value(*modelp, r2, &r64v2num, &r64v2den)

	AssertEqual(t, r64v1num, 7, "r64v1num == 7")
	AssertEqual(t, r64v1den, 2, "r64v1den == 2")
	AssertEqual(t, r64v2num, 4, "r64v2num == 4")
	AssertEqual(t, r64v2den, 1, "r64v2den == 1")

	var rdoub1 float64
	var rdoub2 float64

	yapi.GetDoubleValue(*modelp, r1, &rdoub1)
	yapi.GetDoubleValue(*modelp, r2, &rdoub2)

	AssertEqual(t, rdoub1, 3.5, "rdoub1 == 3.5")
	AssertEqual(t, rdoub2, 4.0, "rdoub2 == 4.0")

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func TestMpzModels(t *testing.T) {

	cfg, ctx, params := setup()

	intT := yapi.IntType()

	i1 := defineConst("i1", intT)
	i2 := defineConst("i2", intT)

	parseAssert("(> i1 987654321987654321987654321)", ctx)
	parseAssert("(< i2 i1)", ctx)

	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	mstr := yapi.ModelToString(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= i1 987654321987654321987654322)\n(= i2 987654321987654321987654321)")

	var i32val1 int32
	errcode := yapi.GetInt32Value(*modelp, i1, &i32val1)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yapi.ErrorString(), "eval error: the term value does not fit the expected type")
	yerror1 := yapi.YicesError()

	yapi.ClearError()

	var i32val2 int32
	errcode = yapi.GetInt32Value(*modelp, i2, &i32val2)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yapi.ErrorString(), "eval error: the term value does not fit the expected type")
	yerror2 := yapi.YicesError()

	AssertEqual(t, yerror1, yerror2)

	var mpzval1 yapi.MpzT
	errcode = yapi.GetMpzValue(*modelp, i1, &mpzval1)
	AssertEqual(t, errcode, 0)

	mpz1 := yapi.Mpz(&mpzval1)
	AssertEqual(t, yapi.TermToString(mpz1, 200, 10, 0), "987654321987654321987654322")

	var mpzval2 yapi.MpzT
	errcode = yapi.GetMpzValue(*modelp, i2, &mpzval2)
	AssertEqual(t, errcode, 0)

	mpz2 := yapi.Mpz(&mpzval2)
	AssertEqual(t, yapi.TermToString(mpz2, 200, 10, 0), "987654321987654321987654321")

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func TestMpqModels(t *testing.T) {

	cfg, ctx, params := setup()

	realT := yapi.RealType()

	r1 := defineConst("r1", realT)
	r2 := defineConst("r2", realT)

	parseAssert("(> (* r1 3456666334217777794) 987654321987654321987654321)", ctx)
	parseAssert("(< r2 r1)", ctx)

	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	mstr := yapi.ModelToString(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= r1 987654325444320656205432115/3456666334217777794)\n(= r2 987654321987654321987654321/3456666334217777794)")

	var r32num1 int32
	var r32den1 uint32
	errcode := yapi.GetRational32Value(*modelp, r1, &r32num1, &r32den1)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yapi.ErrorString(), "eval error: the term value does not fit the expected type")
	yerror1 := yapi.YicesError()

	var r64num2 int64
	var r64den2 uint64
	errcode = yapi.GetRational64Value(*modelp, r2, &r64num2, &r64den2)
	AssertEqual(t, errcode, -1)
	AssertEqual(t, yapi.ErrorString(), "eval error: the term value does not fit the expected type")
	yerror2 := yapi.YicesError()

	AssertEqual(t, yerror1, yerror2)

	var mpqval1 yapi.MpqT
	errcode = yapi.GetMpqValue(*modelp, r1, &mpqval1)
	AssertEqual(t, errcode, 0)

	mpq1 := yapi.Mpq(&mpqval1)
	AssertEqual(t, yapi.TermToString(mpq1, 200, 10, 0), "987654325444320656205432115/3456666334217777794")

	var mpqval2 yapi.MpqT
	errcode = yapi.GetMpqValue(*modelp, r2, &mpqval2)
	AssertEqual(t, errcode, 0)

	mpq2 := yapi.Mpq(&mpqval2)
	AssertEqual(t, yapi.TermToString(mpq2, 200, 10, 0), "987654321987654321987654321/3456666334217777794")

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func TestAlgebraicModels(t *testing.T) {
	yapi.Init()
	if yapi.HasMcsat() == int32(0) {
		fmt.Println("TestAlgebraicModels skipped because no mcsat.")
		return
	}
	realT := yapi.RealType()
	var cfg yapi.ConfigT
	var ctx yapi.ContextT
	var params yapi.ParamT
	yapi.InitConfig(&cfg)
	yapi.DefaultConfigForLogic(cfg, "QF_NRA")
	yapi.SetConfig(cfg, "mode", "one-shot")
	yapi.InitContext(cfg, &ctx)
	x := defineConst("x", realT)
	parseAssert("(= (* x x) 2)", ctx)
	stat := yapi.CheckContext(ctx, params) //params == NULL in the C
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	yapi.PrintModel(os.Stdout, *modelp)
	var xf float64
	yapi.GetDoubleValue(*modelp, x, &xf)
	AssertEqual(t, xf, -1.414213562373095, "xf == -1.414213562373095")
	yapi.CloseModel(modelp)
	yapi.CloseConfig(&cfg)
	yapi.CloseContext(&ctx)
	yapi.Exit()
}

func TestBvModels(t *testing.T) {

	cfg, ctx, params := setup()

	bvT := yapi.BvType(3)
	bv1 := defineConst("bv1", bvT)
	bv2 := defineConst("bv2", bvT)
	bv3 := defineConst("bv3", bvT)
	parseAssert("(= bv1 (bv-add bv2 bv3))", ctx)
	parseAssert("(bv-gt bv2 0b000)", ctx)
	parseAssert("(bv-gt bv3 0b000)", ctx)

	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	bval1 := []int32{0, 0, 0}
	bval2 := []int32{0, 0, 0}
	bval3 := []int32{0, 0, 0}

	errcode := yapi.GetBvValue(*modelp, bv1, bval1)
	AssertEqual(t, errcode, 0, "errcode == 0")
	fmt.Printf("bval1 = %v\n", bval1)
	AssertEqual(t, bval1, []int32{0, 0, 0}, "bval1 == []int32{0, 0, 0}")

	errcode = yapi.GetBvValue(*modelp, bv2, bval2)
	AssertEqual(t, errcode, 0, "errcode == 0")
	fmt.Printf("bval2 = %v\n", bval2)
	AssertEqual(t, bval2, []int32{0, 0, 1}, "bval2 == []int32{0, 0, 1}")

	errcode = yapi.GetBvValue(*modelp, bv3, bval3)
	AssertEqual(t, errcode, 0, "errcode == 0")
	fmt.Printf("bval3 = %v\n", bval3)
	AssertEqual(t, bval3, []int32{0, 0, 1}, "bval1 == []int32{0, 0, 1}")

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func TestTupleModels(t *testing.T) {

	cfg, ctx, params := setup()

	boolT := yapi.BoolType()
	intT := yapi.IntType()
	realT := yapi.RealType()
	tupT := yapi.TupleType3(boolT, realT, intT)
	t1 := defineConst("t1", tupT)
	parseAssert("(ite (select t1 1) (< (select t1 2) (select t1 3)) (> (select t1 2) (select t1 3)))", ctx)
	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")
	mstr := yapi.ModelToString(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= t1 (mk-tuple false 1 0))")
	var yval yapi.YvalT
	yapi.GetValue(*modelp, t1, &yval)
	AssertEqual(t, yapi.GetTag(yval), yapi.YvalTuple)
	AssertEqual(t, yapi.ValTupleArity(*modelp, &yval), 3)

	yvec := make([]yapi.YvalT, 3)
	yapi.ValExpandTuple(*modelp, &yval, yvec)
	AssertEqual(t, yapi.GetTag(yvec[0]), yapi.YvalBool)
	var bval int32
	var ival int32
	yapi.ValGetBool(*modelp, &yvec[0], &bval)
	yapi.ValGetInt32(*modelp, &yvec[1], &ival)
	AssertEqual(t, bval, 0)
	AssertEqual(t, ival, 1)

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func Test_function_models(t *testing.T) {

	cfg, ctx, params := setup()

	boolT := yapi.BoolType()
	intT := yapi.IntType()
	realT := yapi.RealType()
	funT := yapi.FunctionType3(intT, boolT, realT, realT)

	fstr := yapi.TypeToString(funT, 100, 80, 0)

	AssertEqual(t, fstr, "(-> int bool real real)")

	fn := defineConst("fn", funT)
	i1 := defineConst("i1", intT)
	//b1 :=
	defineConst("b1", boolT)
	r1 := defineConst("r1", realT)

	parseAssert("(> (fn i1 b1 r1) (fn (+ i1 1) (not b1) (- r1 i1)))", ctx)

	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	mstr := yapi.ModelToString(*modelp, 80, 100, 0)
	AssertEqual(t, mstr, "(= b1 false)\n(= i1 1463)\n(= r1 -579)\n(function fn\n (type (-> int bool real real))\n (= (fn 1463 false -579) 1)\n (= (fn 1464 true -2042) 0)\n (default 2))")

	var yval yapi.YvalT
	yapi.GetValue(*modelp, fn, &yval)
	AssertEqual(t, yapi.GetTag(yval), yapi.YvalFunction)
	AssertEqual(t, yapi.ValFunctionArity(*modelp, &yval), 3)

	var ydef yapi.YvalT

	yvec := yapi.ValExpandFunction(*modelp, &yval, &ydef)
	AssertNotEqual(t, yvec, nil)
	AssertEqual(t, yapi.GetTag(ydef), yapi.YvalRational)

	var def32val int32
	yapi.ValGetInt32(*modelp, &ydef, &def32val)
	AssertEqual(t, def32val, 2)
	AssertEqual(t, len(yvec), 2)
	map1 := yvec[0]
	map2 := yvec[1]
	AssertEqual(t, yapi.GetTag(map1), yapi.YvalMapping)
	AssertEqual(t, yapi.GetTag(map2), yapi.YvalMapping)
	AssertEqual(t, yapi.ValMappingArity(*modelp, &map1), 3)
	AssertEqual(t, yapi.ValMappingArity(*modelp, &map2), 3)

	var yval1, yval2 yapi.YvalT

	vec1 := yapi.ValExpandMapping(*modelp, &map1, &yval1)
	vec2 := yapi.ValExpandMapping(*modelp, &map2, &yval2)

	AssertEqual(t, yapi.GetTag(yval1), yapi.YvalRational)
	AssertEqual(t, yapi.GetTag(yval2), yapi.YvalRational)

	AssertEqual(t, len(vec1), 3)
	AssertEqual(t, len(vec2), 3)

	var val1, val2 int32

	yapi.ValGetInt32(*modelp, &yval1, &val1)
	yapi.ValGetInt32(*modelp, &yval2, &val2)

	AssertEqual(t, val1, 1)
	AssertEqual(t, val2, 0)

	var arg10, arg20 int32

	yapi.ValGetInt32(*modelp, &vec1[0], &arg10)
	yapi.ValGetInt32(*modelp, &vec2[0], &arg20)
	AssertEqual(t, arg10, 1463)
	AssertEqual(t, arg20, 1464)

	var arg11, arg21 int32

	yapi.ValGetBool(*modelp, &vec1[1], &arg11)
	yapi.ValGetBool(*modelp, &vec2[1], &arg21)
	AssertEqual(t, arg11, 0)
	AssertEqual(t, arg21, 1)

	var arg12, arg22 int32

	yapi.ValGetInt32(*modelp, &vec1[2], &arg12)
	yapi.ValGetInt32(*modelp, &vec2[2], &arg22)
	AssertEqual(t, arg12, -579)
	AssertEqual(t, arg22, -2042)

	fmla := yapi.ParseTerm("(> i1 r1)")

	AssertEqual(t, yapi.FormulaTrueInModel(*modelp, fmla), 1)

	aArr := []yapi.TermT{i1, fmla, r1}
	bArr := make([]yapi.TermT, 3)

	errcode := yapi.TermArrayValue(*modelp, aArr, bArr)

	AssertEqual(t, errcode, 0)

	AssertEqual(t, bArr[0], yapi.Int32(1463))
	AssertEqual(t, bArr[1], yapi.True())
	AssertEqual(t, bArr[2], yapi.Int32(-579))

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func TestScalarModels(t *testing.T) {

	cfg, ctx, params := setup()

	scalarT := yapi.NewScalarType(10)

	sc1 := defineConst("sc1", scalarT)
	sc2 := defineConst("sc2", scalarT)
	sc3 := defineConst("sc3", scalarT)

	parseAssert("(/= sc1 sc2)", ctx)
	parseAssert("(/= sc1 sc3)", ctx)

	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	var val1, val2, val3 int32

	yapi.GetScalarValue(*modelp, sc1, &val1)
	yapi.GetScalarValue(*modelp, sc2, &val2)
	yapi.GetScalarValue(*modelp, sc3, &val3)

	AssertEqual(t, sc1, yapi.TermT(6))
	AssertEqual(t, sc2, yapi.TermT(8))
	AssertEqual(t, sc3, yapi.TermT(10))

	AssertEqual(t, yapi.TermIsScalar(sc1), true)
	AssertEqual(t, yapi.TermIsScalar(sc2), true)
	AssertEqual(t, yapi.TermIsScalar(sc3), true)

	var yval1, yval2, yval3 yapi.YvalT

	AssertEqual(t, yapi.GetValue(*modelp, sc1, &yval1), 0)
	AssertEqual(t, yapi.GetValue(*modelp, sc2, &yval2), 0)
	AssertEqual(t, yapi.GetValue(*modelp, sc3, &yval3), 0)

	AssertEqual(t, yapi.GetTag(yval1), yapi.YvalScalar)
	AssertEqual(t, yapi.GetTag(yval2), yapi.YvalScalar)
	AssertEqual(t, yapi.GetTag(yval3), yapi.YvalScalar)

	var tau1, tau2, tau3 yapi.TypeT

	AssertEqual(t, yapi.ValGetScalar(*modelp, &yval1, &val1, &tau1), 0)
	AssertEqual(t, yapi.ValGetScalar(*modelp, &yval2, &val2, &tau2), 0)
	AssertEqual(t, yapi.ValGetScalar(*modelp, &yval3, &val3, &tau3), 0)

	AssertEqual(t, val1, 9)
	AssertEqual(t, val2, 8)
	AssertEqual(t, val3, 8)

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func TestModelFromMap(t *testing.T) {

	cfg, ctx, params := setup()

	bvT := yapi.BvType(8)
	intT := yapi.IntType()
	realT := yapi.RealType()

	i1 := defineConst("i1", intT)
	r1 := defineConst("r1", realT)
	bv1 := defineConst("bv1", bvT)

	iconst1 := yapi.Int32(42)
	rconst1 := yapi.Rational32(13, 131)
	bvconst1 := yapi.BvconstInt32(8, 134)

	vars := []yapi.TermT{i1, r1, bv1}
	vals := []yapi.TermT{iconst1, rconst1, bvconst1}

	modelp := yapi.ModelFromMap(vars, vals)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	modelstr := yapi.ModelToString(*modelp, 80, 100, 0)

	AssertEqual(t, modelstr, "(= i1 42)\n(= r1 13/131)\n(= bv1 0b10000110)")

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func TestImplicant(t *testing.T) {

	cfg, ctx, params := setup()

	intT := yapi.IntType()

	i1 := defineConst("i1", intT)

	parseAssert("(and (> i1 2) (< i1 8) (/= i1 4))", ctx)

	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	fmla0 := yapi.ParseTerm("(>= i1 3)")

	terms := yapi.ImplicantForFormula(*modelp, fmla0)

	AssertEqual(t, len(terms), 1)

	mdlstr := yapi.ModelToString(*modelp, 80, 100, 0)
	AssertEqual(t, mdlstr, "(= i1 7)")

	implstr := yapi.TermToString(terms[0], 200, 10, 0)
	AssertEqual(t, implstr, "(>= (+ -3 i1) 0)")

	fmla1 := yapi.ParseTerm("(<= i1 9)")
	fmlas := []yapi.TermT{fmla0, fmla1}

	terms = yapi.ImplicantForFormulas(*modelp, fmlas)
	AssertEqual(t, len(terms), 2)

	implstr2 := yapi.TermToString(terms[0], 200, 10, 0)
	AssertEqual(t, implstr2, "(>= (+ -3 i1) 0)")
	implstr3 := yapi.TermToString(terms[1], 200, 10, 0)
	AssertEqual(t, implstr3, "(>= (+ 9 (* -1 i1)) 0)")

	fmlas = yapi.GeneralizeModelArray(*modelp, fmlas, []yapi.TermT{i1}, 0)
	AssertEqual(t, len(fmlas), 2)
	tstr0 := yapi.TermToString(fmlas[0], 200, 10, 0)
	AssertEqual(t, tstr0, "true")
	tstr1 := yapi.TermToString(fmlas[1], 200, 10, 0)
	AssertEqual(t, tstr1, "true")

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}

func TestYvalNumericModels(t *testing.T) {

	cfg, ctx, params := setup()

	intT := yapi.IntType()

	i1 := defineConst("i1", intT)
	i2 := defineConst("i2", intT)

	parseAssert("(> i1 3)", ctx)
	parseAssert("(< i2 i1)", ctx)

	stat := yapi.CheckContext(ctx, params)
	AssertEqual(t, stat, yapi.StatusSat, "stat == yapi.StatusSat")
	modelp := yapi.GetModel(ctx, 1)
	AssertNotEqual(t, modelp, nil, "modelp != nil")

	var y1, y2 yapi.YvalT
	errcode := yapi.GetValue(*modelp, i1, &y1)
	AssertEqual(t, errcode, 0)
	yerror := yapi.YicesError()
	AssertEqual(t, yerror, nil) //aye curumba: hay magic in testlib.go to avoid: (*yapi.YicesErrorT)(nil)

	errcode = yapi.GetValue(*modelp, i2, &y2)
	AssertEqual(t, errcode, 0)
	yerror = yapi.YicesError()
	AssertEqual(t, yerror, nil) //aye curumba. ibid

	AssertEqual(t, yapi.ValIsInt32(*modelp, &y1), 1)
	AssertEqual(t, yapi.ValIsInt64(*modelp, &y1), 1)
	AssertEqual(t, yapi.ValIsRational32(*modelp, &y1), 1)
	AssertEqual(t, yapi.ValIsRational64(*modelp, &y1), 1)
	AssertEqual(t, yapi.ValIsInteger(*modelp, &y1), 1)

	AssertEqual(t, yapi.ValBitsize(*modelp, &y1), 0)
	AssertEqual(t, yapi.ValTupleArity(*modelp, &y1), 0)
	AssertEqual(t, yapi.ValMappingArity(*modelp, &y1), 0)
	AssertEqual(t, yapi.ValFunctionArity(*modelp, &y1), 0)

	var b1 int32

	errcode = yapi.ValGetBool(*modelp, &y1, &b1)
	AssertEqual(t, errcode, -1)
	yerror = yapi.YicesError()
	AssertNotEqual(t, yerror, nil) //aye curumba. ibid
	AssertEqual(t, yerror.ErrorString, "invalid operation on yval")
	AssertEqual(t, yerror.Code, yapi.ErrorYvalInvalidOp)

	var ival32 int32
	errcode = yapi.ValGetInt32(*modelp, &y1, &ival32)
	AssertEqual(t, errcode, 0)
	yerror = yapi.YicesError()
	AssertEqual(t, yerror, nil) //aye curumba. ibid
	AssertEqual(t, ival32, 4)

	var ival64 int64
	errcode = yapi.ValGetInt64(*modelp, &y1, &ival64)
	AssertEqual(t, errcode, 0)
	yerror = yapi.YicesError()
	AssertEqual(t, yerror, nil) //aye curumba. ibid
	AssertEqual(t, ival64, 4)

	var num32 int32
	var den32 uint32
	errcode = yapi.ValGetRational32(*modelp, &y1, &num32, &den32)
	AssertEqual(t, errcode, 0)
	AssertEqual(t, num32, 4)
	AssertEqual(t, den32, 1)

	var num64 int64
	var den64 uint64
	errcode = yapi.ValGetRational64(*modelp, &y1, &num64, &den64)
	AssertEqual(t, errcode, 0)
	AssertEqual(t, num64, 4)
	AssertEqual(t, den64, 1)

	var dval float64
	errcode = yapi.ValGetDouble(*modelp, &y1, &dval)
	AssertEqual(t, errcode, 0)
	AssertEqual(t, dval, 4.0)

	var mpz yapi.MpzT
	yapi.InitMpz(&mpz)
	errcode = yapi.ValGetMpz(*modelp, &y1, &mpz)
	AssertEqual(t, errcode, 0)
	ytz := yapi.Mpz(&mpz)
	AssertEqual(t, yapi.TermToString(ytz, 200, 10, 0), "4")
	yapi.CloseMpz(&mpz)

	var mpq yapi.MpqT
	yapi.InitMpq(&mpq)
	errcode = yapi.ValGetMpq(*modelp, &y1, &mpq)
	AssertEqual(t, errcode, 0)
	ytq := yapi.Mpq(&mpq)
	AssertEqual(t, yapi.TermToString(ytq, 200, 10, 0), "4")
	yapi.CloseMpq(&mpq)

	yapi.CloseModel(modelp)

	cleanup(&cfg, &ctx, &params)

}
