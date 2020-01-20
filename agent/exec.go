/**
 totoroAgent project exec.go
 author:feng
 since:2020-01-08
**/
package agent

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// 执行脚本命令
func Exeshell(task *TaskDetail) {
	ClientLog(task.CmdType, "exec started:")
	resultLog := make(chan string)
	task.ResultCode = 1
	go func() {
		if task.Cmd == "" {
			ClientLog("no command find")
			resultLog <- "success"
		} else {
			comd := strings.Fields(task.Cmd)
			ClientLog("cmd:", task.Cmd)
			result, err := exec.Command(comd[0], comd[1:]...).CombinedOutput()
			ClientLog("cmd execting:")
			task.Status = 0 //命令的返回值，服务端默认为0标识成功
			if err != nil {
				task.ResultCode = 0
				if exitErr, ok := err.(*exec.ExitError); ok {
					if s, ok := exitErr.Sys().(syscall.WaitStatus); ok {
						task.Status = int(s.ExitStatus())
					} else {
						task.Status = -1 // 没有取到返回值
					}
				}
			}
			re := "Command exit code: " + strconv.Itoa(task.Status) + "\n" + string(result)
			resultLog <- re
			ClientLog("shell log:", re)
		}
	}()
	select {
	case res := <-resultLog:
		task.ResultInfo += res
	case <-time.After(time.Second * 300):
		task.ResultInfo += " timeout"
		task.ResultCode = 0
		ClientLog("shell is timeout")
	}
	ClientLog(task.CmdType, task.Cmd, "exec finished")
}
