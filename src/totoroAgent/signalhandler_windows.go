/**
 totoroAgent project aignalhandler_windows.go
 author:feng
 since:2020-01-08
**/
package totoroAgent

import (
	"os"
)

func SignalHandler() {
	for {
	}
}

func SigProcessRestart(p *os.Process) {
	_ = p
}

func IsProcessExist(p *os.Process) (result bool) {
	return true
}
