package controller

import (
	"crypto/md5"
	"fmt"
	"nan_api_main/config"
	"nan_api_main/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r routes) addLogin(rg *gin.RouterGroup) {
	rt := rg.Group("/login")
	rt.POST("/", LoginUser)
}
func LoginUser(c *gin.Context) {
	var u model.LoginUser
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "error": "Invalid json provided"})
		return
	}
	db, err := config.SetConnection()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "error": "not connect dnms"})
		return
	}
	userModel := model.UsersModel{
		Db: db,
	}
	pwd := []byte(u.Password)

	rsUser, err := userModel.FindLogin(u.Username, fmt.Sprintf("%x", md5.Sum(pwd)))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "error": "youer check username and password or disable account"})
		return
	}

	//compare the user from the request, with the one we defined:
	if len(rsUser) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "Please provide valid login details"})
		return
	}

	token, err := model.CreateToken(int(rsUser[0].UserId), rsUser[0].Username, rsUser[0].FirstName+" "+rsUser[0].LastName, rsUser[0].UserType)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"ok": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "token": token, "rows": rsUser})
}
