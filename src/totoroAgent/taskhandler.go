/**
 totoroAgent project taskhandler.go
 author:feng
 since:2020-01-08
**/
package totoroAgent

import (
	"encoding/json"
	"net/http"
)

var MapTasks map[string]*[]TaskDetail = make(map[string]*[]TaskDetail)

func TaskHandler(w http.ResponseWriter, req *http.Request) {
	ClientLog("start process request ")
	body, err := ReadBody(req)
	if err != nil {
		ClientLog("read error: ", err)
		w.Write([]byte("read data error" + err.Error()))
	} else {
		var t []TaskDetail
		ClientLog("content:", string(body))
		err := json.Unmarshal(body, &t)
		if err != nil {
			ClientLog("Unmarshal error: ", err)
			w.Write([]byte(err.Error()))
		} else if len(t) > 0 {
			ClientLog("tasks:", t)
			// 把任务保存下来
			MapTasks[t[0].TaskId] = &t
			go DoTasks(&t)
			ClientLog("content:", t)
			w.Write([]byte("ok"))
		}
	}
	ClientLog("end process ")
}
