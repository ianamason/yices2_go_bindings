package yices2

/*
#cgo CFLAGS: -g -fPIC
#cgo LDFLAGS:  -lyices -lgmp
#include <stdlib.h>
#include <yices.h>
type_t yices_type_vector_get(type_vector_t* vec, uint32_t elem){ return vec->data[elem]; }
*/
import "C"

import "os"
import "unsafe"

/*
 *  See yices.h for comments.
 *
 *  Naming convention:   yices_foo  becomes yices2.Foo  (will probably ditch the 2 at some stage)
 *
 *  bd: - free the strings returned by yices.
 *      - check that some int32 retvals could be bool
 *
 * iam: - free the result of C.CString using C.free (maybe wrapped with unsafe.Pointer?)
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

func Init() {
	C.yices_init()
}

func Exit() {
	C.yices_exit()
}

func Reset() {
	C.yices_reset()
}

//__YICES_DLLSPEC__ extern void yices_free_string(char *s);  iam: should be unnecessary? we only ever return go strings

/***************************
 * OUT-OF-MEMORY CALLBACK  *
 **************************/

//__YICES_DLLSPEC__ extern void yices_set_out_of_mem_callback(void (*callback)(void));  iam: defer this

/*********************
 *  ERROR REPORTING  *
 ********************/

//__YICES_DLLSPEC__ extern error_code_t yices_error_code(void);

func Error_code() int32 {
	return int32(C.yices_error_code())
} //iam: FIXME error_code_t should (associated with a) be a go type

//__YICES_DLLSPEC__ extern error_report_t *yices_error_report(void); //iam: FIXME seem to recall this has unions in it?

func Clear_error() {
	C.yices_clear_error()
}

func Print_error(f *os.File) int32 {
	return int32(C.yices_print_error_fd(C.int(f.Fd())))
} //iam: FIXME error checking and File without the os.

func Error_string() string {
	cs := C.yices_error_string()
	defer C.yices_free_string(cs)
	return C.GoString(cs)
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

const NULL_TYPE Type_t = -1


func Bool_type() Type_t {
	return Type_t(C.yices_bool_type())
}

func Int_type() Type_t {
	return Type_t(C.yices_int_type())
}

func Real_type() Type_t {
	return Type_t(C.yices_real_type())
}

func Bv_type(size uint32) Type_t {
	return Type_t(C.yices_bv_type(C.uint32_t(size)))
}

func New_scalar_type(card uint32) Type_t {
	return Type_t(C.yices_new_scalar_type(C.uint32_t(card)))
}

func New_uninterpreted_type() Type_t {
	return Type_t(C.yices_new_uninterpreted_type())
}

func Tuple_type(tau []Type_t) Type_t {
	tau_len := len(tau)
	return Type_t(C.yices_tuple_type(C.uint32_t(tau_len), (*C.type_t)(&tau[0])))
}

func Tuple_type1(tau1 Type_t) Type_t {
	carr := []C.type_t{C.type_t(tau1)}
	return Type_t(C.yices_tuple_type(C.uint32_t(1), (*C.type_t)(&carr[0])))
}

func Tuple_type2(tau1 Type_t, tau2 Type_t) Type_t {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2)}
	return Type_t(C.yices_tuple_type(C.uint32_t(2), (*C.type_t)(&carr[0])))
}

func Tuple_type3(tau1 Type_t, tau2 Type_t, tau3 Type_t) Type_t {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2), C.type_t(tau3)}
	return Type_t(C.yices_tuple_type(C.uint32_t(3), (*C.type_t)(&carr[0])))
}

func Function_type(dom []Type_t, rng Type_t) Type_t {
	dom_len := len(dom)
	return Type_t(C.yices_function_type(C.uint32_t(dom_len), (*C.type_t)(&dom[0]), C.type_t(rng)))
}

func Function_type1(tau1 Type_t, rng Type_t) Type_t {
	carr := []C.type_t{C.type_t(tau1)}
	return Type_t(C.yices_function_type(C.uint32_t(1), (*C.type_t)(&carr[0]), C.type_t(rng)))
}

func Function_type2(tau1 Type_t, tau2 Type_t, rng Type_t) Type_t {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2)}
	return Type_t(C.yices_function_type(C.uint32_t(2), (*C.type_t)(&carr[0]), C.type_t(rng)))
}

func Function_type3(tau1 Type_t, tau2 Type_t, tau3 Type_t, rng Type_t) Type_t {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2), C.type_t(tau3)}
	return Type_t(C.yices_function_type(C.uint32_t(3), (*C.type_t)(&carr[0]), C.type_t(rng)))
}

/*************************
 *   TYPE EXPLORATION    *
 ************************/

