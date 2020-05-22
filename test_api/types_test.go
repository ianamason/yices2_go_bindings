package tests

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"testing"
)

func TestTypes0(t *testing.T) {
	yapi.Init()

	bool_t := yapi.Bool_type()
	int_t := yapi.Int_type()
	AssertNotEqual(t, bool_t, int_t)
	real_t := yapi.Real_type()
	AssertNotEqual(t, real_t, bool_t)
	AssertNotEqual(t, real_t, int_t)
	bv_t := yapi.Bv_type(8)
	scal_t := yapi.New_scalar_type(12)
	unint_t := yapi.New_uninterpreted_type()
	tup1_t := yapi.Tuple_type1(bool_t)
	tup2_t := yapi.Tuple_type2(int_t, real_t)
	tup3_t := yapi.Tuple_type3(bv_t, scal_t, unint_t)
	tup4_t := yapi.Tuple_type([]yapi.TypeT{bool_t, tup1_t, tup2_t, tup3_t})
	fun1_t := yapi.Function_type1(int_t, bool_t)
	fun2_t := yapi.Function_type2(real_t, bv_t, scal_t)
	fun3_t := yapi.Function_type3(tup1_t, tup2_t, tup3_t, fun1_t)
	fun4_t := yapi.Function_type([]yapi.TypeT{bool_t, tup1_t, tup2_t, tup3_t}, fun3_t)

	AssertTrue(t, yapi.Type_is_bool(bool_t))
	AssertFalse(t, yapi.Type_is_bool(int_t))
	AssertTrue(t, yapi.Type_is_int(int_t))
	AssertTrue(t, yapi.Type_is_real(real_t))
	AssertTrue(t, yapi.Type_is_arithmetic(real_t))
	AssertTrue(t, yapi.Type_is_bitvector(bv_t))
	AssertTrue(t, yapi.Type_is_tuple(tup1_t))
	AssertTrue(t, yapi.Type_is_function(fun2_t))
	AssertTrue(t, yapi.Type_is_function(fun3_t))
	AssertTrue(t, yapi.Type_is_function(fun4_t))
	AssertTrue(t, yapi.Type_is_scalar(scal_t))
	AssertTrue(t, yapi.Type_is_uninterpreted(unint_t))
	AssertTrue(t, yapi.Test_subtype(int_t, real_t))
	AssertFalse(t, yapi.Test_subtype(real_t, int_t))
	AssertEqual(t, yapi.Bvtype_size(bv_t), uint32(8))         //yuk
	AssertEqual(t, yapi.Scalar_type_card(scal_t), uint32(12)) //yuk
	AssertEqual(t, yapi.Type_num_children(tup3_t), int32(3))  //yuk
	AssertEqual(t, yapi.Type_child(tup3_t, 1), scal_t)
	tc := yapi.Type_children(tup4_t)
	AssertEqual(t, len(tc), 4)
	AssertEqual(t, tc[0], bool_t)
	AssertEqual(t, tc[1], tup1_t)
	AssertEqual(t, tc[2], tup2_t)
	AssertEqual(t, tc[3], tup3_t)

	yapi.Exit()

}
