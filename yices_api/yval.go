package yices2

// YvalTagT is the go version of the yval_tag_t enum
type YvalTagT int32

// The corresponding element of the YvalTagT enum
const (
	YvalUnknown YvalTagT = iota
	YvalBool
	YvalRational
	YvalAlgebraic
	YvalBv
	YvalScalar
	YvalTuple
	YvalFunction
	YvalMapping
)