func Type_is_bool(tau Type_t) int32 {
	return int32(C.yices_type_is_bool(C.type_t(tau)))
}

func Type_is_int(tau Type_t) int32 {
	return int32(C.yices_type_is_int(C.type_t(tau)))
}

func Type_is_real(tau Type_t) int32 {
	return int32(C.yices_type_is_real(C.type_t(tau)))
}

func Type_is_arithmetic(tau Type_t) int32 {
	return int32(C.yices_type_is_arithmetic(C.type_t(tau)))
}

func Type_is_bitvector(tau Type_t) int32 {
	return int32(C.yices_type_is_bitvector(C.type_t(tau)))
}

func Type_is_tuple(tau Type_t) int32 {
	return int32(C.yices_type_is_tuple(C.type_t(tau)))
}

func Type_is_function(tau Type_t) int32 {
	return int32(C.yices_type_is_function(C.type_t(tau)))
}

func Type_is_scalar(tau Type_t) int32 {
	return int32(C.yices_type_is_scalar(C.type_t(tau)))
}

func Type_is_uninterpreted(tau Type_t) int32 {
	return int32(C.yices_type_is_uninterpreted(C.type_t(tau)))
}

func Test_subtype(tau Type_t, sigma Type_t) int32 {
	return int32(C.yices_test_subtype(C.type_t(tau), C.type_t(sigma)))
}

func Compatible_types(tau Type_t, sigma Type_t) int32 {
	return int32(C.yices_compatible_types(C.type_t(tau), C.type_t(sigma)))
}

func Bvtype_size(tau Type_t) uint32 {
	return uint32(C.yices_bvtype_size(C.type_t(tau)))
}

func Scalar_type_card(tau Type_t) uint32 {
	return uint32(C.yices_scalar_type_card(C.type_t(tau)))
}

func Type_num_children(tau Type_t) int32 {
	return int32(C.yices_type_num_children(C.type_t(tau)))
}

func Type_child(tau Type_t, i int32) Type_t {
	return Type_t(C.yices_type_child(C.type_t(tau), C.int32_t(i)))
}

func Type_children(tau Type_t) (children []Type_t) {
	//iam: FIXME is there an easier way?
	var tv [1]C.type_vector_t
	C.yices_init_type_vector(&tv[0])
	ycount := int32(C.yices_type_children(C.type_t(tau), &tv[0]))
	if ycount != -1 {
		count := int(tv[0].size)
		children = make([]Type_t, count, count)
		// defined in the preamble yices_type_vector_get(type_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			children[i] = Type_t(C.yices_type_vector_get(&tv[0], C.uint32_t(i)))
		}
	}
	C.yices_delete_type_vector(&tv[0])
	return
}

/***********************
 *  TERM CONSTRUCTORS  *
 **********************/

type Term_t int32

const NULL_TERM Term_t = -1

func True() Term_t {
	return Term_t(C.yices_true())
}

func False() Term_t {
	return Term_t(C.yices_false())
}

func Constant(tau Type_t, index int32) Term_t {
	return Term_t(C.yices_constant(C.type_t(tau), C.int32_t(index)))
}

func New_uninterpreted_term(tau Type_t) Term_t {
	return Term_t(C.yices_new_uninterpreted_term(C.type_t(tau)))
}

func New_variable(tau Type_t) Term_t {
	return Term_t(C.yices_new_variable(C.type_t(tau)))	
}

func Application(fun Term_t, argv []Term_t) Term_t {
	argc := len(argv)
	return Term_t(C.yices_application(C.term_t(fun), C.uint32_t(argc), (*C.term_t)(&argv[0])))
}

func Application1(fun Term_t, arg1 Term_t) Term_t {
	argv := []C.term_t{ C.term_t(arg1) }
	return Term_t(C.yices_application(C.term_t(fun), C.uint32_t(1), (*C.term_t)(&argv[0])))
}

func Application2(fun Term_t, arg1 Term_t, arg2 Term_t) Term_t {
	argv := []C.term_t{ C.term_t(arg1), C.term_t(arg2) }
	return Term_t(C.yices_application(C.term_t(fun), C.uint32_t(2), (*C.term_t)(&argv[0])))
}

func Application3(fun Term_t, arg1 Term_t, arg2 Term_t, arg3 Term_t) Term_t {
	argv := []C.term_t{ C.term_t(arg1), C.term_t(arg2), C.term_t(arg3) }
	return Term_t(C.yices_application(C.term_t(fun), C.uint32_t(3), (*C.term_t)(&argv[0])))
}

func Ite(cond Term_t, then_term Term_t, else_term Term_t)  Term_t {
	return Term_t(C.yices_ite(C.term_t(cond), C.term_t(then_term), C.term_t(else_term)))
}

