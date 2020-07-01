package yices2

// SmtStatusT is the analog of smt_status_t defined in yices_types.h
type SmtStatusT int32

// These are the analogs to the elements of smt_status_t defined in yices_types.h
const (
	StatusIdle SmtStatusT = iota
	StatusSearching
	StatusUnknown
	StatusSat
	StatusUnsat
	StatusInterrupted
	StatusError
)
