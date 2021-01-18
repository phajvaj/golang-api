package controller

import (
	"fmt"
	"nan_api_main/model"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (r routes) addIndex(rg *gin.RouterGroup) {
	rt := rg.Group("/index")

	rt.GET("/", GetIndex)
	rt.POST("/", model.TokenAuthMiddleware(), AddData)
	rt.POST("/upload", SetUpload)
	rt.PUT("/:userId", EditData)
	rt.GET("/pdf/:id", GetPdf)
}

func GetIndex(c *gin.Context) {
	name := c.Query("name")
	c.JSON(http.StatusOK, gin.H{"message": "Hello World My name is: " + name})
}

func AddData(c *gin.Context) {
	var u model.UsersLoginJson
	c.BindJSON(&u)
	c.JSON(http.StatusOK, gin.H{"ok": true, "user": u})
}

func EditData(c *gin.Context) {
	userId := c.Param("userId")
	var u model.UsersLoginJson
	c.BindJSON(&u)
	c.JSON(http.StatusOK, gin.H{"ok": true, "Id": userId, "rows": u})
}

func SetUpload(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, "./public/pdf/"+filename); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "pdf": fmt.Sprintf("File %s uploaded successfully with fields name=%s and email=%s.", file.Filename, name, email)})
}

func GetPdf(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"PDF": id + ".pdf"})
}