func Eq(lhs Term_t, rhs Term_t) Term_t {
	return Term_t(C.yices_eq(C.term_t(lhs), (C.term_t(rhs))))
}

func Neq(lhs Term_t, rhs Term_t) Term_t {
	return Term_t(C.yices_neq(C.term_t(lhs), (C.term_t(rhs))))
}

func Not(arg Term_t) Term_t {
	return Term_t(C.yices_not(C.term_t(arg)))
}

func Or(disjuncts []Term_t) Term_t {
	count := C.uint32_t(len(disjuncts))
	return Term_t(C.yices_or(count, (*C.term_t)(&disjuncts[0])))
}

func And(conjuncts []Term_t) Term_t {
	count := C.uint32_t(len(conjuncts))
	return Term_t(C.yices_and(count, (*C.term_t)(&conjuncts[0])))
}

func Xor(xorjuncts []Term_t) Term_t {
	count := C.uint32_t(len(xorjuncts))
	return Term_t(C.yices_xor(count, (*C.term_t)(&xorjuncts[0])))
}

func Or2(arg1 Term_t, arg2 Term_t) Term_t {
	return Term_t(C.yices_or2(C.term_t(arg1), C.term_t(arg2)))
}

func And2(arg1 Term_t, arg2 Term_t) Term_t {
	return Term_t(C.yices_and2(C.term_t(arg1), C.term_t(arg2)))
}

func Xor2(arg1 Term_t, arg2 Term_t) Term_t {
	return Term_t(C.yices_xor2(C.term_t(arg1), C.term_t(arg2)))
}

