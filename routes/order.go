package routes

import (
	"golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(route *gin.Engine) {
	route.GET("/orders", controllers.GetOrders())
	route.GET("/order/:order_id", controllers.GetOrder())
	route.POST("/orders", controllers.CreateOrder())
	route.POST("/order/:order_id", controllers.UpdateOrder())
}
