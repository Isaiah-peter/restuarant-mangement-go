package main

import (
	"github.com/gin-gonic/gin"
	"golang-restaurant-management/database"
	"golang-restaurant-management/routes"
	"os"

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


}
