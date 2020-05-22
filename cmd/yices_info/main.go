package main

import (
	"fmt"
	"github.com/ianamason/yices2_go_bindings/yices"
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
)

func main() {

	fmt.Printf("Yices version no:           %s\n", yices.Version)
	fmt.Printf("Yices build arch:           %s\n", yices.Build_arch)
	fmt.Printf("Yices build mode:           %s\n", yices.Build_mode)
	fmt.Printf("Yices build date:           %s\n", yices.Build_date)
	fmt.Printf("Yices has mcsat:            %v\n", yices.Has_mcsat)
	fmt.Printf("Yices is thread safe:       %v\n", yices.Is_thread_safe)
	fmt.Printf("Yices has cadical:          %v\n", yapi.Has_delegate("cadical"))
	fmt.Printf("Yices has cryptominisat:    %v\n", yapi.Has_delegate("cryptominisat"))
	fmt.Printf("Yices has y2sat:            %v\n", yapi.Has_delegate("y2sat"))
}
