/**
 totoroAgent project hearthandler.go
 author:feng
 since:2020-01-08
**/
package totoroAgent

import (
	"fmt"
	"net/http"
)

func HeartHandler(w http.ResponseWriter, req *http.Request) {
	ClientLog("start heart log request")
	fmt.Fprintln(w, req.Method)
	ClientLog("end heart log request")
}
