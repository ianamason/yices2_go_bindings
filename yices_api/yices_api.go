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
//iam: passing a mpq_t[] on linux thru cgo seems too hard
term_t yices_poly_mpzp(uint32_t n, mpz_t *z, const term_t t[]){
    return yices_poly_mpz(n, z, t);
}
//iam: passing a mpq_t[] on linux thru cgo seems too hard
term_t yices_poly_mpqp(uint32_t n, mpq_t *q, const term_t t[]){
    return yices_poly_mpq(n, q, t);
}
//iam: passing a  mpz_t thru cgo seems too hard
int32_t yices_model_set_mpzp(model_t *model, term_t var, mpz_t *val){
   return yices_model_set_mpz(model, var, *val);
}
//iam: passing a  mpq_t thru cgo seems too hard
int32_t yices_model_set_mpqp(model_t *model, term_t var, mpq_t *val){
   return yices_model_set_mpq(model, var, *val);
}
//iam: passing a  mpz_t thru cgo seems too hard
int32_t yices_model_set_bv_mpzp(model_t *model, term_t var, mpz_t *val){
   return yices_model_set_bv_mpz(model, var, *val);
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
// iam: have to deal with yices 2.6.4's problem child.
// we can keep the interpolation_context_t hidden from the go user using the following primitives
// and go's ability to return multiple values.
interpolation_context_t* new_interpolation_context(context_t *ctx_A, context_t *ctx_B) {
   interpolation_context_t* ictx = (interpolation_context_t*)calloc(1, sizeof(interpolation_context_t));
   ictx->ctx_A = ctx_A;
   ictx->ctx_B = ctx_B;
   return ictx;
}
term_t get_interpolation_context_interpolant(interpolation_context_t *ictx){
   return ictx->interpolant;
}
model_t* get_interpolation_context_model(interpolation_context_t *ictx){
   return ictx->model;
}
void free_interpolation_context(interpolation_context_t *ictx) {
   free(ictx);
}
*/
import "C"

import "os"
import "fmt"
import "unsafe"

/*
 *  See yices.h for detailed comments.
 *
 *  Naming convention:   yices_foo  becomes yices_api.Foo  (will probably ditch the 2 at some stage)
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

// ErrorString returns a copy of the value of yices_error_string().
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

// TypeT is the analog of term_t defined in yices_types.h. Unclear if this is a good idea.
type TypeT int32

// NullType is the go version of NULL_TYPE.
const NullType TypeT = -1

// BoolType is the go version of yices_bool_type.
func BoolType() TypeT {
	return TypeT(C.yices_bool_type())
}

// IntType is the go version of yices_int_type.
func IntType() TypeT {
	return TypeT(C.yices_int_type())
}

// RealType is the go version of yices_real_type.
func RealType() TypeT {
	return TypeT(C.yices_real_type())
}

// BvType is the go version of yices_bv_type.
func BvType(size uint32) TypeT {
	return TypeT(C.yices_bv_type(C.uint32_t(size)))
}

// NewScalarType is the go version of yices_new_scalar_type.
func NewScalarType(card uint32) TypeT {
	return TypeT(C.yices_new_scalar_type(C.uint32_t(card)))
}

// NewUninterpretedType is the go version of yices_new_uninterpreted_type.
func NewUninterpretedType() TypeT {
	return TypeT(C.yices_new_uninterpreted_type())
}

// TupleType is the go version of yices_tuple_type.
func TupleType(tau []TypeT) TypeT {
	tauLen := len(tau)
	//iam: FIXME need to unify the yices errors and the go errors...
	if tauLen == 0 {
		return NullType
	}
	return TypeT(C.yices_tuple_type(C.uint32_t(tauLen), (*C.type_t)(&tau[0])))
}

// TupleType1 is the unary go version of yices_tuple_type.
func TupleType1(tau1 TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1)}
	return TypeT(C.yices_tuple_type(C.uint32_t(1), (*C.type_t)(&carr[0])))
}

// TupleType2 is the binary go version of yices_tuple_type.
func TupleType2(tau1 TypeT, tau2 TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2)}
	return TypeT(C.yices_tuple_type(C.uint32_t(2), (*C.type_t)(&carr[0])))
}

// TupleType3 is the ternary go version of yices_tuple_type.
func TupleType3(tau1 TypeT, tau2 TypeT, tau3 TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2), C.type_t(tau3)}
	return TypeT(C.yices_tuple_type(C.uint32_t(3), (*C.type_t)(&carr[0])))
}

// FunctionType is the go version of yices_function_type.
func FunctionType(dom []TypeT, rng TypeT) TypeT {
	domLen := len(dom)
	//iam: FIXME need to unify the yices errors and the go errors...
	if domLen == 0 {
		return NullType
	}
	return TypeT(C.yices_function_type(C.uint32_t(domLen), (*C.type_t)(&dom[0]), C.type_t(rng)))
}

// FunctionType1 is the unary go version of yices_function_type.
func FunctionType1(tau1 TypeT, rng TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1)}
	return TypeT(C.yices_function_type(C.uint32_t(1), (*C.type_t)(&carr[0]), C.type_t(rng)))
}

// FunctionType2 is the binary go version of yices_function_type.
func FunctionType2(tau1 TypeT, tau2 TypeT, rng TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2)}
	return TypeT(C.yices_function_type(C.uint32_t(2), (*C.type_t)(&carr[0]), C.type_t(rng)))
}

// FunctionType3 is the ternaery go version of yices_function_type.
func FunctionType3(tau1 TypeT, tau2 TypeT, tau3 TypeT, rng TypeT) TypeT {
	carr := []C.type_t{C.type_t(tau1), C.type_t(tau2), C.type_t(tau3)}
	return TypeT(C.yices_function_type(C.uint32_t(3), (*C.type_t)(&carr[0]), C.type_t(rng)))
}

/*************************
 *   TYPE EXPLORATION    *
 ************************/

// TypeIsBool is the go version of yices_type_is_bool.
func TypeIsBool(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_bool(C.type_t(tau)))
}

// TypeIsInt is the go version of yices_type_is_int.
func TypeIsInt(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_int(C.type_t(tau)))
}

// TypeIsReal is the go version of yices_type_is_real.
func TypeIsReal(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_real(C.type_t(tau)))
}

// TypeIsArithmetic is the go version of yices_type_is_arithmetic.
func TypeIsArithmetic(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_arithmetic(C.type_t(tau)))
}

// TypeIsBitvector is the go version of yices_type_is_bitvector.
func TypeIsBitvector(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_bitvector(C.type_t(tau)))
}

// TypeIsTuple is the go version of yices_type_is_tuple.
func TypeIsTuple(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_tuple(C.type_t(tau)))
}

// TypeIsFunction is the go version of yices_type_is_function.
func TypeIsFunction(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_function(C.type_t(tau)))
}

// TypeIsScalar is the go version of yices_type_is_scalar.
func TypeIsScalar(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_scalar(C.type_t(tau)))
}

// TypeIsUninterpreted is the go version of yices_type_is_uninterpreted.
func TypeIsUninterpreted(tau TypeT) bool {
	return int32(1) == int32(C.yices_type_is_uninterpreted(C.type_t(tau)))
}

// TestSubtype is the go version of yices_test_subtype.
func TestSubtype(tau TypeT, sigma TypeT) bool {
	return int32(1) == int32(C.yices_test_subtype(C.type_t(tau), C.type_t(sigma)))
}

// CompatibleTypes is the go version of yices_compatible_types.
func CompatibleTypes(tau TypeT, sigma TypeT) bool {
	return int32(1) == int32(C.yices_compatible_types(C.type_t(tau), C.type_t(sigma)))
}

// BvtypeSize is the go version of yices_bvtype_size.
func BvtypeSize(tau TypeT) uint32 {
	return uint32(C.yices_bvtype_size(C.type_t(tau)))
}

// ScalarTypeCard is the go version of yices_scalar_type_card.
func ScalarTypeCard(tau TypeT) uint32 {
	return uint32(C.yices_scalar_type_card(C.type_t(tau)))
}

// TypeNumChildren is the go version of yices_type_num_children.
func TypeNumChildren(tau TypeT) int32 {
	return int32(C.yices_type_num_children(C.type_t(tau)))
}

// TypeChild is the go version of yices_type_child.
func TypeChild(tau TypeT, i int32) TypeT {
	return TypeT(C.yices_type_child(C.type_t(tau), C.int32_t(i)))
}

// TypeChildren is the go version of yices_type_children.
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

// TermT is the go analog on term_t defined in yices_types.h. Unclear if this is a wise move.
type TermT int32

// NullTerm is the go analog of NULL_TERM defined in yices_types.h.
const NullTerm TermT = -1

// True is the go version of yices_true.
func True() TermT {
	return TermT(C.yices_true())
}

// False is the go version of yices_false.
func False() TermT {
	return TermT(C.yices_false())
}

// Constant is the go version of yices_constant.
func Constant(tau TypeT, index int32) TermT {
	return TermT(C.yices_constant(C.type_t(tau), C.int32_t(index)))
}

// NewUninterpretedTerm is the go version of yices_new_uninterpreted_term.
func NewUninterpretedTerm(tau TypeT) TermT {
	return TermT(C.yices_new_uninterpreted_term(C.type_t(tau)))
}

// NewVariable is the go version of yices_new_variable.
func NewVariable(tau TypeT) TermT {
	return TermT(C.yices_new_variable(C.type_t(tau)))
}

// Application is the go version of yices_application.
func Application(fun TermT, argv []TermT) TermT {
	argc := len(argv)
	//iam: FIXME need to unify the yices errors and the go errors...
	if argc == 0 {
		return NullTerm
	}
	return TermT(C.yices_application(C.term_t(fun), C.uint32_t(argc), (*C.term_t)(&argv[0])))
}

