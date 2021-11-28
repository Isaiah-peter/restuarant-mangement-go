package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderItemPack struct {
	TableId   *string
	OrderItem []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderItemCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while getting orderitem data"})
			return
		}

		var allOrderItem []bson.M

		if err := result.All(ctx, &allOrderItem); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allOrderItem)
	}
}

func GetOrderItemByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ItemByOrder(id string) (orderitems []primitive.M, err error) {
	return nil, err
}
