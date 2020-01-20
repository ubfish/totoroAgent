/**
 totoroAgent project log.go
 author:feng
 since:2020-01-08
**/
package agent

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var agentLog *log.Logger

func SaveHistoryLog(logPath string) error {
	dir, file := filepath.Split(logPath)
	newName := filepath.Join(dir, file+"."+time.Now().Format("2006-01-02"))
	oldNameResult := CheckFileIsExist(logPath)
	if oldNameResult == false {
		err := CreateDir(logPath)
		if err != nil {
			return err
		}
	}
	newNameResult := CheckFileIsExist(newName)
	if newNameResult == false {
		err := os.Rename(logPath, newName)
		if err != nil {
			ClientLog("rename error:", err)
		}
	}
	return nil
}

func InitLog(logPath string) (file *os.File) {
	log.SetFlags(3)
	logfile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		ClientLog("initLog error:", err)
		return nil
	}
	clientLogger := log.New(logfile, "", log.Ldate|log.Ltime)
	agentLog = clientLogger
	ClientLog("initLog success!")
	return logfile
}

func TickerForLog(logFile *os.File, logPath string) {
	timer1 := time.NewTicker(24 * time.Hour)
	var flag bool = false
	var newLogFile *os.File
	for {
		select {
		case <-timer1.C:
			if flag == false {
				if logFile == nil {
					ClientLog("logFile is nil")
					return
				}
				logFile.Close()
				flag = true
			} else {
				newLogFile.Close()
			}
			err := SaveHistoryLog(logPath)
			if err != nil {
				ClientLog("saveHistoryLog error:", err)
				return
			}
			newLogFile = InitLog(logPath)
			if newLogFile == nil {
				ClientLog("newLogFile is nil")
				return
			}
			ClientLog("ticker for save historyLog success!")
		}
	}
}

func CreateDir(logPath string) (err error) {
	dir := filepath.Dir(logPath)
	err = os.MkdirAll(dir, os.ModePerm)
	return err
}

func LogManager(logPath string) {
	err := CreateDir(logPath)
	if err != nil {
		ClientLog("CreateDir error:", err)
		return
	}
	file := InitLog(logPath)
	if file == nil {
		ClientLog("LogManager initLog error of the file nil")
		return
	}
	go TickerForLog(file, logPath)
}

func CheckFileIsExist(logPath string) (result bool) {
	f, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return true
	}
	defer f.Close()
	return true
}

func ClientLog(result ...interface{}) {
	_, _, line, _ := runtime.Caller(1)
	if agentLog != nil {
		agentLog.Println(line, result)
	}
	log.Println(line, result)
}
