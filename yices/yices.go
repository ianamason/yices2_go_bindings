package yices

import (
	yapi "github.com/ianamason/yices2_go_bindings/yices_api"
)

/*********************
 *  VERSION NUMBERS  *
 ********************/

var Version string

var Build_arch string

var Build_mode string

var Build_date string

var Has_mcsat bool

var Is_thread_safe bool

func init() {

	Version = yapi.Version()
	Build_arch = yapi.Build_arch()
	Build_mode = yapi.Build_mode()
	Build_date = yapi.Build_date()
	Has_mcsat = (yapi.Has_mcsat() == int32(1))
	Is_thread_safe = (yapi.Is_thread_safe() == int32(1))

}

/***************************************
 *  GLOBAL INITIALIZATION AND CLEANUP  *
 **************************************/

func Init() { yapi.Init() }

func Exit() { yapi.Exit() }

func Reset() { yapi.Reset() }
