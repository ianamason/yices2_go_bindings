package yices2

type Smt_status_t int32

const (
  STATUS_IDLE Smt_status_t = iota
  STATUS_SEARCHING
  STATUS_UNKNOWN
  STATUS_SAT
  STATUS_UNSAT
  STATUS_INTERRUPTED
  STATUS_ERROR
)
