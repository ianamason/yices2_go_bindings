package yices2

/*
#cgo CFLAGS: -g -fPIC
#cgo LDFLAGS:  -lyices -lgmp
#include <stdlib.h>
#include <gmp.h>
#include <yices.h>
//iam: avoid ugly pointer arithmetic
type_t yices_type_vector_get(type_vector_t* vec, uint32_t elem){ return vec->data[elem]; }
term_t yices_term_vector_get(term_vector_t* vec, uint32_t elem){ return vec->data[elem]; }
void yices_yval_vector_get(yval_vector_t *vec, uint32_t elem, yval_t* val){
      yval_t *v = &vec->data[elem];
      *val = *v;
}
//iam: hack to get around weird cgo complaints
term_t ympz(uintptr_t ptr){
   return yices_mpz(*((mpz_t *)((void *)ptr)));
}
//iam: hack to get around weird cgo complaints
term_t ympq(uintptr_t ptr){
   return yices_mpq(*((mpq_t *)((void *)ptr)));
}
//iam: minimal interface to gmp's mpz
void init_mpzp(uintptr_t ptr){
   mpz_init(*(mpz_t *)(ptr));
}
void close_mpzp(uintptr_t ptr){
   mpz_clear(*(mpz_t *)(ptr));
}
//iam: minimal interface to gmp's mpq
void init_mpqp(uintptr_t ptr){
   mpq_init(*(mpq_t *)(ptr));
}
void close_mpqp(uintptr_t ptr){
   mpq_clear(*(mpq_t *)(ptr));
}
//iam: passing a mpz_t thru cgo seems too hard
term_t yices_bvconst_mpzp(uint32_t n, mpz_t *x){
   return yices_bvconst_mpz(n, *x);
}
//iam: passing a mpq_t thru cgo seems too hard
int32_t yices_rational_const_valuep(term_t t, mpq_t *q){
   return yices_rational_const_value(t, *q);
}
//iam: passing a mpq_t thru cgo seems too hard
int32_t yices_sum_componentp(term_t t, int32_t i, mpq_t *coeff, term_t *term){
   return yices_sum_component(t, i, *coeff, term);
}
//iam: passing a mpz_t thru cgo seems too hard
int32_t yices_get_mpz_valuep(model_t *mdl, term_t t, mpz_t *val){
   return yices_get_mpz_value(mdl, t, *val);
}
//iam: passing a mpq_t thru cgo seems too hard
int32_t yices_get_mpq_valuep(model_t *mdl, term_t t, mpq_t *val){
   return yices_get_mpq_value(mdl, t, *val);
}
//iam: passing a mpz_t thru cgo seems too hard
int32_t yices_val_get_mpzp(model_t *mdl, const yval_t *v, mpz_t *val){
   return yices_val_get_mpz(mdl, v, *val);
}
//iam: passing a mpq_t thru cgo seems too hard
int32_t yices_val_get_mpqp(model_t *mdl, const yval_t *v, mpq_t *val){
   return yices_val_get_mpq(mdl, v, *val);
}
//iam: is there a better way than this?
void fetchReport(error_code_t *code, uint32_t *line, uint32_t *column, term_t *term1, type_t *type1, term_t *term2, type_t *type2, int64_t *badval){
   error_report_t *report = yices_error_report();
   *code = report->code;
   *line = report->line;
   *column = report->column;
   *term1 = report->term1;
   *type1 = report->type1;
   *term2 = report->term2;
   *type2 = report->type2;
   *badval = report->badval;
}
*/
import "C"

import "os"
import "fmt"
import "unsafe"

/*
 *  See yices.h for comments.
 *
 *  Naming convention:   yices_foo  becomes yices2.Foo  (will probably ditch the 2 at some stage)
 *  The exception to this rule is the new and free routines for contexts configs and params; that have
 *  Init_<kind> and Close_<kind> routines.
 *
 *  This layer (yices_api.go) is a thin wrapper to the yices_api. Maybe a more go-like layer will sit atop this
 *  much like the python version of the API.
 *
 *  Currently backward compatibility is in the "too hard" basket. We are going to assume the yices library is
 *  at least 2.6.2 (look for the "Since 2.6.2" comments)
 *
 */

/*********************
 *  VERSION NUMBERS  *
 ********************/

// Version is the yices2 library version.
func Version() string {
	return C.GoString(C.yices_version)
}

// BuildArch is the yices2 library build architecture.
func BuildArch() string {
	return C.GoString(C.yices_build_arch)
}

// BuildMode is the yices2 library build mode.
func BuildMode() string {
	return C.GoString(C.yices_build_mode)
}

// BuildDate is the yices2 library build date.
func BuildDate() string {
	return C.GoString(C.yices_build_date)
}

// HasMcsat indicates if the yices2 library supports MCSAT.
func HasMcsat() int32 {
	return int32(C.yices_has_mcsat())
}

// IsThreadSafe indicate if the yices2 library was built with thread safety enabled.
// Since 2.6.2
func IsThreadSafe() int32 {
	return int32(C.yices_is_thread_safe())
}

/***************************************
 *  GLOBAL INITIALIZATION AND CLEANUP  *
 **************************************/

// Init initializes the internal yices2 library data structures.
func Init() {
	C.yices_init()
}

// Exit cleans up the internal yices2 library data structures.
func Exit() {
	C.yices_exit()
}

// Reset resets up the internal yices2 library data structures.
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

// ErrorCodeT is the analog of the error_code_t enum of the yices library.
type ErrorCodeT int32

// YicesErrorT is the go representation of the yices_error_t struct of yices library.
type YicesErrorT struct {
	ErrorString string
	Code        ErrorCodeT
	Line        uint32
	Column      uint32
	Term1       TermT
	Type1       TypeT
	Term2       TermT
	Type2       TypeT
	Badval      int64
}

// Error returns the error's current error string.
func (yerror *YicesErrorT) Error() string {
	return yerror.ErrorString
}

func fetchErrorReport(yerror *YicesErrorT) {
	var code C.error_code_t
	var line C.uint32_t
	var column C.uint32_t
	var term1 C.term_t
	var type1 C.type_t
	var term2 C.term_t
	var type2 C.type_t
	var badval C.int64_t
	C.fetchReport(&code, &line, &column, &term1, &type1, &term2, &type2, &badval)
	yerror.Code = ErrorCodeT(code)
	yerror.Line = uint32(line)
	yerror.Column = uint32(column)
	yerror.Term1 = TermT(term1)
	yerror.Type1 = TypeT(type1)
	yerror.Term2 = TermT(term2)
	yerror.Type2 = TypeT(type2)
	yerror.Badval = int64(badval)
	return
}

// YicesError returns a copy of the current error state
func YicesError() (yerror *YicesErrorT) {
	errcode := ErrorCode()
	if errcode != NoError {
		yerror = new(YicesErrorT)
		yerror.ErrorString = ErrorString()
		fetchErrorReport(yerror)
		ClearError()
	}
	return
}

// String returns the error string component of the YicesErrorT struct.
func (yerror *YicesErrorT) String() string {
	return fmt.Sprintf("string = %s code = %d line = %d column = %d term1 = %d type1 = %d term2 = %d type2 = %d badval = %d",
		yerror.ErrorString, yerror.Code, yerror.Line, yerror.Column,
		yerror.Term1, yerror.Type1,
		yerror.Term2, yerror.Type2,
		yerror.Badval)
}

// ErrorCode returns the most recent yices error code.
func ErrorCode() ErrorCodeT {
	return ErrorCodeT(C.yices_error_code())
}

// ClearError reset the current yices error struct.
func ClearError() {
	C.yices_clear_error()
}

// PrintError prints the most recent yices error out to the given file.
func PrintError(f *os.File) int32 {
	return int32(C.yices_print_error_fd(C.int(f.Fd())))
} //iam: FIXME error checking and File without the os.

func ErrorString() string {
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

//iam: we use a type definition, does it improve readability?
// does it make life harder? Lots of casting necessary now.

type TypeT int32

const NullType TypeT = -1

func BoolType() TypeT {
	return TypeT(C.yices_bool_type())
}

func IntType() TypeT {
	return TypeT(C.yices_int_type())
}

func RealType() TypeT {
	return TypeT(C.yices_real_type())
}

func BvType(size uint32) TypeT {
	return TypeT(C.yices_bv_type(C.uint32_t(size)))
}

func NewScalarType(card uint32) TypeT {
	return TypeT(C.yices_new_scalar_type(C.uint32_t(card)))
}

func NewUninterpretedType() TypeT {
	return TypeT(C.yices_new_uninterpreted_type())
}

func TupleType(tau []TypeT) TypeT {
	tauLen := len(tau)
	//iam: FIXME need to unify the yices errors and the go errors...
	if tauLen == 0 {
		return NullType
	}
	return TypeT(C.yices_tuple_type(C.uint32_t(tauLen), (*C.type_t)(&tau[0])))
}

func TupleType1(tau1 TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1)}
	return TypeT(C.yices_tuple_type(C.uint32_t(1), (*C.type_t)(&carr[0])))
}

