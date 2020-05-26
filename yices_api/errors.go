package yices2

const (
	NoError ErrorCodeT = iota
	/*
	 * Errors in type or term construction
	 */
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

const (
	/*
	 * Errors in assertion processing.
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

const (
	/*
	 * Error codes for other operations
	 */
	ErrorCtxInvalidOperation ErrorCodeT = 400 + iota
	ErrorCtxOperationNotSupported
)

const (
	/*
	 * Errors in context configurations and search parameter settings
	 */
	ErrorCtxInvalidConfig ErrorCodeT = 500 + iota
	ErrorCtxUnknownParameter
	ErrorCtxInvalidParameterValue
	ErrorCtxUnknownLogic
)

const (
	/*
	 * Error codes for model queries
	 */
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

const (
	/*
	 * Error codes for model construction
	 */
	ErrorMdlUnintRequired ErrorCodeT = 700 + iota
	ErrorMdlConstantRequired
	ErrorMdlDuplicateVar
	ErrorMdlFtypeNotAllowed
	ErrorMdlConstructionFailed
)

const (
	/*
	 * Error codes in DAG/node queries
	 */
	ErrorYvalInvalidOp ErrorCodeT = 800 + iota
	ErrorYvalOverflow
	ErrorYvalNotSupported
)

const (
	/*
	 * Error codes for model generalization
	 */
	ErrorMdlGenTypeNotSupported ErrorCodeT = 900 + iota
	ErrorMdlGenNonlinear
	ErrorMdlGenFailed
)

const (
	/*
	 * MCSAT error codes
	 */
	ErrorMcsatErrorUnsupportedTheory ErrorCodeT = 1000 + iota
)

const (
	/*
	 * Input/output and system errors
	 */
	ErrorOutputError ErrorCodeT = 9000 + iota
)

const (
	/*
	 * Catch-all code for anything else.
	 * This is a symptom that a bug has been found.
	 */
	ErrorInternalException ErrorCodeT = 9999 + iota
)
