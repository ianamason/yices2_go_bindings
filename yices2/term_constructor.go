package yices2

type Term_constructor_t int32


const (
  CONSTRUCTOR_ERROR Term_constructor_t = iota - 1 // to report an error

  // atomic terms
  BOOL_CONSTANT  Term_constructor_t = iota     // boolean constant
  ARITH_CONSTANT      // rational constant
  BV_CONSTANT         // bitvector constant
  SCALAR_CONSTANT     // constant of uninterpreted/scalar
  VARIABLE            // variable in quantifiers
  UNINTERPRETED_TERM  // (i.e. global variables can't be bound)

  // composite terms
  ITE_TERM            // if-then-else
  APP_TERM            // application of an uninterpreted function
  UPDATE_TERM         // function update
  TUPLE_TERM          // tuple constructor
  EQ_TERM             // equality
  DISTINCT_TERM       // distinct t_1 ... t_n
  FORALL_TERM         // quantifier
  LAMBDA_TERM         // lambda
  NOT_TERM            // (not t)
  OR_TERM             // n-ary OR
  XOR_TERM            // n-ary XOR

  BV_ARRAY            // array of boolean terms
  BV_DIV              // unsigned division
  BV_REM              // unsigned remainder
  BV_SDIV             // signed division
  BV_SREM             // remainder in signed division (rounding to 0)
  BV_SMOD             // remainder in signed division (rounding to -infinity)
  BV_SHL              // shift left (padding with 0)
  BV_LSHR             // logical shift right (padding with 0)
  BV_ASHR             // arithmetic shift right (padding with sign bit)
  BV_GE_ATOM          // unsigned comparison: (t1 >= t2)
  BV_SGE_ATOM         // signed comparison (t1 >= t2)
  ARITH_GE_ATOM       // atom (t1 >= t2) for arithmetic terms: t2 is always 0
  ARITH_ROOT_ATOM     // atom (0 <= k <= root_count(p)) && (x r root(p, k)) for r in <, <=, ==, !=, >, >=


  ABS                 // absolute value
  CEIL                // ceil
  FLOOR               // floor
  RDIV                // real division (as in x/y)
  IDIV                // integer division
  IMOD                // modulo
  IS_INT_ATOM         // integrality test: (is-int t)
  DIVIDES_ATOM        // divisibility test: (divides t1 t2)
  
  // projections
  SELECT_TERM         // tuple projection
  BIT_TERM            // bit-select: extract the i-th bit of a bitvector

  // sums
  BV_SUM              // sum of pairs a * t where a is a bitvector constant (and t is a bitvector term)
  ARITH_SUM           // sum of pairs a * t where a is a rational (and t is an arithmetic term)

  // products
  POWER_PRODUCT        // power products: (t1^d1 * ... * t_n^d_n)

)