// Application1 is the unary go version of yices_application.
func Application1(fun TermT, arg1 TermT) TermT {
	argv := []C.term_t{C.term_t(arg1)}
	return TermT(C.yices_application(C.term_t(fun), C.uint32_t(1), (*C.term_t)(&argv[0])))
}

// Application2 is the binary go version of yices_application.
func Application2(fun TermT, arg1 TermT, arg2 TermT) TermT {
	argv := []C.term_t{C.term_t(arg1), C.term_t(arg2)}
	return TermT(C.yices_application(C.term_t(fun), C.uint32_t(2), (*C.term_t)(&argv[0])))
}

// Application3 is the ternary go version of yices_application.
func Application3(fun TermT, arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	argv := []C.term_t{C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)}
	return TermT(C.yices_application(C.term_t(fun), C.uint32_t(3), (*C.term_t)(&argv[0])))
}

// Ite is the go version of yices_ite.
func Ite(cond TermT, thenTerm TermT, elseTerm TermT) TermT {
	return TermT(C.yices_ite(C.term_t(cond), C.term_t(thenTerm), C.term_t(elseTerm)))
}

// Eq is the go version of yices_eq.
func Eq(lhs TermT, rhs TermT) TermT {
	return TermT(C.yices_eq(C.term_t(lhs), (C.term_t(rhs))))
}

// Neq is the go version of yices_neq.
func Neq(lhs TermT, rhs TermT) TermT {
	return TermT(C.yices_neq(C.term_t(lhs), (C.term_t(rhs))))
}

// Not is the go version of yices_not.
func Not(arg TermT) TermT {
	return TermT(C.yices_not(C.term_t(arg)))
}

// Or is the go version of yices_or.
func Or(disjuncts []TermT) TermT {
	count := C.uint32_t(len(disjuncts))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return TermT(C.yices_false())
	}
	return TermT(C.yices_or(count, (*C.term_t)(&disjuncts[0])))
}

// And is the go version of yices_and.
func And(conjuncts []TermT) TermT {
	count := C.uint32_t(len(conjuncts))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return TermT(C.yices_true())
	}
	return TermT(C.yices_and(count, (*C.term_t)(&conjuncts[0])))
}

// Xor is the go version of yices_xor.
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

// Or2 is the go version of yices_or2.
func Or2(arg1 TermT, arg2 TermT) TermT {
	return TermT(C.yices_or2(C.term_t(arg1), C.term_t(arg2)))
}

// And2 is the go version of yices_and2.
func And2(arg1 TermT, arg2 TermT) TermT {
	return TermT(C.yices_and2(C.term_t(arg1), C.term_t(arg2)))
}

// Xor2 is the go version of yices_xor2.
func Xor2(arg1 TermT, arg2 TermT) TermT {
	return TermT(C.yices_xor2(C.term_t(arg1), C.term_t(arg2)))
}

// Or3 is the go version of yices_or3.
func Or3(arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	return TermT(C.yices_or3(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

// And3 is the go version of yices_and3.
func And3(arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	return TermT(C.yices_and3(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

// Xor3 is the go version of yices_xor3.
func Xor3(arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	return TermT(C.yices_xor3(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

// Iff is the go version of yices_iff.
func Iff(lhs TermT, rhs TermT) TermT {
	return TermT(C.yices_iff(C.term_t(lhs), (C.term_t(rhs))))
}

// Implies is the go version of yices_implies.
func Implies(lhs TermT, rhs TermT) TermT {
	return TermT(C.yices_implies(C.term_t(lhs), (C.term_t(rhs))))
}

// Tuple is the go version of yices_tuple.
func Tuple(argv []TermT) TermT {
	count := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_tuple(count, (*C.term_t)(&argv[0])))
}

// Pair is the go version of yices_pair.
func Pair(arg1 TermT, arg2 TermT) TermT {
	return TermT(C.yices_pair(C.term_t(arg1), C.term_t(arg2)))
}

// Triple is the go version of yices_triple.
func Triple(arg1 TermT, arg2 TermT, arg3 TermT) TermT {
	return TermT(C.yices_triple(C.term_t(arg1), C.term_t(arg2), C.term_t(arg3)))
}

// Select is the go version of yices_select.
func Select(index uint32, tuple TermT) TermT {
	return TermT(C.yices_select(C.uint32_t(index), C.term_t(tuple)))
}

// TupleUpdate is the go version of yices_tuple_update.
func TupleUpdate(tuple TermT, index uint32, value TermT) TermT {
	return TermT(C.yices_tuple_update(C.term_t(tuple), C.uint32_t(index), C.term_t(value)))
}

// Update is the go version of yices_update.
func Update(fun TermT, argv []TermT, value TermT) TermT {
	count := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_update(C.term_t(fun), count, (*C.term_t)(&argv[0]), C.term_t(value)))
}

// Update1 is the go version of yices_update1.
func Update1(fun TermT, arg1 TermT, value TermT) TermT {
	return TermT(C.yices_update1(C.term_t(fun), C.term_t(arg1), C.term_t(value)))
}

// Update2 is the go version of yices_update2.
func Update2(fun TermT, arg1 TermT, arg2 TermT, value TermT) TermT {
	return TermT(C.yices_update2(C.term_t(fun), C.term_t(arg1), C.term_t(arg2), C.term_t(value)))
}

// Update3 is the go version of yices_update3.
func Update3(fun TermT, arg1 TermT, arg2 TermT, arg3 TermT, value TermT) TermT {
	return TermT(C.yices_update3(C.term_t(fun), C.term_t(arg1), C.term_t(arg2), C.term_t(arg3), C.term_t(value)))
}

// Distinct is the go version of yices_distinct.
func Distinct(argv []TermT) TermT {
	n := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if n == 0 {
		return NullTerm
	}
	return TermT(C.yices_distinct(n, (*C.term_t)(&argv[0])))
}

// Forall is the go version of yices_forall.
func Forall(vars []TermT, body TermT) TermT {
	n := C.uint32_t(len(vars))
	//iam: FIXME need to unify the yices errors and the go errors...
	if n == 0 {
		return NullTerm
	}
	return TermT(C.yices_forall(n, (*C.term_t)(&vars[0]), C.term_t(body)))
}

// Exists is the go version of yices_exists.
func Exists(vars []TermT, body TermT) TermT {
	n := C.uint32_t(len(vars))
	//iam: FIXME need to unify the yices errors and the go errors...
	if n == 0 {
		return NullTerm
	}
	return TermT(C.yices_exists(n, (*C.term_t)(&vars[0]), C.term_t(body)))
}

// Lambda is the go version of yices_lambda.
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

// Zero is the go version of yices_zero.
func Zero() TermT {
	return TermT(C.yices_zero())
}

// Int32 is the go version of yices_int32.
func Int32(val int32) TermT {
	return TermT(C.yices_int32(C.int32_t(val)))
}

// Int64 is the go version of yices_int64.
func Int64(val int64) TermT {
	return TermT(C.yices_int64(C.int64_t(val)))
}

// Rational32 is the go version of yices_rational32.
func Rational32(num int32, den uint32) TermT {
	return TermT(C.yices_rational32(C.int32_t(num), C.uint32_t(den)))
}

// Rational64 is the go version of yices_rational64.
func Rational64(num int64, den uint64) TermT {
	return TermT(C.yices_rational64(C.int64_t(num), C.uint64_t(den)))
}

// MpzT is defined so that we can refer to C.mpz_t outside of this package
type MpzT C.mpz_t

// InitMpz creates a MpzT.
func InitMpz(mpz *MpzT) {
	C.init_mpzp(C.uintptr_t(uintptr(unsafe.Pointer(mpz))))
}

// CloseMpz frees MpzT pointer.
func CloseMpz(mpz *MpzT) {
	C.close_mpzp(C.uintptr_t(uintptr(unsafe.Pointer(mpz))))
}

// Mpz is the go version of yices_mpz.
func Mpz(z *MpzT) TermT {
	// some contortions needed here to do the simplest of things
	// similar contortions are needed to "identify" yices2._Ctype_mpz_t
	// with gmp._Ctype_mpz_t, which seems a little odd.
	return TermT(C.ympz(C.uintptr_t(uintptr(unsafe.Pointer(z)))))
}

// MpqT is defined so that we can refer to C.mpq_t outside of this package
type MpqT C.mpq_t

// InitMpq creates a MpqT.
func InitMpq(mpq *MpqT) {
	C.init_mpqp(C.uintptr_t(uintptr(unsafe.Pointer(mpq))))
}

// CloseMpq frees up a MpqT pointer.
func CloseMpq(mpq *MpqT) {
	C.close_mpqp(C.uintptr_t(uintptr(unsafe.Pointer(mpq))))
}

// Mpq is the go version of yices_mpq.
func Mpq(q *MpqT) TermT {
	// some contortions needed here to do the simplest of things
	return TermT(C.ympq(C.uintptr_t(uintptr(unsafe.Pointer(q)))))
}

// ParseRational is the go version of yices_parse_rational.
func ParseRational(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_rational(cs))
}

// ParseFloat is the go version of yices_parse_float.
func ParseFloat(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_float(cs))
}

/*
 * ARITHMETIC OPERATIONS
 */

// Add is the go version of yices_add.
func Add(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_add(C.term_t(t1), C.term_t(t2)))
}

// Sub is the go version of yices_sub.
func Sub(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_sub(C.term_t(t1), C.term_t(t2)))
}

// Neg is the go version of yices_neg.
func Neg(t1 TermT) TermT {
	return TermT(C.yices_neg(C.term_t(t1)))
}

// Mul is the go version of yices_mul.
func Mul(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_mul(C.term_t(t1), C.term_t(t2)))
}

// Square is the go version of yices_square.
func Square(t1 TermT) TermT {
	return TermT(C.yices_square(C.term_t(t1)))
}

// Power is the go version of yices_power.
func Power(t1 TermT, d uint32) TermT {
	return TermT(C.yices_power(C.term_t(t1), C.uint32_t(d)))
}

// Sum is the go version of yices_sum.
func Sum(argv []TermT) TermT {
	count := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_sum(count, (*C.term_t)(&argv[0])))
}

// Product is the go version of yices_product.
func Product(argv []TermT) TermT {
	count := C.uint32_t(len(argv))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return TermT(C.yices_int32(1))
	}
	return TermT(C.yices_product(count, (*C.term_t)(&argv[0])))
}

