package routes

import (
	"golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(route *gin.Engine) {
	route.GET("/tables", controllers.GetTables())
	route.GET("/table/:table_id", controllers.GetTable())
	route.POST("/tables", controllers.CreateTable())
	route.POST("/table/:table_id", controllers.UpdateTable())
}
