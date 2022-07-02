package gee

import (
	"fmt"
	"net/http"
)

//定义了类型HandlerFunc，这是提供给框架用户的，用来定义路由映射的处理方法，
/*
我们在Engine中，添加了一张路由映射表router，key 由请求方法和静态路由地址构成，例如GET-/、GET-/hello、POST-/hello，这样针对相同的路由，
如果请求方法不同,可以映射不同的处理方法(Handler)，value 是用户映射的处理方法
*/
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

/*
当用户调用(*Engine).GET()方法时，
会将路由和处理方法注册到映射表 router 中，(*Engine).Run()方法，是 ListenAndServe 的包装
*/
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

/*
Engine实现的 ServeHTTP 方法的作用就是，解析请求的路径，查找路由映射表，如果查到
，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 not found:%s \n", req.URL)
	}
}

/*
到这里为止，整个Gee框架的原型已经出来了。实现了路由映射表，
提供了用户注册静态路由的方法，包装了启动服务的函数。当然，到目前为止，我们还没有实现比net/http标准库更强大的能力，
不用担心，很快就可以将动态路由、中间件等功能添加上去了。
*/
