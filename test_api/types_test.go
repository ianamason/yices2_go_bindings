package tests

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"testing"
)

func TestTypes0(t *testing.T) {
	yapi.Init()

	boolT := yapi.BoolType()
	intT := yapi.IntType()
	AssertNotEqual(t, boolT, intT)
	realT := yapi.RealType()
	AssertNotEqual(t, realT, boolT)
	AssertNotEqual(t, realT, intT)
	bvT := yapi.BvType(8)
	scalT := yapi.NewScalarType(12)
	unintT := yapi.NewUninterpretedType()
	tup1T := yapi.TupleType1(boolT)
	tup2T := yapi.TupleType2(intT, realT)
	tup3T := yapi.TupleType3(bvT, scalT, unintT)
	tup4T := yapi.TupleType([]yapi.TypeT{boolT, tup1T, tup2T, tup3T})
	fun1T := yapi.FunctionType1(intT, boolT)
	fun2T := yapi.FunctionType2(realT, bvT, scalT)
	fun3T := yapi.FunctionType3(tup1T, tup2T, tup3T, fun1T)
	fun4T := yapi.FunctionType([]yapi.TypeT{boolT, tup1T, tup2T, tup3T}, fun3T)

	AssertTrue(t, yapi.TypeIsBool(boolT))
	AssertFalse(t, yapi.TypeIsBool(intT))
	AssertTrue(t, yapi.TypeIsInt(intT))
	AssertTrue(t, yapi.TypeIsReal(realT))
	AssertTrue(t, yapi.TypeIsArithmetic(realT))
	AssertTrue(t, yapi.TypeIsBitvector(bvT))
	AssertTrue(t, yapi.TypeIsTuple(tup1T))
	AssertTrue(t, yapi.TypeIsFunction(fun2T))
	AssertTrue(t, yapi.TypeIsFunction(fun3T))
	AssertTrue(t, yapi.TypeIsFunction(fun4T))
	AssertTrue(t, yapi.TypeIsScalar(scalT))
	AssertTrue(t, yapi.TypeIsUninterpreted(unintT))
	AssertTrue(t, yapi.TestSubtype(intT, realT))
	AssertFalse(t, yapi.TestSubtype(realT, intT))
	AssertEqual(t, yapi.BvtypeSize(bvT), uint32(8))        //yuk
	AssertEqual(t, yapi.ScalarTypeCard(scalT), uint32(12)) //yuk
	AssertEqual(t, yapi.TypeNumChildren(tup3T), int32(3))  //yuk
	AssertEqual(t, yapi.TypeChild(tup3T, 1), scalT)
	tc := yapi.TypeChildren(tup4T)
	AssertEqual(t, len(tc), 4)
	AssertEqual(t, tc[0], boolT)
	AssertEqual(t, tc[1], tup1T)
	AssertEqual(t, tc[2], tup2T)
	AssertEqual(t, tc[3], tup3T)

	yapi.Exit()

}
