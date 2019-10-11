package yices2

/*
#cgo CFLAGS: -g -fPIC
#cgo LDFLAGS:  -lyices -lgmp
#include <yices.h>
type_t yices_type_vector_get(type_vector_t* vec, uint32_t elem){ return vec->data[elem]; }
*/
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
	return Term_t(C.yices_parse_rational(C.CString(s)))
}

func Parse_float(s string) Term_t {
	return Term_t(C.yices_parse_float(C.CString(s)))
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
