/**
 totoroAgent project exechandler.go
 author:feng
 since:2020-01-08
**/
package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//监听服务器需要执行命令的请求
func ExecHandler(w http.ResponseWriter, req *http.Request) {
	ClientLog("start execution command")
	body, err := ReadBody(req)
	if err != nil {
		ClientLog("read error:", err)
		w.Write([]byte("read data error" + err.Error()))
	} else {
		var taskDetail TaskDetail
		err := json.Unmarshal(body, &taskDetail)
		if err != nil {
			ClientLog("taskDetail unmarshal error:", err)
		} else {
			ClientLog("taskDdetail:" + taskDetail.Cmd)
			Exeshell(&taskDetail)
			//exeResult := MarshalTask(&taskDetail)
			fmt.Fprintln(w, string(taskDetail.ResultInfo))
			ClientLog("execution command success")
		}
	}
}