func Or3(arg1 Term_t, arg2 Term_t, arg3 Term_t) Term_t {
	return Term_t(C.yices_or3(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

func And3(arg1 Term_t, arg2 Term_t, arg3 Term_t) Term_t {
	return Term_t(C.yices_and3(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

func Xor3(arg1 Term_t, arg2 Term_t, arg3 Term_t) Term_t {
	return Term_t(C.yices_xor3(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

func Iff(lhs Term_t, rhs Term_t) Term_t {
	return Term_t(C.yices_iff(C.term_t(lhs), (C.term_t(rhs))))
}

func Implies(lhs Term_t, rhs Term_t) Term_t {
	return Term_t(C.yices_implies(C.term_t(lhs), (C.term_t(rhs))))
}


func Tuple(argv []Term_t) Term_t {
	count := C.uint32_t(len(argv))
	return Term_t(C.yices_tuple(count, (*C.term_t)(&argv[0])))
}


func Pair(arg1 Term_t, arg2 Term_t) Term_t {
	return Term_t(C.yices_pair(C.term_t(arg1), C.term_t(arg2)))
}

func Triple(arg1 Term_t, arg2 Term_t, arg3 Term_t) Term_t {
	return Term_t(C.yices_triple(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

func Select(index uint32, tuple Term_t) Term_t {
	return Term_t(C.yices_select(C.uint32_t(index), C.term_t(tuple)))
}

func Tuple_update(tuple Term_t,  index uint32, value Term_t) Term_t {
	return Term_t(C.yices_tuple_update(C.term_t(tuple), C.uint32_t(index), C.term_t(value)))
}

func Update(fun Term_t, argv []Term_t, value Term_t) Term_t {
	count := C.uint32_t(len(argv))
	return  Term_t(C.yices_update(C.term_t(fun), count, (*C.term_t)(&argv[0]), C.term_t(value)))
}

func Update1(fun Term_t, arg1 Term_t, value Term_t) Term_t {
	return Term_t(C.yices_update1(C.term_t(fun), C.term_t(arg1), C.term_t(value)))
}

func Update2(fun Term_t, arg1 Term_t, arg2 Term_t, value Term_t) Term_t {
	return Term_t(C.yices_update2(C.term_t(fun), C.term_t(arg1), C.term_t(arg2), C.term_t(value)))
}

func Update3(fun Term_t, arg1 Term_t, arg2 Term_t, arg3 Term_t, value Term_t) Term_t {
	return Term_t(C.yices_update3(C.term_t(fun), C.term_t(arg1), C.term_t(arg2), C.term_t(arg3), C.term_t(value)))
}


func Distinct(argv []Term_t) Term_t {
	n := C.uint32_t(len(argv))
	return Term_t(C.yices_distinct(n, (*C.term_t)(&argv[0])))
}

func Forall(vars []Term_t, body Term_t) Term_t {
	n := C.uint32_t(len(vars))
	return Term_t(C.yices_forall(n, (*C.term_t)(&vars[0]), C.term_t(body)))
}

func Exists(vars []Term_t, body Term_t) Term_t {
	n := C.uint32_t(len(vars))
	return Term_t(C.yices_exists(n, (*C.term_t)(&vars[0]), C.term_t(body)))
}

func Lambda(vars []Term_t, body Term_t) Term_t {
	n := C.uint32_t(len(vars))
	return Term_t(C.yices_lambda(n, (*C.term_t)(&vars[0]), C.term_t(body)))
}

/**********************************
 *  ARITHMETIC TERM CONSTRUCTORS  *
 *********************************/

func Zero() Term_t {
	return Term_t(C.yices_zero())
}

func Int32(val int32) Term_t {
	return Term_t(C.yices_int32(C.int32_t(val)))
}

func Int64(val int64) Term_t {
	return Term_t(C.yices_int64(C.int64_t(val)))
}

func Rational32(num int32, den uint32) Term_t {
	return Term_t(C.yices_rational32(C.int32_t(num), C.uint32_t(den)))
}

func Rational64(num int64, den uint64) Term_t {
	return Term_t(C.yices_rational64(C.int64_t(num), C.uint64_t(den)))
}

/* iam: FIXME in the too hard basket for now.
https://github.com/golang/go/blob/master/misc/cgo/gmp/gmp.go
#ifdef __GMP_H__
__YICES_DLLSPEC__ extern term_t yices_mpz(const mpz_t z);
__YICES_DLLSPEC__ extern term_t yices_mpq(const mpq_t q);
#endif
*/

func Parse_rational(s string) Term_t {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return Term_t(C.yices_parse_rational(cs))
}

func Parse_float(s string) Term_t {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return Term_t(C.yices_parse_float(cs))
}

/*
 * ARITHMETIC OPERATIONS
 */

func Add(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_add(C.term_t(t1), C.term_t(t2)))
}

func Sub(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_sub(C.term_t(t1), C.term_t(t2)))
}

func Neg(t1 Term_t) Term_t {
	return Term_t(C.yices_neg(C.term_t(t1)))
}

func Mul(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_mul(C.term_t(t1), C.term_t(t2)))
}

func Square(t1 Term_t) Term_t {
	return Term_t(C.yices_square(C.term_t(t1)))
}

func Power(t1 Term_t, d uint32) Term_t {
	return Term_t(C.yices_power(C.term_t(t1), C.uint32_t(d)))
}

func Sum(argv []Term_t) Term_t {
	count := C.uint32_t(len(argv))
	return Term_t(C.yices_sum(count, (*C.term_t)(&argv[0])))
}

func Product(argv []Term_t) Term_t {
	count := C.uint32_t(len(argv))
	return Term_t(C.yices_product(count, (*C.term_t)(&argv[0])))
}

func Division(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_division(C.term_t(t1), C.term_t(t2)))
}

func Idiv(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_idiv(C.term_t(t1), C.term_t(t2)))
}

func Imod(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_imod(C.term_t(t1), C.term_t(t2)))
}

func Divides_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_divides_atom(C.term_t(t1), C.term_t(t2)))
}

func Is_int_atom(t Term_t) Term_t {
	return Term_t(C.yices_is_int_atom(C.term_t(t)))
}

func Abs(t1 Term_t) Term_t {
	return Term_t(C.yices_abs(C.term_t(t1)))
}

func Floor(t1 Term_t) Term_t {
	return Term_t(C.yices_floor(C.term_t(t1)))
}

func Ceil(t1 Term_t) Term_t {
	return Term_t(C.yices_ceil(C.term_t(t1)))
}

/*
 * POLYNOMIALS
 */

func Poly_int32(a []int32, t []Term_t) Term_t {
	count := C.uint32_t(len(a))
	return Term_t(C.yices_poly_int32(count, (*C.int32_t)(&a[0]), (*C.term_t)(&t[0])))
}

func Poly_int64(a []int64, t []Term_t) Term_t {
	count := C.uint32_t(len(a))
	return Term_t(C.yices_poly_int64(count, (*C.int64_t)(&a[0]), (*C.term_t)(&t[0])))
}

func Poly_rational32(num []int32, den []uint32, t []Term_t) Term_t {
	count := C.uint32_t(len(num))
	return Term_t(C.yices_poly_rational32(count, (*C.int32_t)(&num[0]), (*C.uint32_t)(&den[0]), (*C.term_t)(&t[0])))
}

func Poly_rational64(num []int64, den []uint64, t []Term_t) Term_t {
	count := C.uint32_t(len(num))
	return Term_t(C.yices_poly_rational64(count, (*C.int64_t)(&num[0]), (*C.uint64_t)(&den[0]), (*C.term_t)(&t[0])))
}


/* iam: FIXME in the too hard basket for now.
https://github.com/golang/go/blob/master/misc/cgo/gmp/gmp.go
#ifdef __GMP_H__
__YICES_DLLSPEC__ extern term_t yices_poly_mpz(uint32_t n, const mpz_t z[], const term_t t[]);
__YICES_DLLSPEC__ extern term_t yices_poly_mpq(uint32_t n, const mpq_t q[], const term_t t[]);
#endif
*/

	
/*
 * ARITHMETIC ATOMS
 */

func Arith_eq_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_arith_eq_atom(C.term_t(t1), C.term_t(t2)))
}

func Arith_neq_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_arith_neq_atom(C.term_t(t1), C.term_t(t2)))
}