// Division is the go version of yices_division.
func Division(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_division(C.term_t(t1), C.term_t(t2)))
}

// Idiv is the go version of yices_idiv.
func Idiv(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_idiv(C.term_t(t1), C.term_t(t2)))
}

// Imod is the go version of yices_imod.
func Imod(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_imod(C.term_t(t1), C.term_t(t2)))
}

// DividesAtom is the go version of yices_divides_atom.
func DividesAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_divides_atom(C.term_t(t1), C.term_t(t2)))
}

// IsIntAtom is the go version of yices_is_int_atom.
func IsIntAtom(t TermT) TermT {
	return TermT(C.yices_is_int_atom(C.term_t(t)))
}

// Abs is the go version of yices_abs.
func Abs(t1 TermT) TermT {
	return TermT(C.yices_abs(C.term_t(t1)))
}

// Floor is the go version of yices_floor.
func Floor(t1 TermT) TermT {
	return TermT(C.yices_floor(C.term_t(t1)))
}

// Ceil is the go version of yices_ceil.
func Ceil(t1 TermT) TermT {
	return TermT(C.yices_ceil(C.term_t(t1)))
}

/*
 * POLYNOMIALS
 */

// PolyInt32 is the go version of yices_poly_int32.
func PolyInt32(a []int32, t []TermT) TermT {
	count := C.uint32_t(len(a))
	//iam: FIXME need to unify the yices errors and the go errors...
	// do we want to be nannies here?
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_int32(count, (*C.int32_t)(&a[0]), (*C.term_t)(&t[0])))
}

// PolyInt64 is the go version of yices_poly_int64.
func PolyInt64(a []int64, t []TermT) TermT {
	count := C.uint32_t(len(a))
	//iam: FIXME need to unify the yices errors and the go errors...
	// do we want to be nannies here?
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_int64(count, (*C.int64_t)(&a[0]), (*C.term_t)(&t[0])))
}

// PolyRational32 is the go version of yices_poly_rational32.
func PolyRational32(num []int32, den []uint32, t []TermT) TermT {
	count := C.uint32_t(len(num))
	//iam: FIXME need to unify the yices errors and the go errors...
	// do we want to be nannies here?
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_rational32(count, (*C.int32_t)(&num[0]), (*C.uint32_t)(&den[0]), (*C.term_t)(&t[0])))
}

// PolyRational64 is the go version of yices_poly_rational64.
func PolyRational64(num []int64, den []uint64, t []TermT) TermT {
	count := C.uint32_t(len(num))
	//iam: FIXME need to unify the yices errors and the go errors...
	// do we want to be nannies here?
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_rational64(count, (*C.int64_t)(&num[0]), (*C.uint64_t)(&den[0]), (*C.term_t)(&t[0])))
}

// PolyMpz is the go version of yices_poly_mpz.
func PolyMpz(z []MpzT, t []TermT) TermT {
	count := C.uint32_t(len(z))
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_mpzp(count, (*C.mpz_t)(&z[0]), (*C.term_t)(&t[0])))
}

// PolyMpq is the go version of yices_poly_mpq.
func PolyMpq(q []MpqT, t []TermT) TermT {
	count := C.uint32_t(len(q))
	if count == 0 {
		return TermT(C.yices_zero())
	}
	return TermT(C.yices_poly_mpqp(count, (*C.mpq_t)(&q[0]), (*C.term_t)(&t[0])))
}

/*
 * ARITHMETIC ATOMS
 */

// ArithEqAtom is the go version of yices_arith_eq_atom.
func ArithEqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_eq_atom(C.term_t(t1), C.term_t(t2)))
}

// ArithNeqAtom is the go version of yices_arith_neq_atom.
func ArithNeqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_neq_atom(C.term_t(t1), C.term_t(t2)))
}

// ArithGeqAtom is the go version of yices_arith_geq_atom.
func ArithGeqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_geq_atom(C.term_t(t1), C.term_t(t2)))
}

// ArithLeqAtom is the go version of yices_arith_leq_atom.
func ArithLeqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_leq_atom(C.term_t(t1), C.term_t(t2)))
}

// ArithGtAtom is the go version of yices_arith_gt_atom.
func ArithGtAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_gt_atom(C.term_t(t1), C.term_t(t2)))
}

// ArithLtAtom is the go version of yices_arith_lt_atom.
func ArithLtAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_arith_lt_atom(C.term_t(t1), C.term_t(t2)))
}

// ArithEq0Atom is the go version of yices_arith_eq0_atom.
func ArithEq0Atom(t TermT) TermT {
	return TermT(C.yices_arith_eq0_atom(C.term_t(t)))
}

// ArithNeq0Atom is the go version of yices_arith_neq0_atom.
func ArithNeq0Atom(t TermT) TermT {
	return TermT(C.yices_arith_neq0_atom(C.term_t(t)))
}

// ArithGeq0Atom is the go version of yices_arith_geq0_atom.
func ArithGeq0Atom(t TermT) TermT {
	return TermT(C.yices_arith_geq0_atom(C.term_t(t)))
}

// ArithLeq0Atom is the go version of yices_arith_leq0_atom.
func ArithLeq0Atom(t TermT) TermT {
	return TermT(C.yices_arith_leq0_atom(C.term_t(t)))
}

// ArithGt0Atom is the go version of yices_arith_gt0_atom.
func ArithGt0Atom(t TermT) TermT {
	return TermT(C.yices_arith_gt0_atom(C.term_t(t)))
}

// ArithLt0Atom is the go version of yices_arith_lt0_atom.
func ArithLt0Atom(t TermT) TermT {
	return TermT(C.yices_arith_lt0_atom(C.term_t(t)))
}

/*********************************
 *  BITVECTOR TERM CONSTRUCTORS  *
 ********************************/

// BvconstUint32 is the go version of yices_bvconst_uint32.
func BvconstUint32(bits uint32, x uint32) TermT {
	return TermT(C.yices_bvconst_uint32(C.uint32_t(bits), C.uint32_t(x)))
}

// BvconstUint64 is the go version of yices_bvconst_uint64.
func BvconstUint64(bits uint32, x uint64) TermT {
	return TermT(C.yices_bvconst_uint64(C.uint32_t(bits), C.uint64_t(x)))
}

// BvconstInt32 is the go version of yices_bvconst_int32.
func BvconstInt32(bits uint32, x int32) TermT {
	return TermT(C.yices_bvconst_int32(C.uint32_t(bits), C.int32_t(x)))
}

// BvconstInt64 is the go version of yices_bvconst_int64.
func BvconstInt64(bits uint32, x int64) TermT {
	return TermT(C.yices_bvconst_int64(C.uint32_t(bits), C.int64_t(x)))
}

// BvconstMpz is the go version of yices_bvconst_mpz.
func BvconstMpz(bits uint32, z MpzT) TermT {
	return TermT(C.yices_bvconst_mpzp(C.uint32_t(bits), (*C.mpz_t)(unsafe.Pointer(&z))))
}

// BvconstZero is the go version of yices_bvconst_zero.
func BvconstZero(bits uint32) TermT {
	return TermT(C.yices_bvconst_zero(C.uint32_t(bits)))
}

// BvconstOne is the go version of yices_bvconst_one.
func BvconstOne(bits uint32) TermT {
	return TermT(C.yices_bvconst_one(C.uint32_t(bits)))
}

// BvconstMinusOne is the go version of yices_bvconst_minus_one.
func BvconstMinusOne(bits uint32) TermT {
	return TermT(C.yices_bvconst_minus_one(C.uint32_t(bits)))
}

// BvconstFromArray is the go version of yices_bvconst_from_array.
//iam: FIXME check that bits is restricted to len(a)
func BvconstFromArray(a []int32) TermT {
	bits := C.uint32_t(len(a))
	//iam: FIXME need to unify the yices errors and the go errors...
	if bits == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvconst_from_array(bits, (*C.int32_t)(&a[0])))
}

// ParseBvbin is the go version of yices_parse_bvbin.
func ParseBvbin(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_bvbin(cs))
}

// ParseBvhex is the go version of yices_parse_bvhex.
func ParseBvhex(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_bvhex(cs))
}

// Bvadd is the go version of yices_bvadd.
func Bvadd(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvadd(C.term_t(t1), C.term_t(t2)))
}

// Bvsub is the go version of yices_bvsub.
func Bvsub(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsub(C.term_t(t1), C.term_t(t2)))
}

// Bvneg is the go version of yices_bvneg.
func Bvneg(t TermT) TermT {
	return TermT(C.yices_bvneg(C.term_t(t)))
}

// Bvmul is the go version of yices_bvmul.
func Bvmul(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvmul(C.term_t(t1), C.term_t(t2)))
}

// Bvsquare is the go version of yices_bvsquare.
func Bvsquare(t TermT) TermT {
	return TermT(C.yices_bvsquare(C.term_t(t)))
}

