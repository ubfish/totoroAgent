/**
 totoraAgent project main.go
 author:feng
 since:2020-01-08
**/
package main

import (
	"totoroAgent"
	"log"
)

func main() {
	// 加载配置文件,没有配置文件加载默认配置
	config, s, err := totoroAgent.InitConfig()
	totoroAgent.ClientLog(config.Version)
	if err != nil {
		totoroAgent.ClientLog("load config error ", err)
		log.Fatal("Can't load config")
	}
	totoroAgent.AppConfig = config
	totoroAgent.AppSignal = s

	// 初始化日志环境
	totoroAgent.LogManager(config.LogPath)

	totoroAgent.ClientLog("config content ", config)
	totoroAgent.ClientLog("signal ", s)
	// 初始化运行环境
	totoroAgent.InitEnviron()

	//开启监听端口
	totoroAgent.StartListen(config.ListenPort)
	totoroAgent.SignalHandler()
}
