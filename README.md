# totoroAgent Go语言编写的适用于Linux和Windows的多线程agent软件

## 项目简介
使用Go语言开发的多线程agent
1.使用http方式调用
2.支持加密传输方式
3.支持调用底层命令与脚本
4.支持同步与异步调用方式
5.支持心跳功能
6.支持windows和linux多操作系统，

## 代码架构：
totoroAgent.go  //入口方法，初始化程序
/totoroAgent    //处理程序，包含

## 编译

### LINUX编译

安装好golang后执行build.sh
编译成功后，执行/bin目录下的totoroAgent运行


### Windows编译

配置好golang和goPath后
执行
'''
go install totoroAgent totoroAgent.go
'''
编译完成后，执行totoroAgent.exe


## 运行
1.直接执行totoroAgent，使用默认配置（默认端口10099）
2.自定义配置(使用json配置启动参数)
'''
./totoroAgent -c config.json
'''

## 使用
启动应用后，访问http://{agent地址}:10099/version，返回版本号即成功

### 执行命令
POST 访问http://{agent地址}:10099/exec, body内容为
'''
{
  "actionType" : "exec",   
  "cmd" : "ls /export"
}
'''
返回内容
'''
Command exit code: 0
code
go
logs
var
'''
第一行返回命令执行结果，后面返回命令执行的内容

### 异步命令