// Bvpower is the go version of yices_bvpower.
func Bvpower(t1 TermT, d uint32) TermT {
	return TermT(C.yices_bvpower(C.term_t(t1), C.uint32_t(d)))
}

// Bvdiv is the go version of yices_bvdiv.
func Bvdiv(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvdiv(C.term_t(t1), C.term_t(t2)))
}

// Bvrem is the go version of yices_bvrem.
func Bvrem(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvrem(C.term_t(t1), C.term_t(t2)))
}

// Bvsdiv is the go version of yices_bvsdiv.
func Bvsdiv(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsdiv(C.term_t(t1), C.term_t(t2)))
}

// Bvsrem is the go version of yices_bvsrem.
func Bvsrem(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsrem(C.term_t(t1), C.term_t(t2)))
}

// Bvsmod is the go version of yices_bvsmod.
func Bvsmod(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsmod(C.term_t(t1), C.term_t(t2)))
}

// Bvnot is the go version of yices_bvnot.
func Bvnot(t TermT) TermT {
	return TermT(C.yices_bvnot(C.term_t(t)))
}

// Bvnand is the go version of yices_bvnand.
func Bvnand(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvnand(C.term_t(t1), C.term_t(t2)))
}

// Bvnor is the go version of yices_bvnor.
func Bvnor(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvnor(C.term_t(t1), C.term_t(t2)))
}

// Bvxnor is the go version of yices_bvxnor.
func Bvxnor(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvxnor(C.term_t(t1), C.term_t(t2)))
}

// Bvshl is the go version of yices_bvshl.
func Bvshl(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvshl(C.term_t(t1), C.term_t(t2)))
}

// Bvlshr is the go version of yices_bvlshr.
func Bvlshr(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvlshr(C.term_t(t1), C.term_t(t2)))
}

// Bvashr is the go version of yices_bvashr.
func Bvashr(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvashr(C.term_t(t1), C.term_t(t2)))
}

// Bvand is the go version of yices_bvand.
func Bvand(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvand(count, (*C.term_t)(&t[0])))
}

// Bvor is the go version of yices_bvor.
func Bvor(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvor(count, (*C.term_t)(&t[0])))
}

// Bvxor is the go version of yices_bvxor.
func Bvxor(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvxor(count, (*C.term_t)(&t[0])))
}

// Bvand2 is the go version of yices_bvand2.
func Bvand2(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvand2(C.term_t(t1), C.term_t(t2)))
}

// Bvor2 is the go version of yices_bvor2.
func Bvor2(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvor2(C.term_t(t1), C.term_t(t2)))
}

// Bvxor2 is the go version of yices_bvxor2.
func Bvxor2(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvxor2(C.term_t(t1), C.term_t(t2)))
}

// Bvand3 is the go version of yices_bvand3.
func Bvand3(t1 TermT, t2 TermT, t3 TermT) TermT {
	return TermT(C.yices_bvand3(C.term_t(t1), C.term_t(t2), C.term_t(t3)))
}

// Bvor3 is the go version of yices_bvor3.
func Bvor3(t1 TermT, t2 TermT, t3 TermT) TermT {
	return TermT(C.yices_bvor3(C.term_t(t1), C.term_t(t2), C.term_t(t3)))
}

// Bvxor3 is the go version of yices_bvxor3.
func Bvxor3(t1 TermT, t2 TermT, t3 TermT) TermT {
	return TermT(C.yices_bvxor3(C.term_t(t1), C.term_t(t2), C.term_t(t3)))
}

// Bvsum is the go version of yices_bvsum.
func Bvsum(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvsum(count, (*C.term_t)(&t[0])))
}

// Bvproduct is the go version of yices_bvproduct.
func Bvproduct(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvproduct(count, (*C.term_t)(&t[0])))
}

// ShiftLeft0 is the go version of yices_shift_left0.
func ShiftLeft0(t TermT, n uint32) TermT {
	return TermT(C.yices_shift_left0(C.term_t(t), C.uint32_t(n)))
}

// ShiftLeft1 is the go version of yices_shift_left1.
func ShiftLeft1(t TermT, n uint32) TermT {
	return TermT(C.yices_shift_left1(C.term_t(t), C.uint32_t(n)))
}

// ShiftRight0 is the go version of yices_shift_right0.
func ShiftRight0(t TermT, n uint32) TermT {
	return TermT(C.yices_shift_right0(C.term_t(t), C.uint32_t(n)))
}

// ShiftRight1 is the go version of yices_shift_right1.
func ShiftRight1(t TermT, n uint32) TermT {
	return TermT(C.yices_shift_right1(C.term_t(t), C.uint32_t(n)))
}

// AshiftRight is the go version of yices_ashift_right.
func AshiftRight(t TermT, n uint32) TermT {
	return TermT(C.yices_ashift_right(C.term_t(t), C.uint32_t(n)))
}

// RotateLeft is the go version of yices_rotate_left.
func RotateLeft(t TermT, n uint32) TermT {
	return TermT(C.yices_rotate_left(C.term_t(t), C.uint32_t(n)))
}

// RotateRight is the go version of yices_rotate_right.
func RotateRight(t TermT, n uint32) TermT {
	return TermT(C.yices_rotate_right(C.term_t(t), C.uint32_t(n)))
}

// Bvextract is the go version of yices_bvextract.
func Bvextract(t TermT, i uint32, j uint32) TermT {
	return TermT(C.yices_bvextract(C.term_t(t), C.uint32_t(i), C.uint32_t(j)))
}

// Bvconcat2 is the go version of yices_bvconcat2.
func Bvconcat2(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvconcat2(C.term_t(t1), C.term_t(t2)))
}

// Bvconcat is the go version of yices_bvconcat.
func Bvconcat(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvconcat(count, (*C.term_t)(&t[0])))
}

// Bvrepeat is the go version of yices_bvrepeat.
func Bvrepeat(t TermT, n uint32) TermT {
	return TermT(C.yices_bvrepeat(C.term_t(t), C.uint32_t(n)))
}

// SignExtend is the go version of yices_sign_extend.
func SignExtend(t TermT, n uint32) TermT {
	return TermT(C.yices_sign_extend(C.term_t(t), C.uint32_t(n)))
}

// ZeroExtend is the go version of yices_zero_extend.
func ZeroExtend(t TermT, n uint32) TermT {
	return TermT(C.yices_zero_extend(C.term_t(t), C.uint32_t(n)))
}

// Redand is the go version of yices_redand.
func Redand(t TermT) TermT {
	return TermT(C.yices_redand(C.term_t(t)))
}

// Redor is the go version of yices_redor.
func Redor(t TermT) TermT {
	return TermT(C.yices_redor(C.term_t(t)))
}

// Redcomp is the go version of yices_redcomp.
func Redcomp(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_redcomp(C.term_t(t1), C.term_t(t2)))
}

// Bvarray is the go version of yices_bvarray.
func Bvarray(t []TermT) TermT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_bvarray(count, (*C.term_t)(&t[0])))
}

// Bitextract is the go version of yices_bitextract.
func Bitextract(t TermT, n uint32) TermT {
	return TermT(C.yices_bitextract(C.term_t(t), C.uint32_t(n)))
}

// BveqAtom is the go version of yices_bveq_atom.
func BveqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bveq_atom(C.term_t(t1), C.term_t(t2)))
}

// BvneqAtom is the go version of yices_bvneq_atom.
func BvneqAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvneq_atom(C.term_t(t1), C.term_t(t2)))
}

// BvgeAtom is the go version of yices_bvge_atom.
func BvgeAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvge_atom(C.term_t(t1), C.term_t(t2)))
}

// BvgtAtom is the go version of yices_bvgt_atom.
func BvgtAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvgt_atom(C.term_t(t1), C.term_t(t2)))
}

// BvleAtom is the go version of yices_bvle_atom.
func BvleAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvle_atom(C.term_t(t1), C.term_t(t2)))
}

// BvltAtom is the go version of yices_bvlt_atom.
func BvltAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvlt_atom(C.term_t(t1), C.term_t(t2)))
}

// BvsgeAtom is the go version of yices_bvsge_atom.
func BvsgeAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsge_atom(C.term_t(t1), C.term_t(t2)))
}

// BvsgtAtom is the go version of yices_bvsgt_atom.
func BvsgtAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsgt_atom(C.term_t(t1), C.term_t(t2)))
}

// BvsleAtom is the go version of yices_bvsle_atom.
func BvsleAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvsle_atom(C.term_t(t1), C.term_t(t2)))
}

// BvsltAtom is the go version of yices_bvslt_atom.
func BvsltAtom(t1 TermT, t2 TermT) TermT {
	return TermT(C.yices_bvslt_atom(C.term_t(t1), C.term_t(t2)))
}

/**************
 *  PARSING   *
 *************/

// ParseType is the go version of yices_parse_type.
func ParseType(s string) TypeT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TypeT(C.yices_parse_type(cs))
}

// ParseTerm is the go version of yices_parse_term.
func ParseTerm(s string) TermT {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_parse_term(cs))
}

/*******************
 *  SUBSTITUTIONS  *
 ******************/

// SubstTerm is the go version of yices_subst_term.
func SubstTerm(vars []TermT, vals []TermT, t TermT) TermT {
	count := C.uint32_t(len(vars))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return NullTerm
	}
	return TermT(C.yices_subst_term(count, (*C.term_t)(&vars[0]), (*C.term_t)(&vals[0]), C.term_t(t)))
}

// SubstTermArray is the go version of yices_subst_term_array.
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

// SetTypeName is the go version of yices_set_type_name.
func SetTypeName(tau TypeT, name string) int32 {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return int32(C.yices_set_type_name(C.type_t(tau), cs))
}

