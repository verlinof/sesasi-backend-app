// Run projects
// compiledaemon --command="./sesasi-backend-app"
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/verlinof/sesasi-backend-app/controllers/user_controller"
	"github.com/verlinof/sesasi-backend-app/initializers"
)

func init() {
	//Load .env file
	initializers.LoadEnvVariables()
	//Connect database
	initializers.ConnectToDB()
	//Migrate
	initializers.Migrate()
}

func main() {
	route := gin.Default()

	route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	route.POST("/signup", user_controller.SignUp)
	route.POST("/login", user_controller.Login)

	route.Run()
}
