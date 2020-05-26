package tests

import (
	"fmt"
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"os"
	"testing"
)

func TestBasis0(t *testing.T) {
	yapi.Init()

	bvt := yapi.BoolType()
	ivt := yapi.IntType()
	rvt := yapi.RealType()

	AssertNotEqual(t, bvt, ivt, "bvt != ivt")
	AssertNotEqual(t, bvt, rvt, "bvt != rvt")
	AssertNotEqual(t, ivt, rvt, "ivt != ivt")

	fmt.Printf("BoolType(): %v %v\n", bvt, yapi.TypeIsBool(bvt))
	fmt.Printf("IntType(): %v %v\n", ivt, yapi.TypeIsInt(ivt))
	fmt.Printf("RealType(): %v %v\n", rvt, yapi.TypeIsReal(rvt))

	typs := []yapi.TypeT{bvt, ivt, rvt}

	tupt := yapi.TupleType(typs)

	fmt.Printf("Tuple type(): %v %v\n", tupt, yapi.TypeIsTuple(tupt))

	yapi.PpType(os.Stdout, tupt, 80, 80, 10)

	children := yapi.TypeChildren(tupt)

	fmt.Printf("children: %v\n", children)

	AssertEqual(t, typs, children, "typs == children")

	fmt.Println("Exiting...")
	yapi.Exit()

}