// SetTermName is the go version of yices_set_term_name.
func SetTermName(t TermT, name string) int32 {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return int32(C.yices_set_term_name(C.term_t(t), cs))
}

// RemoveTypeName is the go version of yices_remove_type_name.
func RemoveTypeName(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.yices_remove_type_name(cs)
}

// RemoveTermName is the go version of yices_remove_term_name.
func RemoveTermName(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.yices_remove_term_name(cs)
}

// GetTypeByName is the go version of yices_get_type_by_name.
func GetTypeByName(name string) TypeT {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return TypeT(C.yices_get_type_by_name(cs))
}

// GetTermByName is the go version of yices_get_term_by_name.
func GetTermByName(name string) TermT {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	return TermT(C.yices_get_term_by_name(cs))
}

// ClearTypeName is the go version of yices_clear_type_name.
func ClearTypeName(tau TypeT) int32 {
	return int32(C.yices_clear_type_name(C.type_t(tau)))
}

// ClearTermName is the go version of yices_clear_term_name.
func ClearTermName(t TermT) int32 {
	return int32(C.yices_clear_term_name(C.term_t(t)))
}

// GetTypeName is the go version of yices_get_type_name.
func GetTypeName(tau TypeT) string {
	//FIXME: check if the name needs to be freed
	return C.GoString(C.yices_get_type_name(C.type_t(tau)))
}

// GetTermName is the go version of yices_get_term_name.
func GetTermName(t TermT) string {
	//FIXME: check if the name needs to be freed
	return C.GoString(C.yices_get_term_name(C.term_t(t)))
}

/***********************
 *  TERM EXPLORATION   *
 **********************/

// TypeOfTerm is the go version of yices_type_of_term.
func TypeOfTerm(t TermT) TypeT {
	return TypeT(C.yices_type_of_term(C.term_t(t)))
}

// TermIsBool is the go version of yices_term_is_bool.
func TermIsBool(t TermT) bool {
	return C.yices_term_is_bool(C.term_t(t)) == C.int32_t(1)
}

// TermIsInt is the go version of yices_term_is_int.
func TermIsInt(t TermT) bool {
	return C.yices_term_is_int(C.term_t(t)) == C.int32_t(1)
}

// TermIsReal is the go version of yices_term_is_real.
func TermIsReal(t TermT) bool {
	return C.yices_term_is_real(C.term_t(t)) == C.int32_t(1)
}

// TermIsArithmetic is the go version of yices_term_is_arithmetic.
func TermIsArithmetic(t TermT) bool {
	return C.yices_term_is_arithmetic(C.term_t(t)) == C.int32_t(1)
}

// TermIsBitvector is the go version of yices_term_is_bitvector.
func TermIsBitvector(t TermT) bool {
	return C.yices_term_is_bitvector(C.term_t(t)) == C.int32_t(1)
}

// TermIsTuple is the go version of yices_term_is_tuple.
func TermIsTuple(t TermT) bool {
	return C.yices_term_is_tuple(C.term_t(t)) == C.int32_t(1)
}

// TermIsFunction is the go version of yices_term_is_function.
func TermIsFunction(t TermT) bool {
	return C.yices_term_is_function(C.term_t(t)) == C.int32_t(1)
}

// TermIsScalar is the go version of yices_term_is_scalar.
func TermIsScalar(t TermT) bool {
	return C.yices_term_is_scalar(C.term_t(t)) == C.int32_t(1)
}

// TermBitsize is the go version of yices_term_bitsize.
func TermBitsize(t TermT) uint32 {
	return uint32(C.yices_term_bitsize(C.term_t(t)))
}

// TermIsGround is the go version of yices_term_is_ground.
func TermIsGround(t TermT) bool {
	return C.yices_term_is_ground(C.term_t(t)) == C.int32_t(1)
}

// TermIsAtomic is the go version of yices_term_is_atomic.
func TermIsAtomic(t TermT) bool {
	return C.yices_term_is_atomic(C.term_t(t)) == C.int32_t(1)
}

// TermIsComposite is the go version of yices_term_is_composite.
func TermIsComposite(t TermT) bool {
	return C.yices_term_is_composite(C.term_t(t)) == C.int32_t(1)
}

// TermIsProjection is the go version of yices_term_is_projection.
func TermIsProjection(t TermT) bool {
	return C.yices_term_is_projection(C.term_t(t)) == C.int32_t(1)
}

// TermIsSum is the go version of yices_term_is_sum.
func TermIsSum(t TermT) bool {
	return C.yices_term_is_sum(C.term_t(t)) == C.int32_t(1)
}

// TermIsBvsum is the go version of yices_term_is_bvsum.
func TermIsBvsum(t TermT) bool {
	return C.yices_term_is_bvsum(C.term_t(t)) == C.int32_t(1)
}

// TermIsProduct is the go version of yices_term_is_product.
func TermIsProduct(t TermT) bool {
	return C.yices_term_is_product(C.term_t(t)) == C.int32_t(1)
}

// TermConstructor is the go version of yices_term_constructor.
func TermConstructor(t TermT) TermConstructorT {
	return TermConstructorT(C.yices_term_constructor(C.term_t(t)))
}

// TermNumChildren is the go version of yices_term_num_children.
func TermNumChildren(t TermT) int32 {
	return int32(C.yices_term_num_children(C.term_t(t)))
}

// TermChild is the go version of yices_term_child.
func TermChild(t TermT, i int32) TermT {
	return TermT(C.yices_term_child(C.term_t(t), C.int32_t(i)))
}

// TermChildren is the go version of yices_term_children.
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

// ProjIndex is the go version of yices_proj_index.
func ProjIndex(t TermT) int32 {
	return int32(C.yices_proj_index(C.term_t(t)))
}

// ProjArg is the go version of yices_proj_arg.
func ProjArg(t TermT) TermT {
	return TermT(C.yices_proj_arg(C.term_t(t)))
}

// BoolConstValue is the go version of yices_bool_const_value.
func BoolConstValue(t TermT, val *int32) int32 {
	return int32(C.yices_bool_const_value(C.term_t(t), (*C.int32_t)(val)))
}

// BvConstValue is the go version of yices_bv_const_value.
func BvConstValue(t TermT, val []int32) int32 {
	return int32(C.yices_bv_const_value(C.term_t(t), (*C.int32_t)(&val[0])))
}

// ScalarConstValue is the go version of yices_scalar_const_value.
func ScalarConstValue(t TermT, val *int32) int32 {
	return int32(C.yices_scalar_const_value(C.term_t(t), (*C.int32_t)(val)))
}

// RationalConstValue is the go version of yices_rational_const_value.
func RationalConstValue(t TermT, q *MpqT) int32 {
	return int32(C.yices_rational_const_valuep(C.term_t(t), (*C.mpq_t)(unsafe.Pointer(q))))
}

// SumComponent is the go version of yices_sum_componentp.
func SumComponent(t TermT, i int32, coeff *MpqT, term *TermT) int32 {
	return int32(C.yices_sum_componentp(C.term_t(t), C.int32_t(i), (*C.mpq_t)(unsafe.Pointer(coeff)), (*C.term_t)(term)))
}

// BvsumComponent is the go version of yices_bvsum_component.
func BvsumComponent(t TermT, i int32, val []int32, term *TermT) int32 {
	return int32(C.yices_bvsum_component(C.term_t(t), C.int32_t(i), (*C.int32_t)(&val[0]), (*C.term_t)(term)))
}

// ProductComponent is the go version of yices_product_component.
func ProductComponent(t TermT, i int32, term *TermT, exp *uint32) int32 {
	return int32(C.yices_product_component(C.term_t(t), C.int32_t(i), (*C.term_t)(term), (*C.uint32_t)(exp)))
}

/*************************
 *  GARBAGE COLLECTION   *
 ************************/

// NumTerms is the go version of yices_num_terms.
func NumTerms() uint32 {
	return uint32(C.yices_num_terms())
}

// NumTypes is the go version of yices_num_types.
func NumTypes() uint32 {
	return uint32(C.yices_num_types())
}

// IncrefTerm is the go version of yices_incref_term.
func IncrefTerm(t TermT) int32 {
	return int32(C.yices_incref_term(C.term_t(t)))
}

// DecrefTerm is the go version of yices_decref_term.
func DecrefTerm(t TermT) int32 {
	return int32(C.yices_decref_term(C.term_t(t)))
}

// IncrefType is the go version of yices_incref_type.
func IncrefType(tau TypeT) int32 {
	return int32(C.yices_incref_type(C.type_t(tau)))
}

// DecrefType is the go version of yices_decref_type.
func DecrefType(tau TypeT) int32 {
	return int32(C.yices_decref_type(C.type_t(tau)))
}

// NumPosrefTerms is the go version of yices_num_posref_terms.
func NumPosrefTerms() uint32 {
	return uint32(C.yices_num_posref_terms())
}

// NumPosrefTypes is the go version of yices_num_posref_types.
func NumPosrefTypes() uint32 {
	return uint32(C.yices_num_posref_types())
}

// GarbageCollect is the go version of yices_garbage_collect.
func GarbageCollect(ts []TermT, taus []TypeT, keepNamed int32) {
	tCount := C.uint32_t(len(ts))
	tauCount := C.uint32_t(len(taus))
	C.yices_garbage_collect((*C.term_t)(&ts[0]), tCount, (*C.type_t)(&taus[0]), tauCount, C.int32_t(keepNamed))
}

/****************************
 *  CONTEXT CONFIGURATION   *
 ***************************/

// ConfigT is a thin wrapper around a C pointer to a ctx_config_t struct. It is stored as a uintptr to avoid the scrutiny of the go GC.
type ConfigT struct {
	raw uintptr // actually *C.ctx_config_t
}

