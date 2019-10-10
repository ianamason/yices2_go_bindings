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

/********************************
 *  VECTORS OF TERMS AND TYPES  *
 *******************************/

//iam: probably not going to export these 
//__YICES_DLLSPEC__ extern void yices_init_term_vector(term_vector_t *v);
//__YICES_DLLSPEC__ extern void yices_init_type_vector(type_vector_t *v);
//__YICES_DLLSPEC__ extern void yices_delete_term_vector(term_vector_t *v);
//__YICES_DLLSPEC__ extern void yices_delete_type_vector(type_vector_t *v);
//__YICES_DLLSPEC__ extern void yices_reset_term_vector(term_vector_t *v);
//__YICES_DLLSPEC__ extern void yices_reset_type_vector(type_vector_t *v);


/***********************
 *  TYPE CONSTRUCTORS  *
 **********************/

//iam: we use a type definition, does it improve readbility?

type Type_t int32 

func Bool_type() Type_t { return Type_t(C.yices_bool_type()) }

func Int_type() Type_t { return Type_t(C.yices_int_type()) }

func Real_type() Type_t { return Type_t(C.yices_real_type()) }

func Bv_type(size uint32) Type_t  { return Type_t(C.yices_bv_type(C.uint32_t(size))) }

func New_scalar_type(card uint32) Type_t  { return Type_t(C.yices_new_scalar_type(C.uint32_t(card))) }

func New_uninterpreted_type() Type_t { return Type_t(C.yices_new_uninterpreted_type()) }

/*
//iam: icky
func type_slice_to_C_array(s []Type_t) *C.int {
	var carr = make([]C.int, len(s), len(s))
	for i := 0; i < len(s); i++ {
		carr[i] = C.int(s[i])
	}
	return &carr[0]
}
func Tuple_type(tau []Type_t) Type_t {
	tau_len := len(tau)
	return Type_t(C.yices_tuple_type(C.uint32_t(tau_len), type_slice_to_C_array(tau)))
}
*/

func Tuple_type(tau []Type_t) Type_t {
	tau_len := len(tau)
	return Type_t(C.yices_tuple_type(C.uint32_t(tau_len), (*C.int32_t)(&tau[0])))
}

func Tuple_type1(tau1 Type_t) Type_t {
	carr := []C.int32_t{ C.int32_t(tau1) }
	return Type_t(C.yices_tuple_type(C.uint32_t(1), (*C.int32_t)(&carr[0])))
}

func Tuple_type2(tau1 Type_t, tau2 Type_t) Type_t {
	carr := []C.int32_t{ C.int32_t(tau1), C.int32_t(tau2) }
	return Type_t(C.yices_tuple_type(C.uint32_t(2), (*C.int32_t)(&carr[0])))
}

func Tuple_type3(tau1 Type_t, tau2 Type_t, tau3 Type_t) Type_t {
	carr := []C.int32_t{ C.int32_t(tau1), C.int32_t(tau2), C.int32_t(tau3) }
	return Type_t(C.yices_tuple_type(C.uint32_t(3), (*C.int32_t)(&carr[0])))
}


func Function_type(dom []Type_t, rng Type_t) Type_t {
	dom_len := len(dom)
	return Type_t(C.yices_function_type(C.uint32_t(dom_len), (*C.int32_t)(&dom[0]), C.int32_t(rng)))
}

func Function_type1(tau1 Type_t, rng Type_t) Type_t {
	carr := []C.int32_t{ C.int32_t(tau1) }
	return Type_t(C.yices_function_type(C.uint32_t(1), (*C.int32_t)(&carr[0]), C.int32_t(rng)))
}

func Function_type2(tau1 Type_t, tau2 Type_t, rng Type_t) Type_t {
	carr := []C.int32_t{ C.int32_t(tau1), C.int32_t(tau2) }
	return Type_t(C.yices_function_type(C.uint32_t(2), (*C.int32_t)(&carr[0]), C.int32_t(rng)))
}

func Function_type3(tau1 Type_t, tau2 Type_t, tau3 Type_t, rng Type_t) Type_t {
	carr := []C.int32_t{ C.int32_t(tau1), C.int32_t(tau2), C.int32_t(tau3) }
	return Type_t(C.yices_function_type(C.uint32_t(3), (*C.int32_t)(&carr[0]), C.int32_t(rng)))
}

/*************************
 *   TYPE EXPLORATION    *
 ************************/
