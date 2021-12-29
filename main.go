package main

import (
	"golang-management-restaurant/database"
	"golang-management-restaurant/middlewares"
	routes "golang-management-restaurant/routes"
	"os"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = ":4000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middlewares.Authentication())

	routes.FoodRoutes(router)
	routes.InvoiceRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderItemRoutes(router)
	routes.OrderRoutes(router)
	routes.TableRoutes(router)

	router.Run(":" + PORT)
}