func Arith_geq_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_arith_geq_atom(C.term_t(t1), C.term_t(t2)))
}

func Arith_leq_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_arith_leq_atom(C.term_t(t1), C.term_t(t2)))
}

func Arith_gt_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_arith_gt_atom(C.term_t(t1), C.term_t(t2)))
}

func Arith_lt_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_arith_lt_atom(C.term_t(t1), C.term_t(t2)))
}

func Arith_eq0_atom(t Term_t) Term_t {
	return Term_t(C.yices_arith_eq0_atom(C.term_t(t)))
}

func Arith_neq0_atom(t Term_t) Term_t {
	return Term_t(C.yices_arith_neq0_atom(C.term_t(t)))
}

func Arith_geq0_atom(t Term_t) Term_t {
	return Term_t(C.yices_arith_geq0_atom(C.term_t(t)))
}

func Arith_leq0_atom(t Term_t) Term_t {
	return Term_t(C.yices_arith_leq0_atom(C.term_t(t)))
}

func Arith_gt0_atom(t Term_t) Term_t {
	return Term_t(C.yices_arith_gt0_atom(C.term_t(t)))
}

func Arith_lt0_atom(t Term_t) Term_t {
	return Term_t(C.yices_arith_lt0_atom(C.term_t(t)))
}

/*********************************
 *  BITVECTOR TERM CONSTRUCTORS  *
 ********************************/


func Bvconst_uint32(bits uint32,  x uint32) Term_t {
	return Term_t(C.yices_bvconst_uint32(C.uint32_t(bits), C.uint32_t(x)))
}

func Bvconst_uint64(bits uint32,  x uint64) Term_t {
	return Term_t(C.yices_bvconst_uint64(C.uint32_t(bits), C.uint64_t(x)))
}

func Bvconst_int32(bits uint32,  x int32) Term_t {
	return Term_t(C.yices_bvconst_int32(C.uint32_t(bits), C.int32_t(x)))
}

func Bvconst_int64(bits uint32,  x int64) Term_t {
	return Term_t(C.yices_bvconst_int64(C.uint32_t(bits), C.int64_t(x)))
}

/* iam: FIXME
#ifdef __GMP_H__
__YICES_DLLSPEC__ extern term_t yices_bvconst_mpz(uint32_t n, const mpz_t x);
#endif
*/

func Bvconst_zero(bits uint32) Term_t {
	return Term_t(C.yices_bvconst_zero(C.uint32_t(bits)))
}

func Bvconst_one(bits uint32) Term_t {
	return Term_t(C.yices_bvconst_one(C.uint32_t(bits)))
}

func Bvconst_minus_one(bits uint32) Term_t {
	return Term_t(C.yices_bvconst_minus_one(C.uint32_t(bits)))
}

//iam: FIXME check that bits is restricted to len(a)
func Bvconstr_from_array(a []int32) Term_t {
	bits := C.uint32_t(len(a))
	return Term_t(C.yices_bvconst_from_array(bits, (*C.int32_t)(&a[0])))
}


func Parse_bvbin(s string) Term_t {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return Term_t(C.yices_parse_bvbin(cs))
}

func Parse_bvhex(s string) Term_t {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return Term_t(C.yices_parse_bvhex(cs))
}

func Bvadd(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvadd(C.term_t(t1), C.term_t(t2)))
}

func Bsubv(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvsub(C.term_t(t1), C.term_t(t2)))
}

func Bvneg(t Term_t) Term_t {
	return Term_t(C.yices_bvneg(C.term_t(t)))
}

func Bvmul(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvmul(C.term_t(t1), C.term_t(t2)))
}

func Bvsquare(t Term_t) Term_t {
	return Term_t(C.yices_bvsquare(C.term_t(t)))
}

func Bvpower(t1 Term_t, d uint32) Term_t {
	return Term_t(C.yices_bvpower(C.term_t(t1), C.uint32_t(d)))
}

