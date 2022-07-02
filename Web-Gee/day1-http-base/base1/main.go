package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
简单介绍net/http库以及http.Handler接口
搭建Gee框架的雏形
*/

/*
标准库启动web服务，Go语言内置了net/http库，封装了HTTP网络编程的基础的接口，我们实现的Gee Web 框架便是基于net/http的。
我们接下来通过一个例子，简单介绍下这个库的使用
*/

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)

	/*
		log.Fatal() 函数在控制台屏幕上打印带有时间戳的指定消息。
		log.Fatal() 类似于 log.Print() 函数，后跟调用 os.Exit(1) 函数
	*/
	log.Fatal(http.ListenAndServe(":9999", nil))
}

/*
用法:

func Fprintf(w io.Writer, format string, a ...any)(n int, err error)
Fprintf 根据格式说明符格式化并写入 w。它返回写入的字节数和遇到的任何写入错误。

*/
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.PATH = %q\n", req.URL.Path)
}

/*
字符串和字节片

%s	输出字符串表示（string类型或[]byte)	Printf("%s", []byte(“Go语言”))	Go语言
%q	双引号围绕的字符串，由Go语法安全地转义	Printf("%q", “Go语言”)	“Go语言”


*/
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header [%q] = %q\n", k, v)
	}
}

/*
我们设置了2个路由，/和/hello，分别绑定 indexHandler 和 helloHandler ， 根据不同的HTTP请求
会调用不同的处理函数。访问/，响应是URL.Path = /，而/hello的响应则是请求头(header)中的键值对信息。
使用curl工具测试

$ curl http://localhost:9999/
URL.Path = "/"
$ curl http://localhost:9999/hello
Header["Accept"] = ...
Header["User-Agent"] = ["curl/7.54.0"]

main 函数的最后一行，是用来启动 Web 服务的，第一个参数是地址，:9999表示在 9999 端口监听。而第二个参数则代表处理所有的HTTP请求的实例，nil
代表使用标准库中的实例处理。第二个参数，则是我们基于net/http标准库实现Web框架的入口。
*/