func ycfg(cfg ConfigT) *C.ctx_config_t {
	return (*C.ctx_config_t)(unsafe.Pointer(cfg.raw))
}

// InitConfig is the go version of yices_new_config.
func InitConfig(cfg *ConfigT) {
	cfg.raw = uintptr(unsafe.Pointer(C.yices_new_config()))
}

// CloseConfig is the go version of yices_free_config.
func CloseConfig(cfg *ConfigT) {
	C.yices_free_config(ycfg(*cfg))
	cfg.raw = 0
}

// SetConfig is the go version of yices_set_config.
func SetConfig(cfg ConfigT, name string, value string) int32 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	return int32(C.yices_set_config(ycfg(cfg), cname, cvalue))
}

// DefaultConfigForLogic is the go version of yices_default_config_for_logic.
func DefaultConfigForLogic(cfg ConfigT, logic string) int32 {
	clogic := C.CString(logic)
	defer C.free(unsafe.Pointer(clogic))
	return int32(C.yices_default_config_for_logic(ycfg(cfg), clogic))
}

/***************
 *  CONTEXTS   *
 **************/

// ContextT is a thin wrapper around a C pointer to a context_t struct. It is stored as a uintptr to avoid the scrutiny of the go GC.
type ContextT struct {
	raw uintptr // actually *C.context_t
}

func yctx(ctx ContextT) *C.context_t {
	return (*C.context_t)(unsafe.Pointer(ctx.raw))
}

// InitContext is the go version of yices_new_context.
func InitContext(cfg ConfigT, ctx *ContextT) {
	ctx.raw = uintptr(unsafe.Pointer(C.yices_new_context(ycfg(cfg))))
}

// CloseContext is the go version of yices_free_context.
func CloseContext(ctx *ContextT) {
	C.yices_free_context(yctx(*ctx))
	ctx.raw = 0
}

// ContextStatus is the go version of yices_context_status.
func ContextStatus(ctx ContextT) SmtStatusT {
	return SmtStatusT(C.yices_context_status(yctx(ctx)))
}

// ResetContext is the go version of yices_reset_context.
func ResetContext(ctx ContextT) {
	C.yices_reset_context(yctx(ctx))
}

// Push is the go version of yices_push.
func Push(ctx ContextT) int32 {
	return int32(C.yices_push(yctx(ctx)))
}

// Pop is the go version of yices_pop.
func Pop(ctx ContextT) int32 {
	return int32(C.yices_pop(yctx(ctx)))
}

// ContextEnableOption is the go version of yices_context_enable_option.
func ContextEnableOption(ctx ContextT, option string) int32 {
	coption := C.CString(option)
	defer C.free(unsafe.Pointer(coption))
	return int32(C.yices_context_enable_option(yctx(ctx), coption))
}

// ContextDisableOption is the go version of yices_context_enable_option.
func ContextDisableOption(ctx ContextT, option string) int32 {
	coption := C.CString(option)
	defer C.free(unsafe.Pointer(coption))
	return int32(C.yices_context_enable_option(yctx(ctx), coption))
}

// AssertFormula is the go version of yices_assert_formula.
func AssertFormula(ctx ContextT, t TermT) int32 {
	return int32(C.yices_assert_formula(yctx(ctx), C.term_t(t)))

}

// AssertFormulas is the go version of yices_assert_formulas.
func AssertFormulas(ctx ContextT, t []TermT) int32 {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return -1
	}
	return int32(C.yices_assert_formulas(yctx(ctx), count, (*C.term_t)(&t[0])))
}

// CheckContext is the go version of yices_check_context.
func CheckContext(ctx ContextT, params ParamT) SmtStatusT {
	return SmtStatusT(C.yices_check_context(yctx(ctx), yparam(params)))
}

// CheckContextWithAssumptions is the go version of yices_check_context_with_assumptions.
func CheckContextWithAssumptions(ctx ContextT, params ParamT, t []TermT) SmtStatusT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return SmtStatusT(StatusError)
	}
	return SmtStatusT(C.yices_check_context_with_assumptions(yctx(ctx), yparam(params), count, (*C.term_t)(&t[0])))
}

// CheckContextWithModel is the go version of yices_check_context_with_model (new in 2.6.4).
func CheckContextWithModel(ctx ContextT, params ParamT, model ModelT, t []TermT) SmtStatusT {
	count := C.uint32_t(len(t))
	//iam: FIXME need to unify the yices errors and the go errors...
	if count == 0 {
		return SmtStatusT(StatusError)
	}
	return SmtStatusT(C.yices_check_context_with_model(yctx(ctx), yparam(params), ymodel(model), count, (*C.term_t)(&t[0])))
}

// CheckContextWithInterpolation is the go version of yices_check_context_with_interpolation (new in 2.6.4).
func CheckContextWithInterpolation(ctxA ContextT, ctxB ContextT, params ParamT, buildModel bool) (status SmtStatusT, modelp *ModelT, interpolantp *TermT) {
	modelp = nil
	interpolantp = nil
	ictx := C.new_interpolation_context(yctx(ctxA), yctx(ctxB))
	defer C.free_interpolation_context(ictx)
	bm := 0
	if buildModel {
		bm = 1
	}
	cstatus := C.yices_check_context_with_interpolation(ictx, yparam(params), C.int32_t(bm))

	if cstatus == C.STATUS_SAT {
		modelp = &ModelT{uintptr(unsafe.Pointer(C.get_interpolation_context_model(ictx)))}
	} else if cstatus == C.STATUS_UNSAT {
		term := TermT(C.get_interpolation_context_interpolant(ictx))
		interpolantp = &term
	}
	status = SmtStatusT(cstatus)
	return
}

// AssertBlockingClause is the go version of yices_assert_blocking_clause.
func AssertBlockingClause(ctx ContextT) int32 {
	return int32(C.yices_assert_blocking_clause(yctx(ctx)))

}

// StopSearch is the go version of yices_stop_search.
func StopSearch(ctx ContextT) {
	C.yices_stop_search(yctx(ctx))
}

/*
 * SEARCH PARAMETERS
 */

// ParamT is a thin wrapper around a C pointer to a param_t struct. It is stored as a uintptr to avoid the scrutiny of the go GC.
type ParamT struct {
	raw uintptr // actually *C.param_t
}

func yparam(params ParamT) *C.param_t {
	return (*C.param_t)(unsafe.Pointer(params.raw))
}

// InitParamRecord is the go version of yices_new_param_record.
func InitParamRecord(params *ParamT) {
	params.raw = uintptr(unsafe.Pointer(C.yices_new_param_record()))
}

// CloseParamRecord is the go version of yices_free_param_record.
func CloseParamRecord(params *ParamT) {
	C.yices_free_param_record(yparam(*params))
	params.raw = 0
}

// DefaultParamsForContext is the go version of yices_default_params_for_context.
func DefaultParamsForContext(ctx ContextT, params ParamT) {
	C.yices_default_params_for_context(yctx(ctx), yparam(params))
}

// SetParam is the go version of yices_set_param.
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

// GetUnsatCore is the go version of yices_get_unsat_core.
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

// GetModelInterpolant is the go version of yices_get_model_interpolant (new in 2.6.4).
func GetModelInterpolant(ctx ContextT) (interpolant TermT) {
	interpolant = TermT(C.yices_get_model_interpolant(yctx(ctx)))
	return
}

/**************
 *   MODELS   *
 *************/

// ModelT is a thin wrapper around a C pointer to a model_t struct. It is stored as a uintptr to avoid the scrutiny of the go GC.
type ModelT struct {
	raw uintptr // actually *C.model_t
}

func ymodel(model ModelT) *C.model_t {
	return (*C.model_t)(unsafe.Pointer(model.raw))
}

// NewModel is the go version of yices_new_model (new in 2.6.4).
func NewModel() *ModelT {
	return &ModelT{uintptr(unsafe.Pointer(C.yices_new_model()))}
}

// GetModel is the go version of yices_get_model.
func GetModel(ctx ContextT, keepSubst int32) *ModelT {
	//yes golang lets you return stuff allocated on the stack
	return &ModelT{uintptr(unsafe.Pointer(C.yices_get_model(yctx(ctx), C.int32_t(keepSubst))))}
}

// CloseModel is the go version of yices_free_model.
func CloseModel(model *ModelT) {
	C.yices_free_model(ymodel(*model))
	model.raw = 0
}

// ModelFromMap is the go version of yices_model_from_map.
func ModelFromMap(vars []TermT, vals []TermT) *ModelT {
	vcount := C.uint32_t(len(vals))
	return &ModelT{uintptr(unsafe.Pointer(C.yices_model_from_map(vcount, (*C.term_t)(&vars[0]), (*C.term_t)(&vals[0]))))}
}

// ModelCollectDefinedTerms is the go version of yices_model_collect_defined_terms.
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

// GetBoolValue is the go version of yices_get_bool_value.
func GetBoolValue(model ModelT, t TermT, val *int32) int32 {
	return int32(C.yices_get_bool_value(ymodel(model), C.term_t(t), (*C.int32_t)(val)))
}

// GetInt32Value is the go version of yices_get_int32_value.
func GetInt32Value(model ModelT, t TermT, val *int32) int32 {
	return int32(C.yices_get_int32_value(ymodel(model), C.term_t(t), (*C.int32_t)(val)))
}

// GetInt64Value is the go version of yices_get_int64_value.
func GetInt64Value(model ModelT, t TermT, val *int64) int32 {
	return int32(C.yices_get_int64_value(ymodel(model), C.term_t(t), (*C.int64_t)(val)))
}