func TupleType2(tau1 TypeT, tau2 TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2)}
	return TypeT(C.yices_tuple_type(C.uint32_t(2), (*C.type_t)(&carr[0])))
}

func TupleType3(tau1 TypeT, tau2 TypeT, tau3 TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2), C.type_t(tau3)}
	return TypeT(C.yices_tuple_type(C.uint32_t(3), (*C.type_t)(&carr[0])))
}

func FunctionType(dom []TypeT, rng TypeT) TypeT {
	domLen := len(dom)
	//iam: FIXME need to unify the yices errors and the go errors...
	if domLen == 0 {
		return NullType
	}
	return TypeT(C.yices_function_type(C.uint32_t(domLen), (*C.type_t)(&dom[0]), C.type_t(rng)))
}

func FunctionType1(tau1 TypeT, rng TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1)}
	return TypeT(C.yices_function_type(C.uint32_t(1), (*C.type_t)(&carr[0]), C.type_t(rng)))
}

func FunctionType2(tau1 TypeT, tau2 TypeT, rng TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2)}
	return TypeT(C.yices_function_type(C.uint32_t(2), (*C.type_t)(&carr[0]), C.type_t(rng)))
}

func FunctionType3(tau1 TypeT, tau2 TypeT, tau3 TypeT, rng TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2), C.type_t(tau3)}
	return TypeT(C.yices_function_type(C.uint32_t(3), (*C.type_t)(&carr[0]), C.type_t(rng)))
}

/*************************
 *   TYPE EXPLORATION    *
 ************************/

func TypeIsBool(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_bool(C.type_t(tau)))
}

func TypeIsInt(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_int(C.type_t(tau)))
}

func TypeIsReal(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_real(C.type_t(tau)))
}

func TypeIsArithmetic(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_arithmetic(C.type_t(tau)))
}

func TypeIsBitvector(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_bitvector(C.type_t(tau)))
}

func TypeIsTuple(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_tuple(C.type_t(tau)))
}

func TypeIsFunction(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_function(C.type_t(tau)))
}

func TypeIsScalar(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_scalar(C.type_t(tau)))
}

func TypeIsUninterpreted(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_uninterpreted(C.type_t(tau)))
}

func TestSubtype(tau TypeT, sigma TypeT) bool {
	return int32(1) == int32(C.yices_test_subtype(C.type_t(tau), C.type_t(sigma)))
}

func CompatibleTypes(tau TypeT, sigma TypeT) bool {
	return int32(1) == int32(C.yices_compatible_types(C.type_t(tau), C.type_t(sigma)))
}

func BvtypeSize(tau TypeT) uint32 {
	return uint32(C.yices_bvtype_size(C.type_t(tau)))
}

func ScalarTypeCard(tau TypeT) uint32 {
	return uint32(C.yices_scalar_type_card(C.type_t(tau)))
}

func TypeNumChildren(tau TypeT) int32 {
	return int32(C.yices_type_num_children(C.type_t(tau)))
}

func TypeChild(tau TypeT, i int32) TypeT {
	return TypeT(C.yices_type_child(C.type_t(tau), C.int32_t(i)))
}

func TypeChildren(tau TypeT) (children []TypeT) {
	var tv C.type_vector_t
	C.yices_init_type_vector(&tv)
	ycount := int32(C.yices_type_children(C.type_t(tau), &tv))
	if ycount != -1 {
		count := int(tv.size)
		children = make([]TypeT, count, count)
		// defined in the preamble yices_type_vector_get(type_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			children[i] = TypeT(C.yices_type_vector_get(&tv, C.uint32_t(i)))
		}
	}
	C.yices_delete_type_vector(&tv)
	return
}

/***********************
 *  TERM CONSTRUCTORS  *
 **********************/

type TermT int32

const NullTerm TermT = -1

func True() TermT {
	return TermT(C.yices_true())
}

func False() TermT {
	return TermT(C.yices_false())
}

func Constant(tau TypeT, index int32) TermT {
	return TermT(C.yices_constant(C.type_t(tau), C.int32_t(index)))
}

func NewUninterpretedTerm(tau TypeT) TermT {
	return TermT(C.yices_new_uninterpreted_term(C.type_t(tau)))
}

func NewVariable(tau TypeT) TermT {
	return TermT(C.yices_new_variable(C.type_t(tau)))
}

func Application(fun TermT, argv []TermT) TermT {
	argc := len(argv)
	//iam: FIXME need to unify the yices errors and the go errors...
	if argc == 0 {
		return NullTerm
	}
	return TermT(C.yices_application(C.term_t(fun), C.uint32_t(argc), (*C.term_t)(&argv[0])))
}

func Application1(fun TermT, arg1 TermT) TermT {
	argv := []C.term_t{C.term_t(arg1)}
	return TermT(C.yices_application(C.term_t(fun), C.uint32_t(1), (*C.term_t)(&argv[0])))
}

func Application2(fun TermT, arg1 TermT, arg2 TermT) TermT {
	argv := []C.term_t{C.term_t(arg1), C.term_t(arg2)}
	return TermT(C.yices_application(C.term_t(fun), C.uint32_t(2), (*C.term_t)(&argv[0])))
}

func Application3(fun TermT, arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	argv := []C.term_t{C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)}
	return TermT(C.yices_application(C.term_t(fun), C.uint32_t(3), (*C.term_t)(&argv[0])))
}

func Ite(cond TermT, thenTerm TermT, elseTerm TermT) TermT {
	return TermT(C.yices_ite(C.term_t(cond), C.term_t(thenTerm), C.term_t(elseTerm)))
}

func Eq(lhs TermT, rhs TermT) TermT {
	return TermT(C.yices_eq(C.term_t(lhs), (C.term_t(rhs))))
}

func Neq(lhs TermT, rhs TermT) TermT {
	return TermT(C.yices_neq(C.term_t(lhs), (C.term_t(rhs))))
}

func Not(arg TermT) TermT {
	return TermT(C.yices_not(C.term_t(arg)))
}

func Or(disjuncts []TermT) TermT {
	count := C.uint32_t(len(disjuncts))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return TermT(C.yices_false())
	}
	return TermT(C.yices_or(count, (*C.term_t)(&disjuncts[0])))
}

func And(conjuncts []TermT) TermT {
	count := C.uint32_t(len(conjuncts))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return TermT(C.yices_true())
	}
	return TermT(C.yices_and(count, (*C.term_t)(&conjuncts[0])))
}

func Xor(xorjuncts []TermT) TermT {
	count := C.uint32_t(len(xorjuncts))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		//FIXME what is xor of an empty array
		var dummy = C.yices_true()
		return TermT(C.yices_xor(count, &dummy))
	}
	return TermT(C.yices_xor(count, (*C.term_t)(&xorjuncts[0])))
}

func Or2(arg1 TermT, arg2 TermT) TermT {
	return TermT(C.yices_or2(C.term_t(arg1), C.term_t(arg2)))
}

func And2(arg1 TermT, arg2 TermT) TermT {
	return TermT(C.yices_and2(C.term_t(arg1), C.term_t(arg2)))
}

func Xor2(arg1 TermT, arg2 TermT) TermT {
	return TermT(C.yices_xor2(C.term_t(arg1), C.term_t(arg2)))
}

