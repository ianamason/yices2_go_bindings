package main

import (
	"fmt"
	"github.com/ianamason/yices2_go_bindings/yices2"
)

func main() {

	fmt.Printf("Yices version no:     %s\n", yices2.Version())
	fmt.Printf("Yices build arch:     %s\n", yices2.Build_arch())
	fmt.Printf("Yices build mode:     %s\n", yices2.Build_mode())
	fmt.Printf("Yices build date:     %s\n", yices2.Build_date())
	fmt.Printf("Yices has mcsat:      %d\n", yices2.Has_mcsat())
	fmt.Printf("Yices is thread safe: %d\n", yices2.Is_thread_safe())

}