// GetRational32Value is the go version of yices_get_rational32_value.
func GetRational32Value(model ModelT, t TermT, num *int32, den *uint32) int32 {
	return int32(C.yices_get_rational32_value(ymodel(model), C.term_t(t), (*C.int32_t)(num), (*C.uint32_t)(den)))
}

// GetRational64Value is the go version of yices_get_rational64_value.
func GetRational64Value(model ModelT, t TermT, num *int64, den *uint64) int32 {
	return int32(C.yices_get_rational64_value(ymodel(model), C.term_t(t), (*C.int64_t)(num), (*C.uint64_t)(den)))
}

// GetDoubleValue is the go version of yices_get_double_value.
func GetDoubleValue(model ModelT, t TermT, val *float64) int32 {
	return int32(C.yices_get_double_value(ymodel(model), C.term_t(t), (*C.double)(val)))
}

// GetMpzValue is the go version of yices_get_mpz_value.
func GetMpzValue(model ModelT, t TermT, val *MpzT) int32 {
	return int32(C.yices_get_mpz_valuep(ymodel(model), C.term_t(t), (*C.mpz_t)(unsafe.Pointer(val))))
}

// GetMpqValue is the go version of yices_get_mpq_value.
func GetMpqValue(model ModelT, t TermT, val *MpqT) int32 {
	return int32(C.yices_get_mpq_valuep(ymodel(model), C.term_t(t), (*C.mpq_t)(unsafe.Pointer(val))))
}

/*
//iam: not gonna assume mcsat
#ifdef LIBPOLY_VERSION
__YICES_DLLSPEC__ extern int32_t yices_get_algebraic_number_value(model_t *mdl, term_t t, lp_algebraic_number_t *a);
#endif
*/

// GetBvValue is the go version of yices_get_bv_value.
func GetBvValue(model ModelT, t TermT, val []int32) int32 {
	return int32(C.yices_get_bv_value(ymodel(model), C.term_t(t), (*C.int32_t)(&val[0])))
}

// GetScalarValue is the go version of yices_get_scalar_value.
func GetScalarValue(model ModelT, t TermT, val *int32) int32 {
	return int32(C.yices_get_scalar_value(ymodel(model), C.term_t(t), (*C.int32_t)(val)))
}

/******************************
 *  SETTING VALUES IN A MODEL *
 ******************************/

// SetBoolValue is the go version of yices_model_set_bool (new in 2.6.4).
func SetBoolValue(model ModelT, t TermT, val int32) int32 {
	return int32(C.yices_model_set_bool(ymodel(model), C.term_t(t), C.int32_t(val)))
}

// SetInt32Value is the go version of yices_model_set_int32 (new in 2.6.4).
func SetInt32Value(model ModelT, t TermT, val int32) int32 {
	return int32(C.yices_model_set_int32(ymodel(model), C.term_t(t), C.int32_t(val)))
}

// SetInt64Value is the go version of yices_model_set_int64 (new in 2.6.4).
func SetInt64Value(model ModelT, t TermT, val int64) int32 {
	return int32(C.yices_model_set_int64(ymodel(model), C.term_t(t), C.int64_t(val)))
}

// SetRational32Value is the go version of yices_model_set_rational32 (new in 2.6.4).
func SetRational32Value(model ModelT, t TermT, num int32, den int32) int32 {
	return int32(C.yices_model_set_rational32(ymodel(model), C.term_t(t), C.int32_t(num), C.uint32_t(den)))
}

// SetRational64Value is the go version of yices_model_set_rational64 (new in 2.6.4).
func SetRational64Value(model ModelT, t TermT, num int64, den int64) int32 {
	return int32(C.yices_model_set_rational64(ymodel(model), C.term_t(t), C.int64_t(num), C.uint64_t(den)))
}

// SetMpzValue is the go version of yices_model_set_mpz_value.
func SetMpzValue(model ModelT, t TermT, val *MpzT) int32 {
	return int32(C.yices_model_set_mpzp(ymodel(model), C.term_t(t), (*C.mpz_t)(unsafe.Pointer(val))))
}

// SetMpqValue is the go version of yices_model_set_mpq_value.
func SetMpqValue(model ModelT, t TermT, val *MpqT) int32 {
	return int32(C.yices_model_set_mpqp(ymodel(model), C.term_t(t), (*C.mpq_t)(unsafe.Pointer(val))))
}

// iam: not going to assume mcsat
//#ifdef LIBPOLY_VERSION
//__YICES_DLLSPEC__ extern int32_t yices_model_set_algebraic_number(model_t *model, term_t var, const lp_algebraic_number_t *val);
//#endif

// SetBvInt32Value is the go version of yices_model_set_bv_int32 (new in 2.6.4).
func SetBvInt32Value(model ModelT, t TermT, val int32) int32 {
	return int32(C.yices_model_set_bv_int32(ymodel(model), C.term_t(t), C.int32_t(val)))
}

// SetBvInt64Value is the go version of yices_model_set_bv_int64 (new in 2.6.4).
func SetBvInt64Value(model ModelT, t TermT, val int64) int32 {
	return int32(C.yices_model_set_bv_int64(ymodel(model), C.term_t(t), C.int64_t(val)))
}

// SetBvUint32Value is the go version of yices_model_set_bv_uint32 (new in 2.6.4).
func SetBvUint32Value(model ModelT, t TermT, val uint32) int32 {
	return int32(C.yices_model_set_bv_uint32(ymodel(model), C.term_t(t), C.uint32_t(val)))
}

// SetBvUint64Value is the go version of yices_model_set_bv_uint64 (new in 2.6.4).
func SetBvUint64Value(model ModelT, t TermT, val uint64) int32 {
	return int32(C.yices_model_set_bv_uint64(ymodel(model), C.term_t(t), C.uint64_t(val)))
}

// SetBvMpzValue is the go version of yices_model_set_bv_mpz_value.
func SetBvMpzValue(model ModelT, t TermT, val *MpzT) int32 {
	return int32(C.yices_model_set_bv_mpzp(ymodel(model), C.term_t(t), (*C.mpz_t)(unsafe.Pointer(val))))
}

// SetBvFromArray is the go version of yices_model_bv_from_array

func SetBvFromArray(model ModelT, t TermT, a []int32) int32 {
	count := C.uint32_t(len(a))
	return int32(C.yices_model_set_bv_from_array(ymodel(model), C.term_t(t), count, (*C.int32_t)(&a[0])))
}

/*
 * GENERIC FORM: VALUE DESCRIPTORS AND NODES
 */

// YvalT is the go analog of yval_t defined in yices_types.h
type YvalT C.yval_t

// GetTag accesses the node_tag of a YvalT object.
func GetTag(yval YvalT) YvalTagT {
	return YvalTagT(yval.node_tag)
}

// YvalVectorT is the go analog of yval_vector_t defined in yices_types.h
type YvalVectorT C.yval_vector_t

// InitYvalVector is the go version of yices_init_yval_vector.
func InitYvalVector(v *YvalVectorT) {
	C.yices_init_yval_vector((*C.yval_vector_t)(v))
}

// DeleteYvalVector is the go version of yices_delete_yval_vector.
func DeleteYvalVector(v *YvalVectorT) {
	C.yices_delete_yval_vector((*C.yval_vector_t)(v))
}

// ResetYvalVector is the go version of yices_reset_yval_vector.
func ResetYvalVector(v *YvalVectorT) {
	C.yices_reset_yval_vector((*C.yval_vector_t)(v))
}

// GetValue is the go version of yices_get_value.
func GetValue(model ModelT, t TermT, val *YvalT) int32 {
	return int32(C.yices_get_value(ymodel(model), C.term_t(t), (*C.yval_t)(val)))
}

// ValIsInt32 is the go version of yices_val_is_int32.
func ValIsInt32(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_int32(ymodel(model), (*C.yval_t)(val)))
}

// ValIsInt64 is the go version of yices_val_is_int64.
func ValIsInt64(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_int64(ymodel(model), (*C.yval_t)(val)))
}

// ValIsRational32 is the go version of yices_val_is_rational32.
func ValIsRational32(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_rational32(ymodel(model), (*C.yval_t)(val)))
}

// ValIsRational64 is the go version of yices_val_is_rational64.
func ValIsRational64(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_rational64(ymodel(model), (*C.yval_t)(val)))
}

// ValIsInteger is the go version of yices_val_is_integer.
func ValIsInteger(model ModelT, val *YvalT) int32 {
	return int32(C.yices_val_is_integer(ymodel(model), (*C.yval_t)(val)))
}

// ValBitsize is the go version of yices_val_bitsize.
func ValBitsize(model ModelT, val *YvalT) uint32 {
	return uint32(C.yices_val_bitsize(ymodel(model), (*C.yval_t)(val)))
}

// ValTupleArity is the go version of yices_val_tuple_arity.
func ValTupleArity(model ModelT, val *YvalT) uint32 {
	return uint32(C.yices_val_tuple_arity(ymodel(model), (*C.yval_t)(val)))
}

// ValMappingArity is the go version of yices_val_mapping_arity.
func ValMappingArity(model ModelT, val *YvalT) uint32 {
	return uint32(C.yices_val_mapping_arity(ymodel(model), (*C.yval_t)(val)))
}

// ValFunctionArity is the go version of yices_val_function_arity.
func ValFunctionArity(model ModelT, val *YvalT) uint32 {
	return uint32(C.yices_val_function_arity(ymodel(model), (*C.yval_t)(val)))
}

// ValFunctionType is the go version of yices_val_function_type.
// Since 2.6.2
func ValFunctionType(model ModelT, val *YvalT) TypeT {
	return TypeT(C.yices_val_function_type(ymodel(model), (*C.yval_t)(val)))
}