func Or3(arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	return TermT(C.yices_or3(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

func And3(arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	return TermT(C.yices_and3(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

func Xor3(arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	return TermT(C.yices_xor3(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

func Iff(lhs TermT, rhs TermT) TermT {
	return TermT(C.yices_iff(C.term_t(lhs), (C.term_t(rhs))))
}

func Implies(lhs TermT, rhs TermT) TermT {
	return TermT(C.yices_implies(C.term_t(lhs), (C.term_t(rhs))))
}

func Tuple(argv []TermT) TermT {
	count := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_tuple(count, (*C.term_t)(&argv[0])))
}

func Pair(arg1 TermT, arg2 TermT) TermT {
	return TermT(C.yices_pair(C.term_t(arg1), C.term_t(arg2)))
}

func Triple(arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	return TermT(C.yices_triple(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

func Select(index uint32, tuple TermT) TermT {
	return TermT(C.yices_select(C.uint32_t(index), C.term_t(tuple)))
}

func TupleUpdate(tuple TermT, index uint32, value TermT) TermT {
	return TermT(C.yices_tuple_update(C.term_t(tuple), C.uint32_t(index), C.term_t(value)))
}

func Update(fun TermT, argv []TermT, value TermT) TermT {
	count := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_update(C.term_t(fun), count, (*C.term_t)(&argv[0]), C.term_t(value)))
}

func Update1(fun TermT, arg1 TermT, value TermT) TermT {
	return TermT(C.yices_update1(C.term_t(fun), C.term_t(arg1), C.term_t(value)))
}

func Update2(fun TermT, arg1 TermT, arg2 TermT, value TermT) TermT {
	return TermT(C.yices_update2(C.term_t(fun), C.term_t(arg1), C.term_t(arg2), C.term_t(value)))
}

func Update3(fun TermT, arg1 TermT, arg2 TermT, arg3 TermT, value TermT) TermT {
	return TermT(C.yices_update3(C.term_t(fun), C.term_t(arg1), C.term_t(arg2), C.term_t(arg3), C.term_t(value)))
}

func Distinct(argv []TermT) TermT {
	n := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if n == 0 {
		return NullTerm
	}
	return TermT(C.yices_distinct(n, (*C.term_t)(&argv[0])))
}

func Forall(vars []TermT, body TermT) TermT {
	n := C.uint32_t(len(vars))
	//iam: FIXME need to unify the yices errors and the go errors...
	if n == 0 {
		return NullTerm
	}
	return TermT(C.yices_forall(n, (*C.term_t)(&vars[0]), C.term_t(body)))
}

func Exists(vars []TermT, body TermT) TermT {
	n := C.uint32_t(len(vars))
	//iam: FIXME need to unify the yices errors and the go errors...
	if n == 0 {
		return NullTerm
	}
	return TermT(C.yices_exists(n, (*C.term_t)(&vars[0]), C.term_t(body)))
}

func Lambda(vars []TermT, body TermT) TermT {
	n := C.uint32_t(len(vars))
	//iam: FIXME need to unify the yices errors and the go errors...
	if n == 0 {
		return NullTerm
	}
	return TermT(C.yices_lambda(n, (*C.term_t)(&vars[0]), C.term_t(body)))
}

/**********************************
 *  ARITHMETIC TERM CONSTRUCTORS  *
 *********************************/

func Zero() TermT {
	return TermT(C.yices_zero())
}

func Int32(val int32) TermT {
	return TermT(C.yices_int32(C.int32_t(val)))
}

func Int64(val int64) TermT {
	return TermT(C.yices_int64(C.int64_t(val)))
}

func Rational32(num int32, den uint32) TermT {
	return TermT(C.yices_rational32(C.int32_t(num), C.uint32_t(den)))
}

func Rational64(num int64, den uint64) TermT {
	return TermT(C.yices_rational64(C.int64_t(num), C.uint64_t(den)))
}

// need to name these to use them outside this package
type MpzT C.mpz_t

func InitMpz(mpz *MpzT) {
	C.init_mpzp(C.uintptr_t(uintptr(unsafe.Pointer(mpz))))
}

func CloseMpz(mpz *MpzT) {
	C.close_mpzp(C.uintptr_t(uintptr(unsafe.Pointer(mpz))))
}

func Mpz(z *MpzT) TermT {
	// some contortions needed here to do the simplest of things
	// similar contortions are needed to "identify" yices2._Ctype_mpz_t
	// with gmp._Ctype_mpz_t, which seems a little odd.
	return TermT(C.ympz(C.uintptr_t(uintptr(unsafe.Pointer(z)))))
}

// need to name these to use them outside this package
type MpqT C.mpq_t

func InitMpq(mpq *MpqT) {
	C.init_mpqp(C.uintptr_t(uintptr(unsafe.Pointer(mpq))))
}

func CloseMpq(mpq *MpqT) {
	C.close_mpqp(C.uintptr_t(uintptr(unsafe.Pointer(mpq))))
}

func Mpq(q *MpqT) TermT {
	// some contortions needed here to do the simplest of things
	return TermT(C.ympq(C.uintptr_t(uintptr(unsafe.Pointer(q)))))
}

func ParseRational(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_rational(cs))
}

func ParseFloat(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_float(cs))
}

/*
 * ARITHMETIC OPERATIONS
 */

func Add(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_add(C.term_t(t1), C.term_t(t2)))
}

func Sub(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_sub(C.term_t(t1), C.term_t(t2)))
}

func Neg(t1 TermT) TermT {
	return TermT(C.yices_neg(C.term_t(t1)))
}

func Mul(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_mul(C.term_t(t1), C.term_t(t2)))
}

func Square(t1 TermT) TermT {
	return TermT(C.yices_square(C.term_t(t1)))
}

func Power(t1 TermT, d uint32) TermT {
	return TermT(C.yices_power(C.term_t(t1), C.uint32_t(d)))
}

func Sum(argv []TermT) TermT {
	count := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_sum(count, (*C.term_t)(&argv[0])))
}

func Product(argv []TermT) TermT {
	count := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return TermT(C.yices_int32(1))
	}
	return TermT(C.yices_product(count, (*C.term_t)(&argv[0])))
}

func Division(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_division(C.term_t(t1), C.term_t(t2)))
}

func Idiv(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_idiv(C.term_t(t1), C.term_t(t2)))
}

func Imod(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_imod(C.term_t(t1), C.term_t(t2)))
}

func DividesAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_divides_atom(C.term_t(t1), C.term_t(t2)))
}

func IsIntAtom(t TermT) TermT {
	return TermT(C.yices_is_int_atom(C.term_t(t)))
}

func Abs(t1 TermT) TermT {
	return TermT(C.yices_abs(C.term_t(t1)))
}

func Floor(t1 TermT) TermT {
	return TermT(C.yices_floor(C.term_t(t1)))
}

func Ceil(t1 TermT) TermT {
	return TermT(C.yices_ceil(C.term_t(t1)))
}

/*
 * POLYNOMIALS
 */

func PolyInt32(a []int32, t []TermT) TermT {
	count := C.uint32_t(len(a))
	//iam: FIXME need to unify the yices errors and the go errors...
	// do we want to be nannies here?
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_int32(count, (*C.int32_t)(&a[0]), (*C.term_t)(&t[0])))
}

func PolyInt64(a []int64, t []TermT) TermT {
	count := C.uint32_t(len(a))
	//iam: FIXME need to unify the yices errors and the go errors...
	// do we want to be nannies here?
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_int64(count, (*C.int64_t)(&a[0]), (*C.term_t)(&t[0])))
}

func PolyRational32(num []int32, den []uint32, t []TermT) TermT {
	count := C.uint32_t(len(num))
	//iam: FIXME need to unify the yices errors and the go errors...
	// do we want to be nannies here?
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_rational32(count, (*C.int32_t)(&num[0]), (*C.uint32_t)(&den[0]), (*C.term_t)(&t[0])))
}

func PolyRational64(num []int64, den []uint64, t []TermT) TermT {
	count := C.uint32_t(len(num))
	//iam: FIXME need to unify the yices errors and the go errors...
	// do we want to be nannies here?
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_rational64(count, (*C.int64_t)(&num[0]), (*C.uint64_t)(&den[0]), (*C.term_t)(&t[0])))
}

func PolyMpz(z []MpzT, t []TermT) TermT {
	count := C.uint32_t(len(z))
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_mpz(count, (*C.mpz_t)(&z[0]), (*C.term_t)(&t[0])))
}

func PolyMpq(q []MpqT, t []TermT) TermT {
	count := C.uint32_t(len(q))
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_mpq(count, (*C.mpq_t)(&q[0]), (*C.term_t)(&t[0])))
}

/*
 * ARITHMETIC ATOMS
 */

func ArithEqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_eq_atom(C.term_t(t1), C.term_t(t2)))
}

func ArithNeqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_neq_atom(C.term_t(t1), C.term_t(t2)))
}

func ArithGeqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_geq_atom(C.term_t(t1), C.term_t(t2)))
}

func ArithLeqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_leq_atom(C.term_t(t1), C.term_t(t2)))
}

func ArithGtAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_gt_atom(C.term_t(t1), C.term_t(t2)))
}

func ArithLtAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_lt_atom(C.term_t(t1), C.term_t(t2)))
}

func ArithEq0Atom(t TermT) TermT {
	return TermT(C.yices_arith_eq0_atom(C.term_t(t)))
}

func ArithNeq0Atom(t TermT) TermT {
	return TermT(C.yices_arith_neq0_atom(C.term_t(t)))
}

func ArithGeq0Atom(t TermT) TermT {
	return TermT(C.yices_arith_geq0_atom(C.term_t(t)))
}

