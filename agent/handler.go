/**
 totoroAgent project handler.go
 author:feng
 since:2020-01-08
**/
package agent

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var NetLisener net.Listener
var ChanNetClose = make(chan string)

// 读取请求数据，如果有加密就解密
func ReadBody(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	if req.Header.Get(SECURE_HEADER) == "TRUE" {
		descrypts, err := TripleDesDecrypt(body, AGENT_DES_KEY)
		if err != nil {
			return nil, err
		}
		return descrypts, nil
	}
	return body, nil
}

// 提交返回数据
func postResult(task *TaskDetail) {
	ClientLog(task.CmdType, "post data started:")
	j := MarshalTask(task)
	ClientLog("Post json string ", j)
	body := url.Values{"json": {j}}
	bodyContent := body.Encode()
	ClientLog("task type: ", task.CmdType, j)
	client := http.DefaultClient
	count := 0
	for {
		req, err := http.NewRequest("POST", task.Url, strings.NewReader(bodyContent))
		if err != nil {
			ClientLog("Http.NewRequest error: ", err)
		} else {
			//req.Host = "domain"
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			result := WaitForSever(*task, req, client)
			if result {
				break
			}
		}
		count++
		if count >= 2 {
			ClientLog("操作重试次数2，当前发布步骤回传失败: ", *task)
			break
		}
		time.Sleep(time.Second * 3)
	}

}

func WaitForSever(task TaskDetail, req *http.Request, client *http.Client) (result bool) {
	res, err := client.Do(req)
	if err != nil {
		ClientLog("client do request error:", err)
	} else {
		defer res.Body.Close()
		ClientLog("response status:", res.Status)
		if res.StatusCode == 204 {
			ClientLog("post success:", task.Url)
			return true
		} else {
			ClientLog("cannot connect to server:", task.Id, res.Status)
		}
	}
	return false
}

func DefaultHandler(w http.ResponseWriter, req *http.Request) {
	ClientLog("default handler...")
	fmt.Fprintln(w, AppConfig.Version)
	ClientLog("return version:" + AppConfig.Version)
}

func ListenAndServe() {
	http.Handle("/", http.HandlerFunc(DefaultHandler))
	http.Handle("/heart", http.HandlerFunc(HeartHandler))
	http.Handle("/version", http.HandlerFunc(VersionHandler))
	http.Handle("/exec", http.HandlerFunc(ExecHandler))
	http.Handle("/tasks", http.HandlerFunc(TaskHandler))
	ClientLog("start http server")
	go func() {
		go http.Serve(NetLisener, nil)
		select {
		case <-ChanNetClose:
			ClientLog("Process pid ", os.Getpid(),
				"stop listen on ", AppConfig.ListenPort)
		}
	}()
}

func StartListen(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		ClientLog("listen error:", err)
		os.Exit(6)
	}
	NetLisener = lis
	ClientLog("start listen on port ", addr)
	ListenAndServe()
}
