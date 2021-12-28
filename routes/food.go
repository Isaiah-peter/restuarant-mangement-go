package routes

import (
	"golang-management-restaurant/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(route *gin.Engine) {
	route.GET("/foods", controllers.GetFoods())
	route.GET("/food/:food_id", controllers.GetFood())
	route.POST("/foods", controllers.CreateFood())
	route.POST("/food/:food_id", controllers.UpdateFood())
}