// ValGetBool is the go version of yices_val_get_bool.
func ValGetBool(model ModelT, yval *YvalT, val *int32) int32 {
	return int32(C.yices_val_get_bool(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(val)))
}

// ValGetInt32 is the go version of yices_val_get_int32.
func ValGetInt32(model ModelT, yval *YvalT, val *int32) int32 {
	return int32(C.yices_val_get_int32(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(val)))
}

// ValGetInt64 is the go version of yices_val_get_int64.
func ValGetInt64(model ModelT, yval *YvalT, val *int64) int32 {
	return int32(C.yices_val_get_int64(ymodel(model), (*C.yval_t)(yval), (*C.int64_t)(val)))
}

// ValGetRational32 is the go version of yices_val_get_rational32.
func ValGetRational32(model ModelT, yval *YvalT, num *int32, den *uint32) int32 {
	return int32(C.yices_val_get_rational32(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(num), (*C.uint32_t)(den)))
}

// ValGetRational64 is the go version of yices_val_get_rational64.
func ValGetRational64(model ModelT, yval *YvalT, num *int64, den *uint64) int32 {
	return int32(C.yices_val_get_rational64(ymodel(model), (*C.yval_t)(yval), (*C.int64_t)(num), (*C.uint64_t)(den)))
}

// ValGetDouble is the go version of yices_val_get_double.
func ValGetDouble(model ModelT, yval *YvalT, val *float64) int32 {
	return int32(C.yices_val_get_double(ymodel(model), (*C.yval_t)(yval), (*C.double)(val)))
}

// ValGetMpz is the go version of yices_val_get_mpz.
func ValGetMpz(model ModelT, yval *YvalT, val *MpzT) int32 {
	return int32(C.yices_val_get_mpzp(ymodel(model), (*C.yval_t)(yval), (*C.mpz_t)(unsafe.Pointer(val))))
}

// ValGetMpq is the go version of yices_val_get_mpq.
func ValGetMpq(model ModelT, yval *YvalT, val *MpqT) int32 {
	return int32(C.yices_val_get_mpqp(ymodel(model), (*C.yval_t)(yval), (*C.mpq_t)(unsafe.Pointer(val))))
}

/*
//iam: not gonna assume mcsat
#ifdef LIBPOLY_VERSION
__YICES_DLLSPEC__ extern int32_t yices_val_get_algebraic_number(model_t *mdl, const yval_t *v, lp_algebraic_number_t *a);
#endif
*/

// ValGetBv is the go version of yices_val_get_bv.
func ValGetBv(model ModelT, yval *YvalT, val []int32) int32 {
	return int32(C.yices_val_get_bv(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(&val[0])))
}

// ValGetScalar is the go version of yices_val_get_scalar.
func ValGetScalar(model ModelT, yval *YvalT, val *int32, tau *TypeT) int32 {
	return int32(C.yices_val_get_scalar(ymodel(model), (*C.yval_t)(yval), (*C.int32_t)(val), (*C.type_t)(tau)))
}

// ValExpandTuple is the go version of yices_val_expand_tuple.
func ValExpandTuple(model ModelT, yval *YvalT, child []YvalT) int32 {
	return int32(C.yices_val_expand_tuple(ymodel(model), (*C.yval_t)(yval), (*C.yval_t)(&child[0])))
}

// ValExpandFunction is the go version of yices_val_expand_function.
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

// ValExpandMapping is the go version of yices_val_mapping_arity.
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

// FormulaTrueInModel is the go version of yices_formula_true_in_model.
func FormulaTrueInModel(model ModelT, t TermT) int32 {
	return int32(C.yices_formula_true_in_model(ymodel(model), C.term_t(t)))
}

// FormulasTrueInModel is the go version of yices_formulas_true_in_model.
func FormulasTrueInModel(model ModelT, t []TermT) int32 {
	tcount := C.uint32_t(len(t))
	return int32(C.yices_formulas_true_in_model(ymodel(model), tcount, (*C.term_t)(&t[0])))
}

/*
 * CONVERSION OF VALUES TO CONSTANT TERMS
 */

// GetValueAsTerm is the go version of yices_get_value_as_term.
func GetValueAsTerm(model ModelT, t TermT) TermT {
	return TermT(C.yices_get_value_as_term(ymodel(model), C.term_t(t)))
}

// TermArrayValue is the go version of yices_term_array_value.
func TermArrayValue(model ModelT, a []TermT, b []TermT) int32 {
	tcount := C.uint32_t(len(a))
	return int32(C.yices_term_array_value(ymodel(model), tcount, (*C.term_t)(&a[0]), (*C.term_t)(&b[0])))
}

/*
 * IMPLICANTS
 */

// ImplicantForFormula is the go version of yices_implicant_for_formula.
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

// ImplicantForFormulas is the go version of yices_implicant_for_formulas.
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

// GenModeT is the go analog of the enum yices_gen_mode_t defined in yices_types.h.
type GenModeT int32

// These are the go analogs of the elements of the enum yices_gen_mode_t defined in yices_types.h.
const (
	GenDefault GenModeT = iota
	GenBySubst
	GenByProj
)

// GeneralizeModel is the go version of yices_generalize_model.
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

// GeneralizeModelArray is the go version of yices_generalize_model_array.
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

// ModelTermSupport is the go version of yices_model_term_support.
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

// ModelTermArraySupport is the go version of yices_model_term_array_support.
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

// HasDelegate is the go version of yices_has_delegate.
// Since 2.6.2
func HasDelegate(delegate string) bool {
	cdelegate := C.CString(delegate)
	defer C.free(unsafe.Pointer(cdelegate))
	has := C.yices_has_delegate(cdelegate)
	return has == 1
}

// CheckFormula is the go version of yices_check_formula.
// Since 2.6.2
func CheckFormula(t TermT, logic string, delegate string, model *ModelT) (status SmtStatusT) {
	clogic := C.CString(logic)
	defer C.free(unsafe.Pointer(clogic))
	cdelegate := C.CString(delegate)
	defer C.free(unsafe.Pointer(cdelegate))
	var cstatus C.smt_status_t
	var cmodel *C.model_t
	if model != nil {
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

// CheckFormulas is the go version of yices_check_formulas.
// Since 2.6.2
func CheckFormulas(t []TermT, logic string, delegate string, model *ModelT) (status SmtStatusT) {
	count := C.uint32_t(len(t))
	clogic := C.CString(logic)
	defer C.free(unsafe.Pointer(clogic))
	cdelegate := C.CString(delegate)
	defer C.free(unsafe.Pointer(cdelegate))
	var cstatus C.smt_status_t
	var cmodel *C.model_t
	if model != nil {
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

// ExportFormulaToDimacs is the go version of yices_export_formula_to_dimacs.
// Since 2.6.2
func ExportFormulaToDimacs(t TermT, filename string, simplify bool, status *SmtStatusT) (errcode int32) {
	path := C.CString(filename)
	defer C.free(unsafe.Pointer(path))
	var csimplify C.int
	if simplify {
		csimplify = 1
	} else {
		csimplify = 0
	}
	var cstatus C.smt_status_t
	errcode = int32(C.yices_export_formula_to_dimacs(C.term_t(t), path, csimplify, &cstatus))
	if errcode == 0 {
		*status = SmtStatusT(cstatus)
	}
	return
}

// ExportFormulasToDimacs is the go version of yices_export_formulas_to_dimacs.
// Since 2.6.2
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
	var cstatus C.smt_status_t
	errcode = int32(C.yices_export_formulas_to_dimacs((*C.term_t)(&t[0]), count, path, csimplify, &cstatus))
	if errcode == 0 {
		*status = SmtStatusT(cstatus)
	}
	return
}

/**********************
 *  PRETTY PRINTING   *
 **********************/

// PpType is the go version of yices_pp_type_fd
func PpType(file *os.File, tau TypeT, width uint32, height uint32, offset uint32) int32 {
	return int32(C.yices_pp_type_fd(C.int(file.Fd()), C.type_t(tau), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset)))
}

// PpTerm is the go version of yices_pp_term_fd
func PpTerm(file *os.File, t TermT, width uint32, height uint32, offset uint32) int32 {
	return int32(C.yices_pp_term_fd(C.int(file.Fd()), C.term_t(t), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset)))
}

// PpTermArray is the go version of yices_pp_term_array_fd
func PpTermArray(file *os.File, t []TermT, width uint32, height uint32, offset uint32, horiz int32) int32 {
	tcount := C.uint32_t(len(t))
	return int32(C.yices_pp_term_array_fd(C.int(file.Fd()), tcount, (*C.term_t)(&t[0]), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset), C.int32_t(horiz)))
}

// PrintModel is the go version of yices_print_model_fd
func PrintModel(file *os.File, model ModelT) int32 {
	return int32(C.yices_print_model_fd(C.int(file.Fd()), ymodel(model)))
}

// PrintTermValues is the go version of yices_print_term_values_fd
// Since 2.6.2
func PrintTermValues(file *os.File, model ModelT, t []TermT) int32 {
	tcount := C.uint32_t(len(t))
	return int32(C.yices_print_term_values_fd(C.int(file.Fd()), ymodel(model), tcount, (*C.term_t)(&t[0])))
}

// PpTermValues is the go version of yices_pp_term_values_fd
// Since 2.6.2
func PpTermValues(file *os.File, model ModelT, t []TermT, width uint32, height uint32, offset uint32) int32 {
	tcount := C.uint32_t(len(t))
	return int32(C.yices_pp_term_values_fd(C.int(file.Fd()), ymodel(model), tcount, (*C.term_t)(&t[0]), C.uint32_t(width), C.uint32_t(height), C.uint32_t(offset)))
}

// PpModel is the go version of yices_pp_model_fd(
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
