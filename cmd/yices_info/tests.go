package main

import (
	"fmt"
	"github.com/ianamason/yices2_go_bindings/yices2"
	"os"
)


func Test0() {
	yices2.Init()

	bvt := yices2.Bool_type()
	ivt := yices2.Int_type()
	rvt := yices2.Real_type()
	
	fmt.Printf("Bool_type(): %v\n", bvt)
	fmt.Printf("Int_type(): %v\n", ivt)
	fmt.Printf("Real_type(): %v\n", rvt)

	typs := []yices2.Type_t { bvt, ivt, rvt }

	tupt := yices2.Tuple_type(typs)

	fmt.Printf("Tuple type(): %v\n", tupt)
	
	yices2.Pp_type(os.Stdout, tupt, 80, 80, 10)
	
	yices2.Exit()
	
}
