package main

import (
	"fmt"
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
)

func main() {

	fmt.Printf("Yices version no:     %s\n", yapi.Version())
	fmt.Printf("Yices build arch:     %s\n", yapi.Build_arch())
	fmt.Printf("Yices build mode:     %s\n", yapi.Build_mode())
	fmt.Printf("Yices build date:     %s\n", yapi.Build_date())
	fmt.Printf("Yices has mcsat:      %d\n", yapi.Has_mcsat())
	fmt.Printf("Yices is thread safe: %d\n", yapi.Is_thread_safe())

}
