/**
 totoroAgent project aignalhandler_windows.go
 author:feng
 since:2020-01-08
**/
package agent

import (
	"fmt"
	"net/http"
)

func VersionHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, AppConfig.Version)
}
