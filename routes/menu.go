package routes

import (
	"golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(route *gin.Engine) {
	route.GET("/menus", controllers.GetMenus())
	route.GET("/menu/:menu_id", controllers.GetMenu())
	route.POST("/menus", controllers.CreateMenu())
	route.POST("/menu/:menu_id", controllers.UpdateMenu())
}
