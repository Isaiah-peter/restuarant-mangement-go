package routes

import (
	"github.com/gin-gonic/gin"
	"golang-restaurant-management/controllers"
)

func UserRoutes(route *gin.Engine) {
	route.GET("/users", controllers.GetUsers())
	route.GET("/user/:user_id", controllers.GetUser())
	route.POST("/register", controllers.SingUp())
	route.POST("/login", controllers.Login())
}