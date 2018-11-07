package main

import (
	"net/http"
	console "shell-exec/http/console"
	"shell-exec/router"
	"time"
)

func main() {

	// 启动一个nsq consumer
	console.ConsumerLogUpload()

	// 启动http服务
	routersInit := router.InitRouter()

	s := &http.Server{
		Addr:           ":8333",
		Handler:        routersInit,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   600 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}
