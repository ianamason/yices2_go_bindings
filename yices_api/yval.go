package yices2

type YvalTagT int32

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
