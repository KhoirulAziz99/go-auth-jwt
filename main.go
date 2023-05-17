package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KhoirulAziz99/go-auth-jwt/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var jwtKey = "SECRET_CODE"

func main() {
	//gin router
	r := gin.Default()
	//setup router
	r.POST("/auth/login", loginHandler)

	userRouter := r.Group("api/v1/users")

	//middleware
	userRouter.Use(auth.AuthMiddleware())

	//setup get user profile
	userRouter.GET("/:id/profile", profileHandler)

	//start server
	r.Run(":8080")
}

func loginHandler(c *gin.Context) {
	var user User

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(user.Username)
	fmt.Println(user.Password)
	//logic authentication(compare username and password)
	if user.Username == "enigma" && user.Password == "12345" {
		//bikin code untuk generate token
		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims) // ini map

		claims["username"] = user.Username
		claims["exp"] = time.Now().Add(time.Minute * 1).Unix() //token akan expired dalam 1 menit

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func profileHandler(c *gin.Context) {
	// Ambil username dari JWT token
	claims := c.MustGet("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	// seharusnya return user  dari database , tapi di contoh ini kita pakai code block
	c.JSON(http.StatusOK, gin.H{"username": username})
}