func ArithLeq0Atom(t TermT) TermT {
	return TermT(C.yices_arith_leq0_atom(C.term_t(t)))
}

func ArithGt0Atom(t TermT) TermT {
	return TermT(C.yices_arith_gt0_atom(C.term_t(t)))
}

func ArithLt0Atom(t TermT) TermT {
	return TermT(C.yices_arith_lt0_atom(C.term_t(t)))
}

/*********************************
 *  BITVECTOR TERM CONSTRUCTORS  *
 ********************************/

func BvconstUint32(bits uint32, x uint32) TermT {
	return TermT(C.yices_bvconst_uint32(C.uint32_t(bits), C.uint32_t(x)))
}

func BvconstUint64(bits uint32, x uint64) TermT {
	return TermT(C.yices_bvconst_uint64(C.uint32_t(bits), C.uint64_t(x)))
}

func BvconstInt32(bits uint32, x int32) TermT {
	return TermT(C.yices_bvconst_int32(C.uint32_t(bits), C.int32_t(x)))
}

func BvconstInt64(bits uint32, x int64) TermT {
	return TermT(C.yices_bvconst_int64(C.uint32_t(bits), C.int64_t(x)))
}

func BvconstMpz(bits uint32, z MpzT) TermT {
	return TermT(C.yices_bvconst_mpzp(C.uint32_t(bits), (*C.mpz_t)(unsafe.Pointer(&z))))
}

func BvconstZero(bits uint32) TermT {
	return TermT(C.yices_bvconst_zero(C.uint32_t(bits)))
}

func BvconstOne(bits uint32) TermT {
	return TermT(C.yices_bvconst_one(C.uint32_t(bits)))
}

func BvconstMinusOne(bits uint32) TermT {
	return TermT(C.yices_bvconst_minus_one(C.uint32_t(bits)))
}

//iam: FIXME check that bits is restricted to len(a)
func BvconstFromArray(a []int32) TermT {
	bits := C.uint32_t(len(a))
	//iam: FIXME need to unify the yices errors and the go errors...
	if bits == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvconst_from_array(bits, (*C.int32_t)(&a[0])))
}

func ParseBvbin(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_bvbin(cs))
}

func ParseBvhex(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_bvhex(cs))
}

func Bvadd(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvadd(C.term_t(t1), C.term_t(t2)))
}

func Bvsub(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsub(C.term_t(t1), C.term_t(t2)))
}

func Bvneg(t TermT) TermT {
	return TermT(C.yices_bvneg(C.term_t(t)))
}

func Bvmul(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvmul(C.term_t(t1), C.term_t(t2)))
}

func Bvsquare(t TermT) TermT {
	return TermT(C.yices_bvsquare(C.term_t(t)))
}

func Bvpower(t1 TermT, d uint32) TermT {
	return TermT(C.yices_bvpower(C.term_t(t1), C.uint32_t(d)))
}

func Bvdiv(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvdiv(C.term_t(t1), C.term_t(t2)))
}

func Bvrem(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvrem(C.term_t(t1), C.term_t(t2)))
}

func Bvsdiv(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsdiv(C.term_t(t1), C.term_t(t2)))
}

func Bvsrem(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsrem(C.term_t(t1), C.term_t(t2)))
}

func Bvsmod(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsmod(C.term_t(t1), C.term_t(t2)))
}

func Bvnot(t TermT) TermT {
	return TermT(C.yices_bvnot(C.term_t(t)))
}

func Bvnand(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvnand(C.term_t(t1), C.term_t(t2)))
}

func Bvnor(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvnor(C.term_t(t1), C.term_t(t2)))
}

func Bvxnor(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvxnor(C.term_t(t1), C.term_t(t2)))
}

func Bvshl(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvshl(C.term_t(t1), C.term_t(t2)))
}

func Bvlshr(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvlshr(C.term_t(t1), C.term_t(t2)))
}

func Bvashr(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvashr(C.term_t(t1), C.term_t(t2)))
}

func Bvand(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvand(count, (*C.term_t)(&t[0])))
}

func Bvor(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvor(count, (*C.term_t)(&t[0])))
}

func Bvxor(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvxor(count, (*C.term_t)(&t[0])))
}

func Bvand2(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvand2(C.term_t(t1), C.term_t(t2)))
}

func Bvor2(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvor2(C.term_t(t1), C.term_t(t2)))
}

func Bvxor2(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvxor2(C.term_t(t1), C.term_t(t2)))
}

func Bvand3(t1 TermT, t2 TermT, t3 TermT) TermT {
	return TermT(C.yices_bvand3(C.term_t(t1), C.term_t(t2), C.term_t(t3)))
}

func Bvor3(t1 TermT, t2 TermT, t3 TermT) TermT {
	return TermT(C.yices_bvor3(C.term_t(t1), C.term_t(t2), C.term_t(t3)))
}

func Bvxor3(t1 TermT, t2 TermT, t3 TermT) TermT {
	return TermT(C.yices_bvxor3(C.term_t(t1), C.term_t(t2), C.term_t(t3)))
}

func Bvsum(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvsum(count, (*C.term_t)(&t[0])))
}

func Bvproduct(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvproduct(count, (*C.term_t)(&t[0])))
}

func ShiftLeft0(t TermT, n uint32) TermT {
	return TermT(C.yices_shift_left0(C.term_t(t), C.uint32_t(n)))
}

func ShiftLeft1(t TermT, n uint32) TermT {
	return TermT(C.yices_shift_left1(C.term_t(t), C.uint32_t(n)))
}

func ShiftRight0(t TermT, n uint32) TermT {
	return TermT(C.yices_shift_right0(C.term_t(t), C.uint32_t(n)))
}

func ShiftRight1(t TermT, n uint32) TermT {
	return TermT(C.yices_shift_right1(C.term_t(t), C.uint32_t(n)))
}

func AshiftRight(t TermT, n uint32) TermT {
	return TermT(C.yices_ashift_right(C.term_t(t), C.uint32_t(n)))
}

func RotateLeft(t TermT, n uint32) TermT {
	return TermT(C.yices_rotate_left(C.term_t(t), C.uint32_t(n)))
}

func RotateRight(t TermT, n uint32) TermT {
	return TermT(C.yices_rotate_right(C.term_t(t), C.uint32_t(n)))
}

func Bvextract(t TermT, i uint32, j uint32) TermT {
	return TermT(C.yices_bvextract(C.term_t(t), C.uint32_t(i), C.uint32_t(j)))
}

func Bvconcat2(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvconcat2(C.term_t(t1), C.term_t(t2)))
}

func Bvconcat(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvconcat(count, (*C.term_t)(&t[0])))
}

func Bvrepeat(t TermT, n uint32) TermT {
	return TermT(C.yices_bvrepeat(C.term_t(t), C.uint32_t(n)))
}

func SignExtend(t TermT, n uint32) TermT {
	return TermT(C.yices_sign_extend(C.term_t(t), C.uint32_t(n)))
}

func ZeroExtend(t TermT, n uint32) TermT {
	return TermT(C.yices_zero_extend(C.term_t(t), C.uint32_t(n)))
}

func Redand(t TermT) TermT {
	return TermT(C.yices_redand(C.term_t(t)))
}

func Redor(t TermT) TermT {
	return TermT(C.yices_redor(C.term_t(t)))
}

func Redcomp(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_redcomp(C.term_t(t1), C.term_t(t2)))
}

func Bvarray(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvarray(count, (*C.term_t)(&t[0])))
}

func Bitextract(t TermT, n uint32) TermT {
	return TermT(C.yices_bitextract(C.term_t(t), C.uint32_t(n)))
}

func BveqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bveq_atom(C.term_t(t1), C.term_t(t2)))
}

func BvneqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvneq_atom(C.term_t(t1), C.term_t(t2)))
}

func BvgeAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvge_atom(C.term_t(t1), C.term_t(t2)))
}

func BvgtAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvgt_atom(C.term_t(t1), C.term_t(t2)))
}

func BvleAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvle_atom(C.term_t(t1), C.term_t(t2)))
}

func BvltAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvlt_atom(C.term_t(t1), C.term_t(t2)))
}

func BvsgeAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsge_atom(C.term_t(t1), C.term_t(t2)))
}

func BvsgtAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsgt_atom(C.term_t(t1), C.term_t(t2)))
}

func BvsleAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsle_atom(C.term_t(t1), C.term_t(t2)))
}

func BvsltAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvslt_atom(C.term_t(t1), C.term_t(t2)))
}

/**************
 *  PARSING   *
 *************/

func ParseType(s string) TypeT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TypeT(C.yices_parse_type(cs))
}

