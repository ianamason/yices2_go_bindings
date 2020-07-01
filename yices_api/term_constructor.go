package yices2

// TermConstructorT is the analog of the enum term_constructor_t defined in yices_types.h
type TermConstructorT int32

// These are the elements of TermConstructorT, the analog of the enum elements of term_constructor_t defined in yices_types.h
const (
	TrmCnstrConstructorError TermConstructorT = iota - 1 // to report an error

	// atomic terms
	TrmCnstrBoolConstant      TermConstructorT = iota // boolean constant
	TrmCnstrArithConstant                             // rational constant
	TrmCnstrBvConstant                                // bitvector constant
	TrmCnstrScalarConstant                            // constant of uninterpreted/scalar
	TrmCnstrVariable                                  // variable in quantifiers
	TrmCnstrUninterpretedTerm                         // (i.e. global variables can't be bound)

	// composite terms
	TrmCnstrIteTerm      // if-then-else
	TrmCnstrAppTerm      // application of an uninterpreted function
	TrmCnstrUpdateTerm   // function update
	TrmCnstrTupleTerm    // tuple constructor
	TrmCnstrEqTerm       // equality
	TrmCnstrDistinctTerm // distinct t_1 ... t_n
	TrmCnstrForallTerm   // quantifier
	TrmCnstrLambdaTerm   // lambda
	TrmCnstrNotTerm      // (not t)
	TrmCnstrOrTerm       // n-ary OR
	TrmCnstrXorTerm      // n-ary XOR

	TrmCnstrBvArray       // array of boolean terms
	TrmCnstrBvDiv         // unsigned division
	TrmCnstrBvRem         // unsigned remainder
	TrmCnstrBvSdiv        // signed division
	TrmCnstrBvSrem        // remainder in signed division (rounding to 0)
	TrmCnstrBvSmod        // remainder in signed division (rounding to -infinity)
	TrmCnstrBvShl         // shift left (padding with 0)
	TrmCnstrBvLshr        // logical shift right (padding with 0)
	TrmCnstrBvAshr        // arithmetic shift right (padding with sign bit)
	TrmCnstrBvGeAtom      // unsigned comparison: (t1 >= t2)
	TrmCnstrBvSgeAtom     // signed comparison (t1 >= t2)
	TrmCnstrArithGeAtom   // atom (t1 >= t2) for arithmetic terms: t2 is always 0
	TrmCnstrArithRootAtom // atom (0 <= k <= root_count(p)) && (x r root(p, k)) for r in <, <=, ==, !=, >, >=

	TrmCnstrAbs         // absolute value
	TrmCnstrCeil        // ceil
	TrmCnstrFloor       // floor
	TrmCnstrRdiv        // real division (as in x/y)
	TrmCnstrIdiv        // integer division
	TrmCnstrImod        // modulo
	TrmCnstrIsIntAtom   // integrality test: (is-int t)
	TrmCnstrDividesAtom // divisibility test: (divides t1 t2)

	// projections
	TrmCnstrSelectTerm // tuple projection
	BTrmCnstritTerm    // bit-select: extract the i-th bit of a bitvector

	// sums
	TrmCnstrBvSum    // sum of pairs a * t where a is a bitvector constant (and t is a bitvector term)
	ATrmCnstrrithSum // sum of pairs a * t where a is a rational (and t is an arithmetic term)

	// products
	TrmCnstrPowerProduct // power products: (t1^d1 * ... * t_n^d_n)

)
