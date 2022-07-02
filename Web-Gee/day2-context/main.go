package day2_context

import (
	"gee"
	"net/http"
)

/*
- 将路由独立出来，方便之后进行增强
- 设计上下文，封装Request和Response，提供对json，html等返回类型的
支持
*/

func main() {

	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		//期望 /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"))
	})

	r.POST("/login", func() {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.Run(":9999")
}

/*

Handler的参数变成成了gee.Context，提供了查询Query/PostForm参数的功能。
gee.Context封装了HTML/String/JSON函数，能够快速构造HTTP响应。

*/
