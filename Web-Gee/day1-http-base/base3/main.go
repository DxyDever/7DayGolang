package base3

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {

	r := gee.New()

	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.PATH = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":9999")
}

/*

看到这里，如果你使用过gin框架的话，肯定会觉得无比的亲切。
gee框架的设计以及API均参考了gin。使用New()创建 gee 的实例，
使用 GET()方法添加路由，最后使用Run()启动Web服务。这里的路由，只是静态路由，不支持/hello/:name这样的动态路由，
动态路由我们将在下一次实现

*/
