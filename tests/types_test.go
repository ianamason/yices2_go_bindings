package tests

import (
	"github.com/ianamason/yices2_go_bindings/yices2"
	"testing"
)

func TestTypes0(t *testing.T) {
	yices2.Init()

	bool_t := yices2.Bool_type()
	int_t := yices2.Int_type()
	AssertNotEqual(t, bool_t, int_t)
	real_t := yices2.Real_type()
	AssertNotEqual(t, real_t, bool_t)
	AssertNotEqual(t, real_t, int_t)
	bv_t := yices2.Bv_type(8)
	scal_t := yices2.New_scalar_type(12)
	unint_t := yices2.New_uninterpreted_type()
	tup1_t := yices2.Tuple_type1(bool_t)
	tup2_t := yices2.Tuple_type2(int_t, real_t)
	tup3_t := yices2.Tuple_type3(bv_t, scal_t, unint_t)
	tup4_t := yices2.Tuple_type([]yices2.Type_t{bool_t, tup1_t, tup2_t, tup3_t})
	fun1_t := yices2.Function_type1(int_t, bool_t)
	fun2_t := yices2.Function_type2(real_t, bv_t, scal_t)
	fun3_t := yices2.Function_type3(tup1_t, tup2_t, tup3_t, fun1_t)
	fun4_t := yices2.Function_type([]yices2.Type_t{bool_t, tup1_t, tup2_t, tup3_t}, fun3_t)

	AssertTrue(t, yices2.Type_is_bool(bool_t))
	AssertFalse(t, yices2.Type_is_bool(int_t))
	AssertTrue(t, yices2.Type_is_int(int_t))
	AssertTrue(t, yices2.Type_is_real(real_t))
	AssertTrue(t, yices2.Type_is_arithmetic(real_t))
	AssertTrue(t, yices2.Type_is_bitvector(bv_t))
	AssertTrue(t, yices2.Type_is_tuple(tup1_t))
	AssertTrue(t, yices2.Type_is_function(fun2_t))
	AssertTrue(t, yices2.Type_is_function(fun3_t))
	AssertTrue(t, yices2.Type_is_function(fun4_t))
	AssertTrue(t, yices2.Type_is_scalar(scal_t))
	AssertTrue(t, yices2.Type_is_uninterpreted(unint_t))
	AssertTrue(t, yices2.Test_subtype(int_t, real_t))
	AssertFalse(t, yices2.Test_subtype(real_t, int_t))
	AssertEqual(t, yices2.Bvtype_size(bv_t), uint32(8)) //yuk
	AssertEqual(t, yices2.Scalar_type_card(scal_t), uint32(12)) //yuk
	AssertEqual(t, yices2.Type_num_children(tup3_t), int32(3)) //yuk
	AssertEqual(t, yices2.Type_child(tup3_t, 1), scal_t)
	tc := yices2.Type_children(tup4_t)
	AssertEqual(t, len(tc), 4)
	AssertEqual(t, tc[0], bool_t)
	AssertEqual(t, tc[1], tup1_t)
	AssertEqual(t, tc[2], tup2_t)
	AssertEqual(t, tc[3], tup3_t)


	yices2.Exit()

}
