package user_controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/verlinof/sesasi-backend-app/initializers"
	"github.com/verlinof/sesasi-backend-app/models"
	"github.com/verlinof/sesasi-backend-app/requests"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	request := new(requests.UserRequest)
	//Get email/pass off req body
	errRequest := c.ShouldBind(&request)

	if errRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Bad Request",
		})
		return
	}

	//Hash password
	hash, errHash := bcrypt.GenerateFromPassword([]byte(request.Password), 10)

	if errHash != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to hash password",
		})
		return
	}

	//Create user
	user := &models.User{
		Email:    request.Email,
		Password: string(hash),
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error:" + result.Error.Error(),
		})
		return
	}

	//Respond
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created",
		"data": result,
	})
}

func Login(c *gin.Context) {
	//Get data from req body	
	request := new(requests.UserRequest)
	errRequest := c.ShouldBind(&request)

	if errRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Bad Request",
		})
		return
	}

	//Find requested User
	user := new(models.User)
	initializers.DB.First(&user, "email = ?", request.Email)

	//Compare password
	errCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if errCompare != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid credentials",
		})
		return
	}

	//Generate JWT
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	secretKey := os.Getenv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error:" + err.Error(),
		})
		return
	}

	//Respond
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token": tokenString,
	})
}