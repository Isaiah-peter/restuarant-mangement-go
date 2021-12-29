package routes

import (
	"golang-management-restaurant/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	route.GET("/users", controllers.GetUsers())
	route.GET("/user/:user_id", controllers.GetUser())
	route.POST("/register", controllers.SingUp())
	route.POST("/login", controllers.Login())
}