func Bvdiv(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvdiv(C.term_t(t1), C.term_t(t2)))
}

func Bvrem(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvrem(C.term_t(t1), C.term_t(t2)))
}

func Bvsdiv(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvsdiv(C.term_t(t1), C.term_t(t2)))
}

func Bvsrem(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvsrem(C.term_t(t1), C.term_t(t2)))
}

func Bvsmod(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvsmod(C.term_t(t1), C.term_t(t2)))
}

func Bvnot(t Term_t) Term_t {
	return Term_t(C.yices_bvnot(C.term_t(t)))
}

func Bvnand(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvnand(C.term_t(t1), C.term_t(t2)))
}

func Bvnor(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvnor(C.term_t(t1), C.term_t(t2)))
}

func Bvxnor(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvxnor(C.term_t(t1), C.term_t(t2)))
}

func Bvshl(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvshl(C.term_t(t1), C.term_t(t2)))
}

func Bvlshr(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvlshr(C.term_t(t1), C.term_t(t2)))
}

func Bvashr(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvashr(C.term_t(t1), C.term_t(t2)))
}

func Bvand(t []Term_t) Term_t {
	count := C.uint32_t(len(t))
	return Term_t(C.yices_bvand(count, (*C.term_t)(&t[0])))
}

func Bvor(t []Term_t) Term_t {
	count := C.uint32_t(len(t))
	return Term_t(C.yices_bvor(count, (*C.term_t)(&t[0])))
}

func Bvxor(t []Term_t) Term_t {
	count := C.uint32_t(len(t))
	return Term_t(C.yices_bvxor(count, (*C.term_t)(&t[0])))
}

func Bvand2(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvand2(C.term_t(t1), C.term_t(t2)))
}

func Bvor2(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvor2(C.term_t(t1), C.term_t(t2)))
}

func Bvxor2(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvxor2(C.term_t(t1), C.term_t(t2)))
}

func Bvand3(t1 Term_t, t2 Term_t, t3 Term_t) Term_t {
	return Term_t(C.yices_bvand3(C.term_t(t1), C.term_t(t2), C.term_t(t3)))
}

func Bvor3(t1 Term_t, t2 Term_t, t3 Term_t) Term_t {
	return Term_t(C.yices_bvor3(C.term_t(t1), C.term_t(t2), C.term_t(t3)))
}

func Bvsum(t []Term_t) Term_t {
	count := C.uint32_t(len(t))
	return Term_t(C.yices_bvsum(count, (*C.term_t)(&t[0])))
}

func Bvproduct(t []Term_t) Term_t {
	count := C.uint32_t(len(t))
	return Term_t(C.yices_bvproduct(count, (*C.term_t)(&t[0])))
}

func Shift_left0(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_shift_left0(C.term_t(t), C.uint32_t(n)))
}

func Shift_left1(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_shift_left1(C.term_t(t), C.uint32_t(n)))
}

func Shift_right0(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_shift_right0(C.term_t(t), C.uint32_t(n)))
}

func Shift_right1(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_shift_right1(C.term_t(t), C.uint32_t(n)))
}

func Ashift_right(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_ashift_right(C.term_t(t), C.uint32_t(n)))
}

func Rotate_left(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_rotate_left(C.term_t(t), C.uint32_t(n)))
}

func Rotate_right(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_rotate_right(C.term_t(t), C.uint32_t(n)))
}

func Bvextract(t Term_t, i uint32, j uint32) Term_t {
	return Term_t(C.yices_bvextract(C.term_t(t), C.uint32_t(i), C.uint32_t(j)))
}

func Bvconcat2(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvconcat2(C.term_t(t1), C.term_t(t2)))
}

func Bvconcat(t []Term_t) Term_t {
	count := C.uint32_t(len(t))
	return Term_t(C.yices_bvconcat(count, (*C.term_t)(&t[0])))
}

func Bvrepeat(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_bvrepeat(C.term_t(t),  C.uint32_t(n)))
}

func Sign_extend(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_sign_extend(C.term_t(t),  C.uint32_t(n)))
}

func Zero_extend(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_zero_extend(C.term_t(t),  C.uint32_t(n)))
}

func Redand(t Term_t) Term_t {
	return Term_t(C.yices_redand(C.term_t(t)))
}

func Redor(t Term_t) Term_t {
	return Term_t(C.yices_redor(C.term_t(t)))
}

func Redcomp(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_redcomp(C.term_t(t1), C.term_t(t2)))
}

func Bvarray(t []Term_t) Term_t {
	count := C.uint32_t(len(t))
	return Term_t(C.yices_bvarray(count, (*C.term_t)(&t[0])))
}

