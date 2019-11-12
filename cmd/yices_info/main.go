package main

import (
	"fmt"
	"github.com/ianamason/yices2_go_bindings/yices"
)

func main() {

	fmt.Printf("Yices version no:     %s\n", yices.Version)
	fmt.Printf("Yices build arch:     %s\n", yices.Build_arch)
	fmt.Printf("Yices build mode:     %s\n", yices.Build_mode)
	fmt.Printf("Yices build date:     %s\n", yices.Build_date)
	fmt.Printf("Yices has mcsat:      %d\n", yices.Has_mcsat)
	fmt.Printf("Yices is thread safe: %d\n", yices.Is_thread_safe)
}