func ParseTerm(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_term(cs))
}

/*******************
 *  SUBSTITUTIONS  *
 ******************/

func SubstTerm(vars []TermT, vals []TermT, t TermT) TermT {
	count := C.uint32_t(len(vars))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_subst_term(count, (*C.term_t)(&vars[0]), (*C.term_t)(&vals[0]), C.term_t(t)))
}

func SubstTermArray(vars []TermT, vals []TermT, t []TermT) TermT {
	count := C.uint32_t(len(vars))
	tcount := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 || tcount == 0 {
		return NullTerm
	}
	return TermT(C.yices_subst_term_array(count, (*C.term_t)(&vars[0]), (*C.term_t)(&vals[0]), tcount, (*C.term_t)(&t[0])))
}

/************
 *  NAMES   *
 ***********/

func SetTypeName(tau TypeT, name string) int32 {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return int32(C.yices_set_type_name(C.type_t(tau), cs))
}

func SetTermName(t TermT, name string) int32 {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return int32(C.yices_set_term_name(C.term_t(t), cs))
}

func RemoveTypeName(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.yices_remove_type_name(cs)
}

func RemoveTermName(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.yices_remove_term_name(cs)
}

func GetTypeByName(name string) TypeT {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return TypeT(C.yices_get_type_by_name(cs))
}

func GetTermByName(name string) TermT {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_get_term_by_name(cs))
}

func ClearTypeName(tau TypeT) int32 {
	return int32(C.yices_clear_type_name(C.type_t(tau)))
}

func ClearTermName(t TermT) int32 {
	return int32(C.yices_clear_term_name(C.term_t(t)))
}

func GetTypeName(tau TypeT) string {
	//FIXME: check if the name needs to be freed
	return C.GoString(C.yices_get_type_name(C.type_t(tau)))
}

func GetTermName(t TermT) string {
	//FIXME: check if the name needs to be freed
	return C.GoString(C.yices_get_term_name(C.term_t(t)))
}

/***********************
 *  TERM EXPLORATION   *
 **********************/

func TypeOfTerm(t TermT) TypeT {
	return TypeT(C.yices_type_of_term(C.term_t(t)))
}

func TermIsBool(t TermT) bool {
	return C.yices_term_is_bool(C.term_t(t)) == C.int32_t(1)
}

func TermIsInt(t TermT) bool {
	return C.yices_term_is_int(C.term_t(t)) == C.int32_t(1)
}

func TermIsReal(t TermT) bool {
	return C.yices_term_is_real(C.term_t(t)) == C.int32_t(1)
}

func TermIsArithmetic(t TermT) bool {
	return C.yices_term_is_arithmetic(C.term_t(t)) == C.int32_t(1)
}

func TermIsBitvector(t TermT) bool {
	return C.yices_term_is_bitvector(C.term_t(t)) == C.int32_t(1)
}

func TermIsTuple(t TermT) bool {
	return C.yices_term_is_tuple(C.term_t(t)) == C.int32_t(1)
}

func TermIsFunction(t TermT) bool {
	return C.yices_term_is_function(C.term_t(t)) == C.int32_t(1)
}

func TermIsScalar(t TermT) bool {
	return C.yices_term_is_scalar(C.term_t(t)) == C.int32_t(1)
}

func TermBitsize(t TermT) uint32 {
	return uint32(C.yices_term_bitsize(C.term_t(t)))
}

func TermIsGround(t TermT) bool {
	return C.yices_term_is_ground(C.term_t(t)) == C.int32_t(1)
}

func TermIsAtomic(t TermT) bool {
	return C.yices_term_is_atomic(C.term_t(t)) == C.int32_t(1)
}

func TermIsComposite(t TermT) bool {
	return C.yices_term_is_composite(C.term_t(t)) == C.int32_t(1)
}

func TermIsProjection(t TermT) bool {
	return C.yices_term_is_projection(C.term_t(t)) == C.int32_t(1)
}

func TermIsSum(t TermT) bool {
	return C.yices_term_is_sum(C.term_t(t)) == C.int32_t(1)
}

func TermIsBvsum(t TermT) bool {
	return C.yices_term_is_bvsum(C.term_t(t)) == C.int32_t(1)
}

func TermIsProduct(t TermT) bool {
	return C.yices_term_is_product(C.term_t(t)) == C.int32_t(1)
}

func TermConstructor(t TermT) TermConstructorT {
	return TermConstructorT(C.yices_term_constructor(C.term_t(t)))
}

func TermNumChildren(t TermT) int32 {
	return int32(C.yices_term_num_children(C.term_t(t)))
}

func TermChild(t TermT, i int32) TermT {
	return TermT(C.yices_term_child(C.term_t(t), C.int32_t(i)))
}

// Since 2.6.2
func TermChildren(t TermT) (children []TermT) {
	var tv C.term_vector_t
	C.yices_init_term_vector(&tv)
	ycount := int32(C.yices_term_children(C.type_t(t), &tv))
	if ycount != -1 {
		count := int(tv.size)
		children = make([]TermT, count, count)
		// defined in the preamble yices_term_vector_get(term_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			children[i] = TermT(C.yices_term_vector_get(&tv, C.uint32_t(i)))
		}
	}
	C.yices_delete_term_vector(&tv)
	return
}

func ProjIndex(t TermT) int32 {
	return int32(C.yices_proj_index(C.term_t(t)))
}

func ProjArg(t TermT) TermT {
	return TermT(C.yices_proj_arg(C.term_t(t)))
}

func BoolConstValue(t TermT, val *int32) int32 {
	return int32(C.yices_bool_const_value(C.term_t(t), (*C.int32_t)(val)))
}

func BvConstValue(t TermT, val []int32) int32 {
	return int32(C.yices_bv_const_value(C.term_t(t), (*C.int32_t)(&val[0])))
}

func ScalarConstValue(t TermT, val *int32) int32 {
	return int32(C.yices_scalar_const_value(C.term_t(t), (*C.int32_t)(val)))
}

func RationalConstValue(t TermT, q *MpqT) int32 {
	return int32(C.yices_rational_const_valuep(C.term_t(t), (*C.mpq_t)(unsafe.Pointer(q))))
}

func SumComponent(t TermT, i int32, coeff *MpqT, term *TermT) int32 {
	return int32(C.yices_sum_componentp(C.term_t(t), C.int32_t(i), (*C.mpq_t)(unsafe.Pointer(coeff)), (*C.term_t)(term)))
}

func BvsumComponent(t TermT, i int32, val []int32, term *TermT) int32 {
	return int32(C.yices_bvsum_component(C.term_t(t), C.int32_t(i), (*C.int32_t)(&val[0]), (*C.term_t)(term)))
}

func ProductComponent(t TermT, i int32, term *TermT, exp *uint32) int32 {
	return int32(C.yices_product_component(C.term_t(t), C.int32_t(i), (*C.term_t)(term), (*C.uint32_t)(exp)))
}

/*************************
 *  GARBAGE COLLECTION   *
 ************************/

func NumTerms() uint32 {
	return uint32(C.yices_num_terms())
}

func NumTypes() uint32 {
	return uint32(C.yices_num_types())
}

func IncrefTerm(t TermT) int32 {
	return int32(C.yices_incref_term(C.term_t(t)))
}

func DecrefTerm(t TermT) int32 {
	return int32(C.yices_decref_term(C.term_t(t)))
}

func IncrefType(tau TypeT) int32 {
	return int32(C.yices_incref_type(C.type_t(tau)))
}

func DecrefType(tau TypeT) int32 {
	return int32(C.yices_decref_type(C.type_t(tau)))
}

func NumPosrefTerms() uint32 {
	return uint32(C.yices_num_posref_terms())
}

func NumPosrefTypes() uint32 {
	return uint32(C.yices_num_posref_types())
}

func GarbageCollect(ts []TermT, taus []TypeT, keepNamed int32) {
	tCount := C.uint32_t(len(ts))
	tauCount := C.uint32_t(len(taus))
	C.yices_garbage_collect((*C.term_t)(&ts[0]), tCount, (*C.type_t)(&taus[0]), tauCount, C.int32_t(keepNamed))
}

/****************************
 *  CONTEXT CONFIGURATION   *
 ***************************/

type ConfigT struct {
	raw uintptr // actually *C.ctx_config_t
}

func ycfg(cfg ConfigT) *C.ctx_config_t {
	return (*C.ctx_config_t)(unsafe.Pointer(cfg.raw))
}

func InitConfig(cfg *ConfigT) {
	cfg.raw = uintptr(unsafe.Pointer(C.yices_new_config()))
}

