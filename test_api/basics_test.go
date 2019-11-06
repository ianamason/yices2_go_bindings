package tests

import (
	"fmt"
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"os"
	"testing"
)

func TestBasis0(t *testing.T) {
	yapi.Init()

	bvt := yapi.Bool_type()
	ivt := yapi.Int_type()
	rvt := yapi.Real_type()

	AssertNotEqual(t, bvt, ivt, "bvt != ivt")
	AssertNotEqual(t, bvt, rvt, "bvt != rvt")
	AssertNotEqual(t, ivt, rvt, "ivt != ivt")

	fmt.Printf("Bool_type(): %v %v\n", bvt, yapi.Type_is_bool(bvt))
	fmt.Printf("Int_type(): %v %v\n", ivt, yapi.Type_is_int(ivt))
	fmt.Printf("Real_type(): %v %v\n", rvt, yapi.Type_is_real(rvt))

	typs := []yapi.Type_t{bvt, ivt, rvt}

	tupt := yapi.Tuple_type(typs)

	fmt.Printf("Tuple type(): %v %v\n", tupt, yapi.Type_is_tuple(tupt))

	yapi.Pp_type(os.Stdout, tupt, 80, 80, 10)

	children := yapi.Type_children(tupt)

	fmt.Printf("children: %v\n", children)

	AssertEqual(t, typs, children, "typs == children")

	fmt.Println("Exiting...")
	yapi.Exit()

}
