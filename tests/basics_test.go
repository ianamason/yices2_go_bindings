package tests

import (
	"fmt"
	"github.com/ianamason/yices2_go_bindings/yices2"
	"os"
	"testing"
)

func TestBasis0(t *testing.T) {
	yices2.Init()

	bvt := yices2.Bool_type()
	ivt := yices2.Int_type()
	rvt := yices2.Real_type()

	AssertNotEqual(t, bvt, ivt, "bvt != ivt")
	AssertNotEqual(t, bvt, rvt, "bvt != rvt")
	AssertNotEqual(t, ivt, rvt, "ivt != ivt")

	fmt.Printf("Bool_type(): %v %v\n", bvt, yices2.Type_is_bool(bvt))
	fmt.Printf("Int_type(): %v %v\n", ivt, yices2.Type_is_int(ivt))
	fmt.Printf("Real_type(): %v %v\n", rvt, yices2.Type_is_real(rvt))

	typs := []yices2.Type_t { bvt, ivt, rvt }

	tupt := yices2.Tuple_type(typs)

	fmt.Printf("Tuple type(): %v %v\n", tupt, yices2.Type_is_tuple(tupt))

	yices2.Pp_type(os.Stdout, tupt, 80, 80, 10)

	children := yices2.Type_children(tupt)

	fmt.Printf("children: %v\n", children)


	AssertEqual(t, typs, children, "typs == children")



	fmt.Println("Exiting...")
	yices2.Exit()

}