func Bitextract(t Term_t, n uint32) Term_t {
	return Term_t(C.yices_bitextract(C.term_t(t),  C.uint32_t(n)))
}

func Bveq_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bveq_atom(C.term_t(t1), C.term_t(t2)))
}

func Bvneq_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvneq_atom(C.term_t(t1), C.term_t(t2)))
}

func Bvge_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvge_atom(C.term_t(t1), C.term_t(t2)))
}

func Bvgt_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvgt_atom(C.term_t(t1), C.term_t(t2)))
}

func Bvle_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvle_atom(C.term_t(t1), C.term_t(t2)))
}

func Bvlt_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvlt_atom(C.term_t(t1), C.term_t(t2)))
}

func Bvsge_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvsge_atom(C.term_t(t1), C.term_t(t2)))
}

func Bvsgt_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvsgt_atom(C.term_t(t1), C.term_t(t2)))
}

func Bvsle_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvsle_atom(C.term_t(t1), C.term_t(t2)))
}

func Bvslt_atom(t1 Term_t, t2 Term_t) Term_t {
	return Term_t(C.yices_bvslt_atom(C.term_t(t1), C.term_t(t2)))
}

/**************
 *  PARSING   *
 *************/

func Parse_type(s string) Type_t {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return Type_t(C.yices_parse_type(cs))
}

func Parse_term(s string) Term_t {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return Term_t(C.yices_parse_term(cs))
}

/*******************
 *  SUBSTITUTIONS  *
 ******************/

func Subst_term(vars []Term_t, vals []Term_t, t Term_t) Term_t {
	count := C.uint32_t(len(vars))
	return Term_t(C.yices_subst_term(count, (*C.term_t)(&vars[0]), (*C.term_t)(&vals[0]), C.term_t(t)))
}

func Subst_term_array(vars []Term_t, vals []Term_t, t []Term_t) Term_t {
	count := C.uint32_t(len(vars))
	tcount := C.uint32_t(len(t))
	return Term_t(C.yices_subst_term_array(count, (*C.term_t)(&vars[0]), (*C.term_t)(&vals[0]), tcount, (*C.term_t)(&t[0])))
}


/************
 *  NAMES   *
 ***********/

func Set_type_name(tau Type_t, name string) int32 {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return int32(C.yices_set_type_name(C.type_t(tau), cs))
}

func Set_term_name(t Term_t, name string) int32 {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return int32(C.yices_set_term_name(C.term_t(t), cs))
}

func Remove_type_name(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.yices_remove_type_name(cs)
}

func Remove_term_name(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.yices_remove_term_name(cs)
}

func Get_type_by_name(name string) Type_t {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return Type_t(C.yices_get_type_by_name(cs))
}

func Get_term_by_name(name string) Term_t {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return Term_t(C.yices_get_term_by_name(cs))
}

func Clear_type_name(tau Type_t) int32 {
	return int32(C.yices_clear_type_name(C.type_t(tau)))
}

func Clear_term_name(t Term_t) int32 {
	return int32(C.yices_clear_term_name(C.term_t(t)))
}

func Get_type_name(tau Type_t) string {
	//FIXME: check if the name needs to be freed
	return C.GoString(C.yices_get_type_name(C.type_t(tau)))
}

/***********************
 *  TERM EXPLORATION   *
 **********************/


func Type_of_term(t Term_t) Type_t {
	return Type_t(C.yices_type_of_term(C.term_t(t)))
}


func Term_is_bool(t Term_t) int32 {
	return int32(C.yices_term_is_bool(C.term_t(t)))
}

func Term_is_int(t Term_t) int32 {
	return int32(C.yices_term_is_int(C.term_t(t)))
}

func Term_is_real(t Term_t) int32 {
	return int32(C.yices_term_is_real(C.term_t(t)))
}

func Term_is_arithmetic(t Term_t) int32 {
	return int32(C.yices_term_is_arithmetic(C.term_t(t)))
}

func Term_is_bitvector(t Term_t) int32 {
	return int32(C.yices_term_is_bitvector(C.term_t(t)))
}

func Term_is_tuple(t Term_t) int32 {
	return int32(C.yices_term_is_tuple(C.term_t(t)))
}

func Term_is_function(t Term_t) int32 {
	return int32(C.yices_term_is_function(C.term_t(t)))
}

func Term_is_scalar(t Term_t) int32 {
	return int32(C.yices_term_is_scalar(C.term_t(t)))
}

func Term_bitsize(t Term_t) uint32 {
	return uint32(C.yices_term_bitsize(C.term_t(t)))
}

func Term_is_ground(t Term_t) int32 {
	return int32(C.yices_term_is_ground(C.term_t(t)))
}

