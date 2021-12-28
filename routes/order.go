package routes

import (
	"golang-management-restaurant/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(route *gin.Engine) {
	route.GET("/orders", controllers.GetOrders())
	route.GET("/order/:order_id", controllers.GetOrder())
	route.POST("/orders", controllers.CreateOrder())
	route.POST("/order/:order_id", controllers.UpdateOrder())
}
