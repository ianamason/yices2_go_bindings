package yices2

// #cgo CFLAGS: -g -fPIC
// #cgo LDFLAGS:  -lyices -lgmp
// #include <yices.h>
import "C"

import "os"

/*
 *  See yices.h for comments.
 *
 *  Naming convention:   yices_foo  becomes yices2.Foo  (will probably ditch the 2 at some stage)
 *
 */ 

/*********************
 *  VERSION NUMBERS  *
 ********************/

func Version() string {
	return C.GoString(C.yices_version)
}

func Build_arch() string {
	return C.GoString(C.yices_build_arch)
}

func Build_mode() string {
	return C.GoString(C.yices_build_mode)
}

func Build_date() string {
	return C.GoString(C.yices_build_date)
}

func Has_mcsat() int32 {
	return int32(C.yices_has_mcsat())
}

func Is_thread_safe() int32 {
	return int32(C.yices_is_thread_safe())
}


/***************************************
 *  GLOBAL INITIALIZATION AND CLEANUP  *
 **************************************/

func Init() { C.yices_init() }

func Exit() { C.yices_exit() }

func Reset() { C.yices_reset() }

//__YICES_DLLSPEC__ extern void yices_free_string(char *s);  iam: should be unnecessary? we only ever return go strings


/***************************
 * OUT-OF-MEMORY CALLBACK  *
 **************************/

//__YICES_DLLSPEC__ extern void yices_set_out_of_mem_callback(void (*callback)(void));  iam: defer this


/*********************
 *  ERROR REPORTING  *
 ********************/

//__YICES_DLLSPEC__ extern error_code_t yices_error_code(void);

func Error_code() int32 {  return int32(C.yices_error_code()) }  //iam: FIXME error_code_t should (associated with a) be a go type

//__YICES_DLLSPEC__ extern error_report_t *yices_error_report(void); //iam: FIXME seem to recall this has unions in it?

func Clear_error() { C.yices_clear_error() }

func Print_error(f *os.File) int32 { return int32(C.yices_print_error_fd(C.int(f.Fd()))) }  //iam: FIXME error checking and File without the os.

func Error_string() string {
	return C.GoString(C.yices_error_string())
}
