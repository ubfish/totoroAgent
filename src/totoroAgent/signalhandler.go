/**
 totoroAgent project signalhanlder.go
 author:feng
 since:2020-01-08
**/
package totoroAgent

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Version string
var PidFile string
var OldPidFile string

var OldPid int
var ExePath string

var AppSignal string

func InitEnviron() {
	exeFile, _ := filepath.Abs(os.Args[0])
	ExePath = exeFile
	makePidPath()
	makeEvnPath()
	ManageProcess()
}

func makePidPath() {
	curPidName := "totoroAgent.pid"
	oldPidName := "totoroAgent.pid.oldbin"
	if strings.LastIndex(AppConfig.PidPath, "/") != (len(AppConfig.PidPath) - 1) {
		curPidName = "/" + curPidName
		oldPidName = "/" + oldPidName
	}
	PidFile = AppConfig.PidPath + curPidName
	OldPidFile = AppConfig.PidPath + oldPidName
	ClientLog(PidFile, OldPidFile)
}

func makeEvnPath() {
	path := os.Getenv("PATH")
	ClientLog("ExePATH: ", ExePath)
	exeDir := filepath.Dir(ExePath)
	if !strings.Contains(path, exeDir) {
		path = exeDir + ":" + path
		os.Setenv("PATH", path)
	}
	ClientLog("new env path :", path)
}

func RemoveAndContinue(pid int) {
	ClientLog("can't find the pid process error remove pid file", pid)
	os.Remove(PidFile)
	ClientLog("continue new process ", os.Getpid())
	CreatePidFile()
}

func ManageProcess() {
	pidExist := CheckFileIsExist(PidFile)
	oldExist := CheckFileIsExist(OldPidFile)

	if !pidExist && !oldExist {
		CreatePidFile()
		ClientLog("ceate pid file", PidFile, " pid ", os.Getpid())
	} else if pidExist && !oldExist { // 如果进程存在，老进程号不存在
		pid, err := ReadPidFile(PidFile)
		if err != nil {
			ClientLog("read pid file error", err)
			os.Exit(6)
		}
		p, err := os.FindProcess(pid)
		if err != nil {
			RemoveAndContinue(pid)
		} else {
			re := IsProcessExist(p)
			if re == false {
				RemoveAndContinue(pid)
			} else {
				if AppSignal == "restart" {
					ClientLog("find the process", pid)
					SigProcessRestart(p)
					ClientLog("exit restart process", os.Getpid())
					os.Exit(0)
				} else {
					ClientLog("process is running exit now", os.Getpid())
					os.Exit(0)
				}
			}
		}

	} else if !pidExist && oldExist {
		ClientLog("signal start process ", os.Getpid())
		CreatePidFile()
	}

}

func ManagePid() (second bool) {
	second = false
	pidExist := CheckFileIsExist(PidFile)
	oldExist := CheckFileIsExist(OldPidFile)
	if pidExist && oldExist {
		ClientLog("pid oldpid file exist. this is the wrong operate. exit code 8")
		os.Exit(8)
	}
	if pidExist {
		err := os.Rename(PidFile, OldPidFile)
		if err != nil {
			ClientLog("rename error", err, "exit code 10")
			os.Exit(10)
		}
		second = true
	}
	return second
}

func CreatePidFile() {
	CreateDir(PidFile)
	newFile, err := os.Create(PidFile)
	if err != nil {
		ClientLog("creat pid file error exit code 10")
		os.Exit(9)
	}
	defer newFile.Close()
	newFile.WriteString(strconv.Itoa(os.Getpid()))
}

func ReadPidFile(filePath string) (pid int, er error) {
	file, err := os.Open(filePath)
	pid = 0
	if err != nil {
		ClientLog("ReadPidFile open pid error:", file, err)
		return pid, err
	}
	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		ClientLog("ReadPidFile read pid file error:", err)
		return pid, err
	}
	data = data[0:count]
	temp, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		ClientLog("ReadPidFile parse int error", err)
		return pid, err
	}
	pid = int(temp)
	return pid, nil
}
