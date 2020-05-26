package yices2

type SmtStatusT int32

const (
	StatusIdle SmtStatusT = iota
	StatusSearching
	StatusUnknown
	StatusSat
	StatusUnsat
	StatusInterrupted
	StatusError
)
