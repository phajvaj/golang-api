package controller

import (
	"nan_api_main/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r routes) addIndex(rg *gin.RouterGroup) {
	rt := rg.Group("/index")

	rt.GET("/", GetIndex)
	rt.POST("/", model.TokenAuthMiddleware(), AddData)
	rt.PUT("/:userId", EditData)
}

func GetIndex(c *gin.Context) {
	name := c.Query("name")
	c.JSON(http.StatusOK, gin.H{"message": "Hello World My name is: " + name})
}

func AddData(c *gin.Context) {
	var u model.UserJson
	c.BindJSON(&u)
	c.JSON(http.StatusOK, gin.H{"ok": true, "user": u})
}

func EditData(c *gin.Context) {
	userId := c.Param("userId")
	var u model.UserJson
	c.BindJSON(&u)
	c.JSON(http.StatusOK, gin.H{"ok": true, "Id": userId, "rows": u})
}
