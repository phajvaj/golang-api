package controller

import (
	"fmt"
	"log"
	"nan_api_main/model"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func pdfPath() string {
	return "./" + viper.GetString("app.public") + "/pdf/"
}

func (r routes) addIndex(rg *gin.RouterGroup) {
	rt := rg.Group("/index")

	rt.GET("/", GetIndex)
	var AdminType = "ADMIN"
	var DoctorType = "DOCTOR"
	rt.POST("/", model.TokenAuthMiddleware(&AdminType), AddData)
	rt.POST("/upload", model.TokenAuthMiddleware(nil), SetUpload)
	rt.PUT("/:userId", EditData)
	rt.GET("/pdf/:id", model.TokenAuthMiddleware(&DoctorType), GetPdf)
	rt.GET("/pdfno/:id", GetPdf)
	rt.GET("/png", GetPng)
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
	uid := c.PostForm("id")

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	_path := pdfPath()
	if _, err := os.Stat(_path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		err = os.Mkdir(_path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	//filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, _path+uid+".pdf"); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "pdf": fmt.Sprintf("File %s uploaded successfully with fields name=%s and email=%s.", file.Filename, name, email)})
}

func GetPdf(c *gin.Context) {
	id := c.Param("id")
	_path := pdfPath()
	afile := _path + id + ".pdf"
	fmt.Println(afile)
	if _, err := os.Stat(afile); os.IsNotExist(err) {
		c.JSON(http.StatusNoContent, gin.H{"error": err})
		return
	}
	c.File(afile)
	//c.JSON(http.StatusOK, gin.H{"PDF": afile})
}

func GetPng(c *gin.Context) {
	response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	defer reader.Close()
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="gopher.png"`,
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
