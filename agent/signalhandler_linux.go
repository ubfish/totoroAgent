/**
 totoroAgent project aignalhandler_linux.go
 author:feng
 since:2020-01-08
**/
package agent

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SignalHandler() {
	for {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGKILL,
			syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)
		sig := <-ch
		ClientLog("Signal received: ", sig)
		switch sig {
		default:
			ClientLog("get sig=", sig)
		case syscall.SIGHUP:
			ClientLog("终端退出，进程不退出")
		case syscall.SIGUSR1:
			ClientLog("usr1 do nothing")
		case syscall.SIGUSR2:
			processRestart()
		case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL:
			exitProcess()
		}

	}
}

func exitProcess() {
	os.Remove(PidFile)
	ClientLog("Remove pid file and exit 0")
	os.Exit(0)
}

func processRestart() {
	ClearTail()
	ClientLog("停掉老的监听")
	ChanNetClose <- "net close"
	err := NetLisener.Close()
	if err != nil {
		ClientLog("old listen close error")
	}
	ClientLog("启用新的代码 exePath", ExePath)
	// 修改老的PID文件为新的
	os.Rename(PidFile, OldPidFile)
	// 启动新进程
	ps, err := os.StartProcess(ExePath, os.Args,
		&os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	if err != nil {
		ClientLog("start new process error", err)
		os.Rename(OldPidFile, PidFile)
	} else {
		ClientLog("restart ok", ps.Pid)
		//退出老进程，删除old pid文件
		os.Remove(OldPidFile)
		// 当前没有任务处理就直接退出，有任务就等待任务结束后退出
		time.Sleep(30 * time.Second)
		size := len(MapTasks)
		ClientLog(os.Getpid(), " have ", size, "task is process")
		if size == 0 {
			os.Exit(0)
		} else {
			waitTaskFinish()
		}
	}
}

func waitTaskFinish() {
	count := 0
	for {
		time.Sleep(30 * time.Second)
		size := len(MapTasks)
		if size == 0 {
			ClientLog("pid task finished exit", os.Getpid())
			os.Exit(0)
		} else {
			ClientLog(os.Getpid(), " have ", size, "task is process ", time.Now().String())
		}
		count++
		// 操作30*10,300秒，5分钟后就强制退出
		if count > 10 {
			os.Exit(0)
		}
	}
}

func SigProcessRestart(p *os.Process) {
	p.Signal(syscall.SIGUSR2)
}

func IsProcessExist(p *os.Process) (result bool) {
	err := p.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}
	return true
}
