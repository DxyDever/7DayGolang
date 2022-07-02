package main

import (
	"log"
	"net/http"
)

type server string

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL.Path)

	w.Write([]byte("hello world"))
}

func main() {
	var s server
	http.ListenAndServe("localhost:9999", &s)
}

/*
创建任意类型 server，并实现 ServeHTTP 方法。
调用 http.ListenAndServe 在 9999 端口启动 http 服务，处理请求的对象为 s server。

接下来我们执行 go run . 启动服务，借助 curl 来测试效果：

$ curl http://localhost:9999
Hello World!
$ curl http://localhost:9999/abc
Hello World!
Go 程序日志输出

2020/02/11 22:56:32 /
2020/02/11 22:56:34 /abc


http.ListenAndServe 接收 2 个参数，第一个参数是服务启动的地址，第二个参数是 Handler，任何实现了 ServeHTTP 方法的对象都可以作为 HTTP 的 Handler。


在标准库中，http.Handler 接口的定义如下：

package http

type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}

*/
