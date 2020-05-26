package yices

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
	"os"
)

/*********************
 *  VERSION NUMBERS  *
 ********************/

// Version is the yices2 library version.
var Version string

// BuildArch is the yices2 library build architecture.
var BuildArch string

// BuildMode is the yices2 library build mode.
var BuildMode string

// BuildDate is the yices2 library build date.
var BuildDate string

// HasMcsat indicates if the yices2 library supports MCSAT.
var HasMcsat bool

// IsThreadSafe indicate if the yices2 library was built with thread safety enabled.
var IsThreadSafe bool

func init() {
	Version = yapi.Version()
	BuildArch = yapi.BuildArch()
	BuildMode = yapi.BuildMode()
	BuildDate = yapi.BuildDate()
	HasMcsat = (yapi.HasMcsat() == int32(1))
	IsThreadSafe = (yapi.IsThreadSafe() == int32(1))
}

/***************************************
 *  GLOBAL INITIALIZATION AND CLEANUP  *
 **************************************/

// Init initializes the internal yices2 library data structures.
func Init() { yapi.Init() }

// Exit cleans up the internal yices2 library data structures.
func Exit() { yapi.Exit() }

// Reset resets up the internal yices2 library data structures.
func Reset() { yapi.Reset() }

/*********************
 *  ERROR REPORTING  *
 ********************/

// Error returns the current yices error structure.
func Error() (yerror *yapi.YicesErrorT) {
	return yapi.YicesError()
}

// ErrorCode returns the most recent yices error code.
func ErrorCode() yapi.ErrorCodeT {
	return yapi.ErrorCode()
}

// ClearError clears the most recent error structure.
func ClearError() {
	yapi.ClearError()
}

// PrintError prints the most recent yices error.
func PrintError(f *os.File) int32 {
	return yapi.PrintError(f)
}

// ErrorString returns a string describing the most recent yices error.
func ErrorString() string {
	return yapi.ErrorString()
}

/***********************
 *  TYPE CONSTRUCTORS  *
 **********************/
