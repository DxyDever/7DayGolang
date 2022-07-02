package main

import (
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {

	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key", key)

			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not Exist", key)
		}))

	addr := "localhost:9999"

	peers := geecache.NewHTTPPool(addr)

	log.Println("geecache is running at", addr)

	log.Fatal(http.ListenAndServe(addr, peers))
}

/*

同样地，我们使用 map 模拟了数据源 db。
创建一个名为 scores 的 Group，若缓存为空，回调函数会从 db 中获取数据并返回。
使用 http.ListenAndServe 在 9999 端口启动了 HTTP 服务。

*/

/*

需要注意的点：
main.go 和 geecache/ 在同级目录，但 go modules 不再支持 import <相对路径>，相对路径需要在 go.mod 中声明：
require geecache v0.0.0
replace geecache => ./geecache

*/

/*

接下来，运行 main 函数，使用 curl 做一些简单测试：

$ curl http://localhost:9999/_geecache/scores/Tom
630
$ curl http://localhost:9999/_geecache/scores/kkk
kkk not exist


GeeCache 的日志输出如下：

2020/02/11 23:28:39 geecache is running at localhost:9999
2020/02/11 23:29:08 [Server localhost:9999] GET /_geecache/scores/Tom
2020/02/11 23:29:08 [SlowDB] search key Tom
2020/02/11 23:29:16 [Server localhost:9999] GET /_geecache/scores/kkk
2020/02/11 23:29:16 [SlowDB] search key kkk
节点间的相互通信不仅需要 HTTP 服务端，还需要 HTTP 客户端，这就是我们下一步需要做的事情。

*/