func CloseConfig(cfg *ConfigT) {
	C.yices_free_config(ycfg(*cfg))
	cfg.raw = 0
}

func SetConfig(cfg ConfigT, name string, value string) int32 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	return int32(C.yices_set_config(ycfg(cfg), cname, cvalue))
}

func DefaultConfigForLogic(cfg ConfigT, logic string) int32 {
	clogic := C.CString(logic)
	defer C.free(unsafe.Pointer(clogic))
	return int32(C.yices_default_config_for_logic(ycfg(cfg), clogic))
}

/***************
 *  CONTEXTS   *
 **************/

type ContextT struct {
	raw uintptr // actually *C.context_t
}

func yctx(ctx ContextT) *C.context_t {
	return (*C.context_t)(unsafe.Pointer(ctx.raw))
}

func InitContext(cfg ConfigT, ctx *ContextT) {
	ctx.raw = uintptr(unsafe.Pointer(C.yices_new_context(ycfg(cfg))))
}

func CloseContext(ctx *ContextT) {
	C.yices_free_context(yctx(*ctx))
	ctx.raw = 0
}

func ContextStatus(ctx ContextT) SmtStatusT {
	return SmtStatusT(C.yices_context_status(yctx(ctx)))
}

func ResetContext(ctx ContextT) {
	C.yices_reset_context(yctx(ctx))
}

func Push(ctx ContextT) int32 {
	return int32(C.yices_push(yctx(ctx)))
}

func Pop(ctx ContextT) int32 {
	return int32(C.yices_pop(yctx(ctx)))
}

func ContextEnableOption(ctx ContextT, option string) int32 {
	coption := C.CString(option)
	defer C.free(unsafe.Pointer(coption))
	return int32(C.yices_context_enable_option(yctx(ctx), coption))
}

func ContextDisableOption(ctx ContextT, option string) int32 {
	coption := C.CString(option)
	defer C.free(unsafe.Pointer(coption))
	return int32(C.yices_context_enable_option(yctx(ctx), coption))
}

func AssertFormula(ctx ContextT, t TermT) int32 {
	return int32(C.yices_assert_formula(yctx(ctx), C.term_t(t)))

}

func AssertFormulas(ctx ContextT, t []TermT) int32 {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return -1
	}
	return int32(C.yices_assert_formulas(yctx(ctx), count, (*C.term_t)(&t[0])))
}

func CheckContext(ctx ContextT, params ParamT) SmtStatusT {
	return SmtStatusT(C.yices_check_context(yctx(ctx), yparam(params)))
}

func CheckContextWithAssumptions(ctx ContextT, params ParamT, t []TermT) SmtStatusT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return SmtStatusT(StatusError)
	}
	return SmtStatusT(C.yices_check_context_with_assumptions(yctx(ctx), yparam(params), count, (*C.term_t)(&t[0])))
}

func AssertBlockingClause(ctx ContextT) int32 {
	return int32(C.yices_assert_blocking_clause(yctx(ctx)))

}

func StopSearch(ctx ContextT) {
	C.yices_stop_search(yctx(ctx))
}

/*
 * SEARCH PARAMETERS
 */

type ParamT struct {
	raw uintptr // actually *C.param_t
}

func yparam(params ParamT) *C.param_t {
	return (*C.param_t)(unsafe.Pointer(params.raw))
}

func InitParamRecord(params *ParamT) {
	params.raw = uintptr(unsafe.Pointer(C.yices_new_param_record()))
}

func CloseParamRecord(params *ParamT) {
	C.yices_free_param_record(yparam(*params))
	params.raw = 0
}

func DefaultParamsForContext(ctx ContextT, params ParamT) {
	C.yices_default_params_for_context(yctx(ctx), yparam(params))
}

func SetParam(params ParamT, pname string, value string) int32 {
	cpname := C.CString(pname)
	defer C.free(unsafe.Pointer(cpname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	return int32(C.yices_set_param(yparam(params), cpname, cvalue))
}

/****************
 *  UNSAT CORE  *
 ***************/

func GetUnsatCore(ctx ContextT) (unsatCore []TermT) {
	var tv C.term_vector_t
	C.yices_init_term_vector(&tv)
	errcode := int32(C.yices_get_unsat_core(yctx(ctx), &tv))
	if errcode != -1 {
		count := int(tv.size)
		unsatCore = make([]TermT, count, count)
		// defined in the preamble yices_term_vector_get(term_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			unsatCore[i] = TermT(C.yices_term_vector_get(&tv, C.uint32_t(i)))
		}
	}
	C.yices_delete_term_vector(&tv)
	return
}

/**************
 *   MODELS   *
 *************/

type ModelT struct {
	raw uintptr // actually *C.model_t
}

func ymodel(model ModelT) *C.model_t {
	return (*C.model_t)(unsafe.Pointer(model.raw))
}

func GetModel(ctx ContextT, keepSubst int32) *ModelT {
	//yes golang lets you return stuff allocated on the stack
	return &ModelT{uintptr(unsafe.Pointer(C.yices_get_model(yctx(ctx), C.int32_t(keepSubst))))}
}

func CloseModel(model *ModelT) {
	C.yices_free_model(ymodel(*model))
	model.raw = 0
}

func ModelFromMap(vars []TermT, vals []TermT) *ModelT {
	vcount := C.uint32_t(len(vals))
	return &ModelT{uintptr(unsafe.Pointer(C.yices_model_from_map(vcount, (*C.term_t)(&vars[0]), (*C.term_t)(&vals[0]))))}
}

func ModelCollectDefinedTerms(model ModelT) (terms []TermT) {
	var tv C.term_vector_t
	C.yices_init_term_vector(&tv)
	C.yices_model_collect_defined_terms(ymodel(model), &tv)
	count := int(tv.size)
	terms = make([]TermT, count, count)
	for i := 0; i < count; i++ {
		terms[i] = TermT(C.yices_term_vector_get(&tv, C.uint32_t(i)))
	}
	C.yices_delete_term_vector(&tv)
	return
}

/***********************
 *  VALUES IN A MODEL  *
 **********************/

/*
 * EVALUATION FOR SIMPLE TYPES
 */

func GetBoolValue(model ModelT, t TermT, val *int32) int32 {
	return int32(C.yices_get_bool_value(ymodel(model), C.term_t(t), (*C.int32_t)(val)))
}

func GetInt32Value(model ModelT, t TermT, val *int32) int32 {
	return int32(C.yices_get_int32_value(ymodel(model), C.term_t(t), (*C.int32_t)(val)))
}

func GetInt64Value(model ModelT, t TermT, val *int64) int32 {
	return int32(C.yices_get_int64_value(ymodel(model), C.term_t(t), (*C.int64_t)(val)))
}

func GetRational32Value(model ModelT, t TermT, num *int32, den *uint32) int32 {
	return int32(C.yices_get_rational32_value(ymodel(model), C.term_t(t), (*C.int32_t)(num), (*C.uint32_t)(den)))
}

func GetRational64Value(model ModelT, t TermT, num *int64, den *uint64) int32 {
	return int32(C.yices_get_rational64_value(ymodel(model), C.term_t(t), (*C.int64_t)(num), (*C.uint64_t)(den)))
}

func GetDoubleValue(model ModelT, t TermT, val *float64) int32 {
	return int32(C.yices_get_double_value(ymodel(model), C.term_t(t), (*C.double)(val)))
}

func GetMpzValue(model ModelT, t TermT, val *MpzT) int32 {
	return int32(C.yices_get_mpz_valuep(ymodel(model), C.term_t(t), (*C.mpz_t)(unsafe.Pointer(val))))
}

func GetMpqValue(model ModelT, t TermT, val *MpqT) int32 {
	return int32(C.yices_get_mpq_valuep(ymodel(model), C.term_t(t), (*C.mpq_t)(unsafe.Pointer(val))))
}

/*
//iam: not gonna assume mcsat
#ifdef LIBPOLY_VERSION
__YICES_DLLSPEC__ extern int32_t yices_get_algebraic_number_value(model_t *mdl, term_t t, lp_algebraic_number_t *a);
#endif
*/

func GetBvValue(model ModelT, t TermT, val []int32) int32 {
	return int32(C.yices_get_bv_value(ymodel(model), C.term_t(t), (*C.int32_t)(&val[0])))
}

func GetScalarValue(model ModelT, t TermT, val *int32) int32 {
	return int32(C.yices_get_scalar_value(ymodel(model), C.term_t(t), (*C.int32_t)(val)))
}

/*
 * GENERIC FORM: VALUE DESCRIPTORS AND NODES
 */

type YvalT C.yval_t

func GetTag(yval YvalT) YvalTagT {
	return YvalTagT(yval.node_tag)
}

type YvalVectorT C.yval_vector_t

func InitYvalVector(v *YvalVectorT) {
	C.yices_init_yval_vector((*C.yval_vector_t)(v))
}

func DeleteYvalVector(v *YvalVectorT) {
	C.yices_delete_yval_vector((*C.yval_vector_t)(v))
}

func ResetYvalVector(v *YvalVectorT) {
	C.yices_reset_yval_vector((*C.yval_vector_t)(v))
}

func GetValue(model ModelT, t TermT, val *YvalT) int32 {
	return int32(C.yices_get_value(ymodel(model), C.term_t(t), (*C.yval_t)(val)))
}

func ValIsInt32(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_int32(ymodel(model), (*C.yval_t)(val)))
}

