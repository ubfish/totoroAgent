/**
 totoroAgent project task.go
 author:feng
 since:2020-01-08
**/
package agent

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mahonia"
	"net/http"
)

type TaskDetail struct {
	Id         int    `json:id`
	TaskId     string `json:taskId`
	ActionType string `json:actionType`
	CmdType    string `json:cmdType`
	Cmd        string `json:cmd`
	Status     int    `json:status`
	ResultCode int    `json:resultCode`
	ResultInfo string `json:resultInfo`
	Url        string `json:url`
}

/**
*
 */
func MarshalTask(task *TaskDetail) string {
	j, err := json.Marshal(task)
	if err != nil {
		ClientLog(err.Error())
		ClientLog("Use GBK encoding marshal task log")
		dec := mahonia.NewDecoder("GBK")
		task.ResultInfo = dec.ConvertString(task.ResultInfo)
		ClientLog("GBK convert log ", task.ResultInfo)
		j, err = json.Marshal(task)
		if err != nil {
			ClientLog("GBK convert log marshal fail ", err.Error())
			task.ResultInfo = err.Error()
			j, err = json.Marshal(task)
			if err != nil {
				ClientLog("json Marshal error ", err.Error())
			}
		}
	}
	return string(j)
}

/**
*  空-不做
*  执行脚本-调用命令
**/
func DoTask(task *TaskDetail) {
	actionType := task.ActionType
	ClientLog("task started and action type is ", actionType)
	if actionType == "" || actionType == NONE {
		ClientLog("ActionType is empty or NONE. Do nothing")
		postResult(task)
	} else if actionType == EXEC {
		Exeshell(task)
		postResult(task)
	} else {
		ClientLog("ActionType is not known. Do nothing")
		postResult(task)
	}
	ClientLog("task finished")
}

func DoTasks(t *[]TaskDetail) {
	for i, task := range *t {
		ClientLog("start do the task ", i, task.CmdType)
		task.ResultCode = 1
		DoTask(&task)
		if task.ResultCode == 0 {
			ClientLog("do task error", task.CmdType)
			break
		}

		ClientLog("end do the task ", i, task.CmdType)
	}
	// 任务执行完后删除map记录
	delete(MapTasks, (*t)[0].TaskId)
}

// 获取任务
func GetTask(taskUrl string) (t []TaskDetail, err error) {
	ClientLog("start request get task")
	client := http.DefaultClient
	req, err := http.NewRequest("GET", taskUrl, nil)
	if err != nil {
		ClientLog("new Request error:", err)
		return t, errors.New("new Request error" + err.Error())
	}
	//req.Host = "domain"
	resp, err := client.Do(req)

	if err != nil {
		ClientLog("get task error: ", err)
		return t, errors.New("get task error:" + err.Error())
	} else {
		defer resp.Body.Close()
		taskContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return t, errors.New("get task body error:" + err.Error())
		}
		ClientLog("get tasks content:", string(taskContent))
		json.Unmarshal(taskContent, &t)
		ClientLog("tasks :", t)
		if len(t) > 0 {
			task := t[0]
			task.ResultCode = 1
			task.ResultInfo = "get task success"
			task.Url = taskUrl
			ClientLog(task.Url, t[0].Url)
			postResult(&task)
			return t, nil
		}
	}
	return t, errors.New("not get tasks.")
}
