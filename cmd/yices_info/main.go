package main

import (
	"fmt"
	"github.com/ianamason/yices2_go_bindings/yices2"
)

func main() {

	fmt.Printf("Yices version no: %s\n", yices2.Version())
	fmt.Printf("Yices build arch: %s\n", yices2.Build_arch())
	fmt.Printf("Yices build mode: %s\n", yices2.Build_mode())
	fmt.Printf("Yices build date: %s\n", yices2.Build_date())

	//while we figure out the memory corruption issue...
	yices2.Init()

	var cfg yices2.Config_t

	yices2.Init_config(&cfg)

	var ctx yices2.Context_t

	yices2.Init_context(cfg, &ctx)

	yices2.Close_config(&cfg)

	bv_t := yices2.Bv_type(3)
	bvvar1 := yices2.New_uninterpreted_term(bv_t)
	yices2.Set_term_name(bvvar1, "x")
	bvvar2 := yices2.New_uninterpreted_term(bv_t)
	yices2.Set_term_name(bvvar2, "y")
	bvvar3 := yices2.New_uninterpreted_term(bv_t)
	yices2.Set_term_name(bvvar3, "z")
	fmla1 := yices2.Parse_term("(= x (bv-add y z))")
	fmla2 := yices2.Parse_term("(bv-gt y 0b000)")
	fmla3 := yices2.Parse_term("(bv-gt z 0b000)")
	yices2.Assert_formula(ctx, fmla1)
	//yices2.Assert_formula(ctx, fmla2)
	//yices2.Assert_formula(ctx, fmla3)
	yices2.Assert_formulas(ctx, []yices2.Term_t{fmla1, fmla2, fmla3})
	smt_stat := yices2.Check_context(ctx, nil)

	fmt.Printf("Context status is yices2.STATUS_SAT: %v\n", smt_stat == yices2.STATUS_SAT)

	yices2.Close_context(&ctx)


	yices2.Exit()
}
