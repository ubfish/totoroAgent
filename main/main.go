/**
 totoraAgent project main.go
 author:feng
 since:2020-01-08
**/
package main

import (
	"agent"
	"log"
)

func main() {
	// 加载配置文件,没有配置文件加载默认配置
	config, s, err := agent.InitConfig()
	agent.ClientLog(config.Version)
	if err != nil {
		agent.ClientLog("load config error ", err)
		log.Fatal("Can't load config")
	}
	agent.AppConfig = config
	agent.AppSignal = s

	// 初始化日志环境
	agent.LogManager(config.LogPath)

	agent.ClientLog("config content ", config)
	agent.ClientLog("signal ", s)
	// 初始化运行环境
	agent.InitEnviron()

	//开启监听端口
	agent.StartListen(config.ListenPort)
	agent.SignalHandler()
}