func ValIsInt64(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_int64(ymodel(model), (*C.yval_t)(val)))
}

func ValIsRational32(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_rational32(ymodel(model), (*C.yval_t)(val)))
}

func ValIsRational64(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_rational64(ymodel(model), (*C.yval_t)(val)))
}

func ValIsInteger(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_integer(ymodel(model), (*C.yval_t)(val)))
}

func ValBitsize(model ModelT, val *YvalT) uint32 {
	return uint32(C.yices_val_bitsize(ymodel(model), (*C.yval_t)(val)))
}

func ValTupleArity(model ModelT, val *YvalT) uint32 {
	return uint32(C.yices_val_tuple_arity(ymodel(model), (*C.yval_t)(val)))
}

func ValMappingArity(model ModelT, val *YvalT) uint32 {
	return uint32(C.yices_val_mapping_arity(ymodel(model), (*C.yval_t)(val)))
}

func ValFunctionArity(model ModelT, val *YvalT) uint32 {
	return uint32(C.yices_val_function_arity(ymodel(model), (*C.yval_t)(val)))
}

// Since 2.6.2
func ValFunctionType(model ModelT, val *YvalT) TypeT {
	return TypeT(C.yices_val_function_type(ymodel(model), (*C.yval_t)(val)))
}

func ValGetBool(model ModelT, yval *YvalT, val *int32) int32 {
	return int32(C.yices_val_get_bool(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(val)))
}

func ValGetInt32(model ModelT, yval *YvalT, val *int32) int32 {
	return int32(C.yices_val_get_int32(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(val)))
}

func ValGetInt64(model ModelT, yval *YvalT, val *int64) int32 {
	return int32(C.yices_val_get_int64(ymodel(model), (*C.yval_t)(yval), (*C.int64_t)(val)))
}

func ValGetRational32(model ModelT, yval *YvalT, num *int32, den *uint32) int32 {
	return int32(C.yices_val_get_rational32(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(num), (*C.uint32_t)(den)))
}

func ValGetRational64(model ModelT, yval *YvalT, num *int64, den *uint64) int32 {
	return int32(C.yices_val_get_rational64(ymodel(model), (*C.yval_t)(yval), (*C.int64_t)(num), (*C.uint64_t)(den)))
}

func ValGetDouble(model ModelT, yval *YvalT, val *float64) int32 {
	return int32(C.yices_val_get_double(ymodel(model), (*C.yval_t)(yval), (*C.double)(val)))
}

func ValGetMpz(model ModelT, yval *YvalT, val *MpzT) int32 {
	return int32(C.yices_val_get_mpzp(ymodel(model), (*C.yval_t)(yval), (*C.mpz_t)(unsafe.Pointer(val))))
}

func ValGetMpq(model ModelT, yval *YvalT, val *MpqT) int32 {
	return int32(C.yices_val_get_mpqp(ymodel(model), (*C.yval_t)(yval), (*C.mpq_t)(unsafe.Pointer(val))))
}

/*
//iam: not gonna assume mcsat
#ifdef LIBPOLY_VERSION
__YICES_DLLSPEC__ extern int32_t yices_val_get_algebraic_number(model_t *mdl, const yval_t *v, lp_algebraic_number_t *a);
#endif
*/

func ValGetBv(model ModelT, yval *YvalT, val []int32) int32 {
	return int32(C.yices_val_get_bv(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(&val[0])))
}

func ValGetScalar(model ModelT, yval *YvalT, val *int32, tau *TypeT) int32 {
	return int32(C.yices_val_get_scalar(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(val), (*C.type_t)(tau)))
}

func ValExpandTuple(model ModelT, yval *YvalT, child []YvalT) int32 {
	return int32(C.yices_val_expand_tuple(ymodel(model), (*C.yval_t)(yval), (*C.yval_t)(&child[0])))
}

func ValExpandFunction(model ModelT, yval *YvalT, def *YvalT) (vector []YvalT) {
	var tv C.yval_vector_t
	C.yices_init_yval_vector(&tv)
	errcode := int32(C.yices_val_expand_function(ymodel(model), (*C.yval_t)(yval), (*C.yval_t)(def), (*C.yval_vector_t)(&tv)))
	if errcode != -1 {
		count := int(tv.size)
		vector = make([]YvalT, count, count)
		// defined in the preamble yices_term_vector_get(term_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			var yv C.yval_t
			C.yices_yval_vector_get(&tv, C.uint32_t(i), (*C.yval_t)(&yv))
			vector[i] = YvalT(yv)
		}
	}
	C.yices_delete_yval_vector(&tv)
	return
}

func ValExpandMapping(model ModelT, m *YvalT, val *YvalT) (vector []YvalT) {
	arity := int(C.yices_val_mapping_arity(ymodel(model), (*C.yval_t)(m)))
	if arity > 0 {
		vec := make([]YvalT, arity, arity)
		errcode := int32(C.yices_val_expand_mapping(ymodel(model), (*C.yval_t)(m), (*C.yval_t)(&vec[0]), (*C.yval_t)(val)))
		if errcode != -1 {
			vector = vec
		}
	}
	return
}

func FormulaTrueInModel(model ModelT, t TermT) int32 {
	return int32(C.yices_formula_true_in_model(ymodel(model), C.term_t(t)))
}

func FormulasTrueInModel(model ModelT, t []TermT) int32 {
	tcount := C.uint32_t(len(t))
	return int32(C.yices_formulas_true_in_model(ymodel(model), tcount, (*C.term_t)(&t[0])))
}

/*
 * CONVERSION OF VALUES TO CONSTANT TERMS
 */

func GetValueAsTerm(model ModelT, t TermT) TermT {
	return TermT(C.yices_get_value_as_term(ymodel(model), C.term_t(t)))
}

func TermArrayValue(model ModelT, a []TermT, b []TermT) int32 {
	tcount := C.uint32_t(len(a))
	return int32(C.yices_term_array_value(ymodel(model), tcount, (*C.term_t)(&a[0]), (*C.term_t)(&b[0])))
}

/*
 * IMPLICANTS
 */

func ImplicantForFormula(model ModelT, t TermT) (literals []TermT) {
	var tv C.term_vector_t
	C.yices_init_term_vector(&tv)
	errcode := int32(C.yices_implicant_for_formula(ymodel(model), C.term_t(t), (*C.term_vector_t)(&tv)))
	if errcode != -1 {
		count := int(tv.size)
		literals = make([]TermT, count, count)
		// defined in the preamble yices_term_vector_get(term_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			literals[i] = TermT(C.yices_term_vector_get(&tv, C.uint32_t(i)))
		}
	}
	C.yices_delete_term_vector(&tv)
	return
}

func ImplicantForFormulas(model ModelT, t []TermT) (literals []TermT) {
	var tv C.term_vector_t
	C.yices_init_term_vector(&tv)
	tcount := C.uint32_t(len(t))
	errcode := int32(C.yices_implicant_for_formulas(ymodel(model), tcount, (*C.term_t)(&t[0]), (*C.term_vector_t)(&tv)))
	if errcode != -1 {
		count := int(tv.size)
		literals = make([]TermT, count, count)
		// defined in the preamble yices_term_vector_get(term_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			literals[i] = TermT(C.yices_term_vector_get(&tv, C.uint32_t(i)))
		}
	}
	C.yices_delete_term_vector(&tv)
	return
}

/************************
 * MODEL GENERALIZATION *
 ************************/

type GenModeT int32

const (
	GenDefault GenModeT = iota
	GenBySubst
	GenByProj
)

func GeneralizeModel(model ModelT, t TermT, elims []TermT, mode GenModeT) (formulas []TermT) {
	var tv C.term_vector_t
	C.yices_init_term_vector(&tv)
	ecount := C.uint32_t(len(elims))
	errcode := int32(C.yices_generalize_model(ymodel(model), C.term_t(t), ecount, (*C.term_t)(&elims[0]), C.yices_gen_mode_t(mode), (*C.term_vector_t)(&tv)))
	if errcode != -1 {
		count := int(tv.size)
		formulas = make([]TermT, count, count)
		// defined in the preamble yices_term_vector_get(term_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			formulas[i] = TermT(C.yices_term_vector_get(&tv, C.uint32_t(i)))
		}
	}
	C.yices_delete_term_vector(&tv)
	return
}

