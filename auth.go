package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var secretkey string = "eventAuthJwtKey-28042021"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func IsAuthorized(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.GetHeader("Token") == "" {
			var err Error
			context.JSON(http.StatusUnauthorized, SetError(err, "No Token Found"))
			return
		}
		var mySigningKey = []byte(secretkey)

		token, err := jwt.Parse(context.GetHeader("Token"), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token.")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err Error
			context.JSON(http.StatusUnauthorized, SetError(err, "Your Token has been expired."))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				context.Request.Header.Set("Role", "admin")
				handler(context)
				return

			} else if claims["role"] == "user" {
				context.Request.Header.Set("Role", "user")
				return
			}
		}
		var reserr Error
		context.JSON(http.StatusUnauthorized, SetError(reserr, "Not Authorized."))
		return
	}
}

func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}
