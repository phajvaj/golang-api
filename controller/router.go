package controller

import (
	"github.com/gin-gonic/gin"
)

//router config
/*func routers() *gin.Engine {
	rt := gin.Default()
	rt.GET("/", GetIndex)
	rt.POST("/", model.TokenAuthMiddleware(), AddData)
	rt.PUT("/:userId", EditData)
	rt.POST("/login", LoginUser)
	//router.GET("/name/:msg", getName)

	return rt
}
*/
type routes struct {
	router *gin.Engine
}

func SetRoutes() routes {
	r := routes{
		router: gin.Default(),
	}
	r.router.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.router.Static("/web", "./public")
	r.router.StaticFile("/favicon.ico", "./images/favicon.ico")

	v1 := r.router.Group("/v1")

	r.addIndex(v1)
	r.addLogin(v1)

	return r
}

func (r routes) Path() *gin.Engine {
	return r.router
}

func (r routes) Run(addr ...string) error {
	return r.router.Run()
}
