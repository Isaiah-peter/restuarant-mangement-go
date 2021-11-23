package routes

import (
	"github.com/gin-gonic/gin"
	"golang-restaurant-management/controllers"
)

func FoodRoutes(route *gin.Engine) {
	route.GET("/foods", controllers.GetFoods())
	route.GET("/food/:food_id", controllers.GetFood())
	route.POST("/foods", controllers.CreateFood())
	route.POST("/food/:food_id", controllers.UpdateFood())
}