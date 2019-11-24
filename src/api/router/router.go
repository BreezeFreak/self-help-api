package router

import (
	"api/app"
	"api/config"
	"api/utils/middleware"
	"api/utils/response"
	"github.com/gin-gonic/gin"
)

var routerMap = map[string]func(c *gin.Context, u *app.UrlDivide){
	"GET":    app.GET,
	"POST":   app.POST,
	"PUT":    app.PUT,
	"PATCH":  app.PATCH,
	"DELETE": app.DELETE,
}

func Init() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Options)
	r.Use(middleware.PrintErrorLog)

	gin.SetMode(config.Server.RunMode)

	// TODO permission by ip? or share everything. record {ip:x, col:x} anyway. use middleware for login and recording

	r.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			app.Index(c)
			return
		}

		var u app.UrlDivide
		if e := u.ParseUrl(c.Request.URL.Path); e != nil {
			ginResponse.BadRequest(c, e)
			return
		}

		if f, ok := routerMap[c.Request.Method]; ok {
			f(c, &u)
			return
		}
		ginResponse.MethodNotAllowed(c)
		return
	})
	return r
}
