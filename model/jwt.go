package model

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func CreateToken(Id int, Username, Fullname, UserType string) (string, error) {
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userId"] = Id
	atClaims["username"] = Username
	atClaims["fullname"] = Fullname
	atClaims["userType"] = UserType
	atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(viper.GetString("app.serctkey")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func TokenAuthMiddleware(uType *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request, uType)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

func TokenValid(r *http.Request, uType *string) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	//fmt.Println(*uType)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//fmt.Println(claims["userType"], claims["username"])
		if uType != nil && *uType != claims["userType"] {
			err := errors.New("Auth middleware is not validate!")
			return err
		}
	} else {
		err := errors.New("Token is not validate!")
		return err
	}

	/*if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}*/

	return nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("app.serctkey")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
