package main

import (
	"flag"
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

func createGroup() *geecache.Group {
	return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowdb] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *geecache.Group) {
	peers := geecache.NewHTTPPool(addr)

	peers.Set(addrs...)
	gee.Registerpeers(peers)
	log.Println("geecache is runnning at ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))

	log.Println("fonted server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {

	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "GeeCache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()

	if api {
		go startAPIServer(apiAddr, gee)
	}

	startCacheServer(addrMap[port], []string(addrs), gee)
}

/*

main 函数的代码比较多，但是逻辑是非常简单的。

startCacheServer() 用来启动缓存服务器：创建 HTTPPool，添加节点信息，注册到 gee 中，启动 HTTP 服务（共3个端口，8001/8002/8003），用户不感知。
startAPIServer() 用来启动一个 API 服务（端口 9999），与用户进行交互，用户感知。
main() 函数需要命令行传入 port 和 api 2 个参数，用来在指定端口启动 HTTP 服务。



*/

/*

为了方便，我们将启动的命令封装为一个 shell 脚本：


#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 &
./server -port=8002 &
./server -port=8003 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait

trap 命令用于在 shell 脚本退出时，删掉临时文件，结束子进程。


$ ./run.sh
2020/02/16 21:17:43 geecache is running at http://localhost:8001
2020/02/16 21:17:43 geecache is running at http://localhost:8002
2020/02/16 21:17:43 geecache is running at http://localhost:8003
2020/02/16 21:17:43 fontend server is running at http://localhost:9999
>>> start test
2020/02/16 21:17:45 [Server http://localhost:8003] Pick peer http://localhost:8001
2020/02/16 21:17:45 [Server http://localhost:8003] Pick peer http://localhost:8001
2020/02/16 21:17:45 [Server http://localhost:8003] Pick peer http://localhost:8001
...
630630630
此时，我们可以打开一个新的 shell，进行测试：


$ curl "http://localhost:9999/api?key=Tom"
630
$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
测试的时候，我们并发了 3 个请求 ?key=Tom，从日志中可以看到，三次均选择了节点 8001，这是一致性哈希算法的功劳。但是有一个问题在于，同时向 8001 发起了 3 次请求。试想，假如有 10 万个在并发请求该数据呢？那就会向 8001 同时发起 10 万次请求，如果 8001 又同时向数据库发起 10 万次查询请求，很容易导致缓存被击穿。

三次请求的结果是一致的，对于相同的 key，能不能只向 8001 发起一次请求？这个问题下一次解决。

*/