func Term_is_atomic(t Term_t) int32 {
	return int32(C.yices_term_is_atomic(C.term_t(t)))
}

func Term_is_composite(t Term_t) int32 {
	return int32(C.yices_term_is_composite(C.term_t(t)))
}

func Term_is_projection(t Term_t) int32 {
	return int32(C.yices_term_is_projection(C.term_t(t)))
}

func Term_is_sum(t Term_t) int32 {
	return int32(C.yices_term_is_sum(C.term_t(t)))
}

func Term_is_bvsum(t Term_t) int32 {
	return int32(C.yices_term_is_bvsum(C.term_t(t)))
}

func Term_is_product(t Term_t) int32 {
	return int32(C.yices_term_is_product(C.term_t(t)))
}

func Term_constructor(t Term_t) Term_constructor_t {
	return Term_constructor_t(C.yices_term_constructor(C.term_t(t)))
}

func Term_num_children(t Term_t) int32 {
	return int32(C.yices_term_num_children(C.term_t(t)))
}

func Term_child(t Term_t, i int32) Term_t {
	return Term_t(C.yices_term_child(C.term_t(t), C.int32_t(i)))
}

func Proj_index(t Term_t) int32 {
	return int32(C.yices_proj_index(C.term_t(t)))
}

func Proj_arg(t Term_t) Term_t {
	return Term_t(C.yices_proj_arg(C.term_t(t)))
}


func Bool_const_value(t Term_t, val *int32) int32  {
	return int32(C.yices_bool_const_value(C.term_t(t), (* C.int32_t)(val)))
}

func Bv_const_value(t Term_t, val []int32) int32 {
	return int32(C.yices_bv_const_value(C.term_t(t), (* C.int32_t)(&val[0])))
}

func Scalar_const_value(t Term_t, val *int32) int32  {
	return int32(C.yices_scalar_const_value(C.term_t(t), (* C.int32_t)(val)))
}

/* iam: FIXME
#ifdef __GMP_H__
__YICES_DLLSPEC__ extern int32_t yices_rational_const_value(term_t t, mpq_t q);
#endif
*/

/* iam: FIXME
#ifdef __GMP_H__
__YICES_DLLSPEC__ extern int32_t yices_sum_component(term_t t, int32_t i, mpq_t coeff, term_t *term);
#endif
*/

func Bvsum_component(t Term_t, i int32, val []int32, term *Term_t) int32 {
	return int32(C.yices_bvsum_component(C.term_t(t), C.int32_t(i), (* C.int32_t)(&val[0]), (*C.term_t)(term)))
}

func Product_component(t Term_t, i int32, term *Term_t, exp *uint32) int32 {
	return int32(C.yices_product_component(C.term_t(t), C.int32_t(i), (* C.term_t)(term), (* C.uint32_t)(exp)))
}

/*************************
 *  GARBAGE COLLECTION   *
 ************************/

func Num_terms() uint32 {
	return uint32(C.yices_num_terms())
}

func Num_types() uint32 {
	return uint32(C.yices_num_types())
}

func Incref_term(t Term_t) {
	C.yices_incref_term(C.term_t(t))
}

func Decref_term(t Term_t) {
	C.yices_decref_term(C.term_t(t))
}

func Incref_type(tau Type_t) {
	C.yices_incref_type(C.type_t(tau))
}

func Decref_type(tau Type_t) {
	C.yices_decref_type(C.type_t(tau))
}

func Num_posref_terms() uint32 {
	return uint32(C.yices_num_posref_terms())
}

func Num_posref_types() uint32 {
	return uint32(C.yices_num_posref_types())
}


func Garbage_collect(ts []Term_t, taus []Type_t,  keep_named int32) {
	t_count := C.uint32_t(len(ts))
	tau_count := C.uint32_t(len(taus))
	C.yices_garbage_collect((* C.term_t)(&ts[0]), t_count, (* C.type_t)(&taus[0]), tau_count, C.int32_t(keep_named))
} 


/****************************
 *  CONTEXT CONFIGURATION   *
 ***************************/


type Config C.ctx_config_t

func New_config() *Config {
	return (* Config)(C.yices_new_config())
}

func Free_config(cfg  *Config) {
	C.yices_free_config((*C.ctx_config_t)(cfg))
}

func Set_config(cfg  *Config, name string, value string) int32 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	return int32(C.yices_set_config((*C.ctx_config_t)(cfg), cname, cvalue))
}

func Default_config_for_logic(cfg  *Config, logic string) int32 {
	clogic := C.CString(logic)
	defer C.free(unsafe.Pointer(clogic))
	return int32(C.yices_default_config_for_logic((*C.ctx_config_t)(cfg), clogic))
}

