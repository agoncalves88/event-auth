package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	auth "github.com/agoncalves88/event-auth"
	configuration "github.com/agoncalves88/event-auth/internal/configuration"
	database "github.com/agoncalves88/event-auth/internal/database"
	"github.com/gin-gonic/gin"
)

var config configuration.Configuration

func main() {
	env := flag.String("env", "", "")
	flag.Parse()
	r := gin.Default()
	//r.Use(gzip.Gzip(gzip.BestSpeed))
	config = configuration.GetConfig(*env)
	database.InitialMigration(config.ConnectionString)
	oauth := r.Group("/OAuth")
	{
		oauth.POST("/SignIn", SignIn)
		oauth.POST("/SignUp", SignUp)
	}
	r.Run(":" + config.Port)
}

func TestAuth(c *gin.Context) {
	c.Status(http.StatusAccepted)
}

func SignUp(c *gin.Context) {
	connection := database.GetDatabase(config.ConnectionString)
	defer database.CloseDatabase(connection)
	c.Request.Header.Set("Content-Type", "application/json")

	var user auth.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		var err auth.Error
		c.JSON(http.StatusUnprocessableEntity, auth.SetError(err, "Error in reading payload."))
		return
	}

	var dbuser auth.User
	connection.Where("email = ?", user.Email).First(&dbuser)

	if dbuser.Email != "" {
		var err auth.Error
		c.JSON(http.StatusBadRequest, auth.SetError(err, "Email already in use"))
		return
	}

	user.Password, err = auth.GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("Error in password hashing.")
		var err auth.Error
		c.JSON(http.StatusBadRequest, auth.SetError(err, "Error in password hashing."))
		return
	}
	connection.Create(&user)
	c.Status(http.StatusCreated)
}

func SignIn(c *gin.Context) {
	connection := database.GetDatabase(config.ConnectionString)
	c.Request.Header.Set("Content-Type", "application/json")
	defer database.CloseDatabase(connection)

	var authDetails auth.Authentication

	err := json.NewDecoder(c.Request.Body).Decode(&authDetails)
	if err != nil {
		var err auth.Error
		c.JSON(http.StatusUnprocessableEntity, auth.SetError(err, "Error in reading payload."))
		return
	}

	var authUser auth.User
	connection.Where("email = 	?", authDetails.Email).First(&authUser)

	if authUser.Email == "" {
		var err auth.Error

		c.JSON(http.StatusUnauthorized, auth.SetError(err, "Username or Password is incorrect"))
		return
	}

	check := auth.CheckPasswordHash(authDetails.Password, authUser.Password)

	if !check {
		var err auth.Error
		c.JSON(http.StatusUnauthorized, auth.SetError(err, "Username or Password is incorrect"))
		return
	}

	validToken, err := auth.GenerateJWT(authUser.Email, authUser.Role)
	if err != nil {
		var err auth.Error
		c.JSON(http.StatusBadRequest, auth.SetError(err, "Failed to generate token"))
		return
	}

	var token auth.Token
	token.Email = authUser.Email
	token.Role = authUser.Role
	token.TokenString = validToken
	c.JSON(http.StatusOK, token)
}
