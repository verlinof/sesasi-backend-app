package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/verlinof/sesasi-backend-app/initializers"
	"github.com/verlinof/sesasi-backend-app/models"
)

func RequireAuth(c *gin.Context) {
	//Get token from header
	tokenString := c.Request.Header.Get("Authorization")
	tokenParts := strings.Split(tokenString, " ")
	if len(tokenParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString = tokenParts[1]

	if tokenString == "" {
		//When Using Cookie
		var errCookie error
		tokenString, errCookie = c.Cookie("Authorization")
		if errCookie != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": errCookie.Error()})
			return
		}
	}

	//Decode/Validate it
	// Parse the token
	token, errToken := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if errToken != nil {
		fmt.Println(token)
		c.AbortWithStatus(http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{"message": errToken.Error()})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {		
		//Find the user
		user := new(models.User)
		initializers.DB.First(&user, claims["id"]) //Claims digunakan untuk ngambil data dari JWT

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//Attach to req
		c.Set("currentUser", user)

		//Continue
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}