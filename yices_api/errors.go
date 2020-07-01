package yices2

// These are the go versions of the yices_types.h error_code (code <= 100 constants: Errors in type or term construction)
const (
	NoError ErrorCodeT = iota
	ErrorInvalidType
	ErrorInvalidTerm
	ErrorInvalidConstantIndex
	ErrorInvalidVarIndex // Not used anymore
	ErrorInvalidTupleIndex
	ErrorInvalidRationalFormat
	ErrorInvalidFloatFormat
	ErrorInvalidBvbinFormat
	ErrorInvalidBvhexFormat
	ErrorInvalidBitshift
	ErrorInvalidBvextract
	ErrorInvalidBitextract // added 2014/02/17
	ErrorTooManyArguments
	ErrorTooManyVars
	ErrorMaxBvsizeExceeded
	ErrorDegreeOverflow
	ErrorDivisionByZero
	ErrorPosIntRequired
	ErrorNonnegIntRequired
	ErrorScalarOrUtypeRequired
	ErrorFunctionRequired
	ErrorTupleRequired
	ErrorVariableRequired
	ErrorArithtermRequired
	ErrorBitvectorRequired
	ErrorScalarTermRequired
	ErrorWrongNumberOfArguments
	ErrorTypeMismatch
	ErrorIncompatibleTypes
	ErrorDuplicateVariable
	ErrorIncompatibleBvsizes
	ErrorEmptyBitvector
	ErrorArithconstantRequired // added 2013/01/23
	ErrorInvalidMacro          // added 2013/03/31
	ErrorTooManyMacroParams    // added 2013/03/31
	ErrorTypeVarRequired       // added 2013/03/31
	ErrorDuplicateTypeVar      // added 2013/03/31
	ErrorBvtypeRequired        // added 2013/05/27
	ErrorBadTermDecref         // added 2013/10/03
	ErrorBadTypeDecref         // added 2013/10/03
	ErrorInvalidTypeOp         // added 2014/12/03
	ErrorInvalidTermOp         // added 2014/12/04
)

// These are the go versions of the yices_types.h error_code (100 < code  <= 300 constants: Parser errors)
const (
	/*
	 * Parser errors
	 */
	ErrorInvalidToken ErrorCodeT = 100 + iota
	ErrorSyntaxError
	ErrorUndefinedTypeName
	ErrorUndefinedTermName
	ErrorRedefinedTypeName
	ErrorRedefinedTermName
	ErrorDuplicateNameInScalar
	ErrorDuplicateVarName
	ErrorIntegerOverflow
	ErrorIntegerRequired
	ErrorRationalRequired
	ErrorSymbolRequired
	ErrorTypeRequired
	ErrorNonConstantDivisor
	ErrorNegativeBvsize
	ErrorInvalidBvconstant
	ErrorTypeMismatchInDef
	ErrorArithError
	ErrorBvarithError
)

// These are the go versions of the yices_types.h error_code (300 < code  <= 400 constants: Errors in assertion processing.)
const (
	/*
	 * These codes mean that the context as configured
	 * cannot process the assertions.
	 */
	ErrorCtxFreeVarInFormula ErrorCodeT = 300 + iota
	ErrorCtxLogicNotSupported
	ErrorCtxUfNotSupported
	ErrorCtxArithNotSupported
	ErrorCtxBvNotSupported
	ErrorCtxArraysNotSupported
	ErrorCtxQuantifiersNotSupported
	ErrorCtxLambdasNotSupported
	ErrorCtxNonlinearArithNotSupported
	ErrorCtxFormulaNotIdl
	ErrorCtxFormulaNotRdl
	ErrorCtxTooManyArithVars
	ErrorCtxTooManyArithAtoms
	ErrorCtxTooManyBvVars
	ErrorCtxTooManyBvAtoms
	ErrorCtxArithSolverException
	ErrorCtxBvSolverException
	ErrorCtxArraySolverException
	ErrorCtxScalarNotSupported // added 2015/03/26
	ErrorCtxTupleNotSupported  // added 2015/03/26
	ErrorCtxUtypeNotSupported  // added 2015/03/26
)

// These are the go versions of the yices_types.h error_code (400 < code  <= 500 constants: Error codes for other operations.)
const (
	ErrorCtxInvalidOperation ErrorCodeT = 400 + iota
	ErrorCtxOperationNotSupported
	ErrorCtxUnknownDelegate      = 420 + iota // Since 2.6.2.
	ErrorCtxDelegateNotAvailable              // Since 2.6.2.
)

// These are the go versions of the yices_types.h error_code (500 < code  <= 600 constants: Errors in context configurations and search parameter settings.)
const (
	ErrorCtxInvalidConfig ErrorCodeT = 500 + iota
	ErrorCtxUnknownParameter
	ErrorCtxInvalidParameterValue
	ErrorCtxUnknownLogic
)

// These are the go versions of the yices_types.h error_code (600 < code  <= 700 constants: Error codes for model queries.)
const (
	ErrorEvalUnknownTerm ErrorCodeT = 600 + iota
	ErrorEvalFreevarInTerm
	ErrorEvalQuantifier
	ErrorEvalLambda
	ErrorEvalOverflow
	ErrorEvalFailed
	ErrorEvalConversionFailed
	ErrorEvalNoImplicant
	ErrorEvalNotSupported
)

// These are the go versions of the yices_types.h error_code (700 < code  <= 800 constants: Error codes for model construction.)
const (
	ErrorMdlUnintRequired ErrorCodeT = 700 + iota
	ErrorMdlConstantRequired
	ErrorMdlDuplicateVar
	ErrorMdlFtypeNotAllowed
	ErrorMdlConstructionFailed
)

// These are the go versions of the yices_types.h error_code (800 < code  <= 900 constants: Error codes in DAG/node queries.)
const (
	ErrorYvalInvalidOp ErrorCodeT = 800 + iota
	ErrorYvalOverflow
	ErrorYvalNotSupported
)

// These are the go versions of the yices_types.h error_code (900 < code  <= 1000 constants: Error codes for model generalization.)
const (
	ErrorMdlGenTypeNotSupported ErrorCodeT = 900 + iota
	ErrorMdlGenNonlinear
	ErrorMdlGenFailed
)

// These are the go versions of the yices_types.h error_code (1000 < code  <= 9000 constants: MCSAT error codes.)
const (
	ErrorMcsatErrorUnsupportedTheory ErrorCodeT = 1000 + iota
)

// These are the go versions of the yices_types.h error_code (9000 < code  <= 9999 constants: Input/output and system errors.)
const (
	ErrorOutputError ErrorCodeT = 9000 + iota
)

// These are the go versions of the yices_types.h error_code (9999< code constants: Catch-all code for anything else.)
const (
	/*
	 * This is a symptom that a bug has been found.
	 */
	ErrorInternalException ErrorCodeT = 9999 + iota
)