func GeneralizeModelArray(model ModelT, a []TermT, elims []TermT, mode GenModeT) (formulas []TermT) {
	var tv C.term_vector_t
	C.yices_init_term_vector(&tv)
	acount := C.uint32_t(len(a))
	ecount := C.uint32_t(len(elims))
	errcode := int32(C.yices_generalize_model_array(ymodel(model), acount, (*C.term_t)(&a[0]), ecount, (*C.term_t)(&elims[0]), C.yices_gen_mode_t(mode), (*C.term_vector_t)(&tv)))
	if errcode != -1 {
		count := int(tv.size)
		formulas = make([]TermT, count, count)
		// defined in the preamble yices_term_vector_get(term_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			formulas[i] = TermT(C.yices_term_vector_get(&tv, C.uint32_t(i)))
		}
	}
	C.yices_delete_term_vector(&tv)
	return
}

/*******************
 *   Model Support *
 *******************/

func ModelTermSupport(model ModelT, t TermT) (support []TermT) {
	var tv C.term_vector_t
	C.yices_init_term_vector(&tv)
	errcode := int32(C.yices_model_term_support(ymodel(model), C.term_t(t), (*C.term_vector_t)(&tv)))
	if errcode != -1 {
		count := int(tv.size)
		support = make([]TermT, count, count)
		// defined in the preamble yices_term_vector_get(term_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			support[i] = TermT(C.yices_term_vector_get(&tv, C.uint32_t(i)))
		}
	}
	C.yices_delete_term_vector(&tv)
	return
}

func ModelTermArraySupport(model ModelT, t []TermT) (support []TermT) {
	var tv C.term_vector_t
	C.yices_init_term_vector(&tv)
	tcount := C.uint32_t(len(t))
	errcode := int32(C.yices_model_term_array_support(ymodel(model), tcount, (*C.term_t)(&t[0]), (*C.term_vector_t)(&tv)))
	if errcode != -1 {
		count := int(tv.size)
		support = make([]TermT, count, count)
		// defined in the preamble yices_term_vector_get(term_vector_t* vec, uint32_t elem)
		for i := 0; i < count; i++ {
			support[i] = TermT(C.yices_term_vector_get(&tv, C.uint32_t(i)))
		}
	}
	C.yices_delete_term_vector(&tv)
	return
}

/***************
 *   DELEGATES *
 **************/

func HasDelegate(delegate string) bool {
	cdelegate := C.CString(delegate)
	defer C.free(unsafe.Pointer(cdelegate))
	has := C.yices_has_delegate(cdelegate)
	return has == 1
}

func CheckFormula(t TermT, logic string, delegate string, model *ModelT) (status SmtStatusT) {
	clogic := C.CString(logic)
	defer C.free(unsafe.Pointer(clogic))
	cdelegate := C.CString(delegate)
	defer C.free(unsafe.Pointer(cdelegate))
	var cstatus C.smt_status_t
	var cmodel *C.model_t = nil
	if model != nil {
		var cmodel *C.model_t
		cstatus = C.yices_check_formula(C.term_t(t), clogic, &cmodel, cdelegate)
	} else {
		cstatus = C.yices_check_formula(C.term_t(t), clogic, (**C.model_t)(C.NULL), cdelegate)
	}
	if cstatus == C.STATUS_SAT {
		status = SmtStatusT(cstatus)
		if model != nil {
			*model = ModelT{uintptr(unsafe.Pointer(cmodel))}
		}
	}
	return
}

func CheckFormulas(t []TermT, logic string, delegate string, model *ModelT) (status SmtStatusT) {
	count := C.uint32_t(len(t))
	clogic := C.CString(logic)
	defer C.free(unsafe.Pointer(clogic))
	cdelegate := C.CString(delegate)
	defer C.free(unsafe.Pointer(cdelegate))
	var cstatus C.smt_status_t
	var cmodel *C.model_t = nil
	if model != nil {
		var cmodel *C.model_t
		cstatus = C.yices_check_formulas((*C.term_t)(&t[0]), count, clogic, &cmodel, cdelegate)
	} else {
		cstatus = C.yices_check_formulas((*C.term_t)(&t[0]), count, clogic, (**C.model_t)(C.NULL), cdelegate)
	}
	if cstatus == C.STATUS_SAT {
		status = SmtStatusT(cstatus)
		if model != nil {
			*model = ModelT{uintptr(unsafe.Pointer(cmodel))}
		}
	}
	return
}

/*************
 *   DIMACS  *
 *************/

func ExportFormulaToDimacs(t TermT, filename string, simplify bool, status *SmtStatusT) (errcode int32) {
	path := C.CString(filename)
	defer C.free(unsafe.Pointer(path))
	var csimplify C.int
	if simplify {
		csimplify = 1
	} else {
		csimplify = 0
	}
	var cstatus C.smt_status_t = 0
	errcode = int32(C.yices_export_formula_to_dimacs(C.term_t(t), path, csimplify, &cstatus))
	if errcode == 0 {
		*status = SmtStatusT(cstatus)
	}
	return
}

func ExportFormulasToDimacs(t []TermT, filename string, simplify bool, status *SmtStatusT) (errcode int32) {
	path := C.CString(filename)
	defer C.free(unsafe.Pointer(path))
	count := C.uint32_t(len(t))
	var csimplify C.int
	if simplify {
		csimplify = 1
	} else {
		csimplify = 0
	}
	var cstatus C.smt_status_t = 0
	errcode = int32(C.yices_export_formulas_to_dimacs((*C.term_t)(&t[0]), count, path, csimplify, &cstatus))
	if errcode == 0 {
		*status = SmtStatusT(cstatus)
	}
	return
}

/**********************
 *  PRETTY PRINTING   *
 **********************/

func PpType(file *os.File, tau TypeT, width uint32, height uint32, offset uint32) int32 {
	return int32(C.yices_pp_type_fd(C.int(file.Fd()), C.type_t(tau), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset)))
}

func PpTerm(file *os.File, t TermT, width uint32, height uint32, offset uint32) int32 {
	return int32(C.yices_pp_term_fd(C.int(file.Fd()), C.term_t(t), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset)))
}

func PpTermArray(file *os.File, t []TermT, width uint32, height uint32, offset uint32, horiz int32) int32 {
	tcount := C.uint32_t(len(t))
	return int32(C.yices_pp_term_array_fd(C.int(file.Fd()), tcount, (*C.term_t)(&t[0]), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset), C.int32_t(horiz)))
}

func PrintModel(file *os.File, model ModelT) int32 {
	return int32(C.yices_print_model_fd(C.int(file.Fd()), ymodel(model)))
}

// Since 2.6.2
func PrintTermValues(file *os.File, model ModelT, t []TermT) int32 {
	tcount := C.uint32_t(len(t))
	return int32(C.yices_print_term_values_fd(C.int(file.Fd()), ymodel(model), tcount, (*C.term_t)(&t[0])))
}

// Since 2.6.2
func PpTermValues(file *os.File, model ModelT, t []TermT, width uint32, height uint32, offset uint32) int32 {
	tcount := C.uint32_t(len(t))
	return int32(C.yices_pp_term_values_fd(C.int(file.Fd()), ymodel(model), tcount, (*C.term_t)(&t[0]), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset)))
}

func PpModel(file *os.File, model ModelT, width uint32, height uint32, offset uint32) int32 {
	return int32(C.yices_pp_model_fd(C.int(file.Fd()), ymodel(model), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset)))
}

// TypeToString returns a pretty printed string representation of the given type.
func TypeToString(tau TypeT, width uint32, height uint32, offset uint32) string {
	cstr := C.yices_type_to_string(C.type_t(tau), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset))
	defer C.yices_free_string(cstr)
	return C.GoString(cstr)
}

// TermToString returns a pretty printed string representation of the given term.
func TermToString(t TermT, width uint32, height uint32, offset uint32) string {
	cstr := C.yices_term_to_string(C.term_t(t), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset))
	defer C.yices_free_string(cstr)
	return C.GoString(cstr)
}

// ModelToString returns a pretty printed string representation of the given model.
func ModelToString(model ModelT, width uint32, height uint32, offset uint32) string {
	cstr := C.yices_model_to_string(ymodel(model), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset))
	defer C.yices_free_string(cstr)
	return C.GoString(cstr)
}
