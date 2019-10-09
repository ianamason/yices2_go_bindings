package yices

// #cgo CFLAGS: -g -fPIC
// #cgo LDFLAGS:  -lyices -lgmp
// #include <yices.h>
import "C"

func version() string {
	return "just little steps for now"
}
