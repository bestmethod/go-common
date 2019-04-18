package gocommon

import (
	"os/exec"
	"syscall"
)

// check return code from exec function. Returns 666 if error existed but retcode could not be extracted
func check_exec_retcode(err error) int {
	if err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if ok == false {
			return 666
		}
		return exiterr.Sys().(syscall.WaitStatus).ExitStatus()
	}
	return 0
}
