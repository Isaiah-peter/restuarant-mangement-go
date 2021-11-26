package controllers

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error why listing order items"})
		}

		var allOrders []bson.M

		if err := result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)

	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occur while get order data"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var order models.Order
		var table models.Table

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		validateErr := validate.Struct(order)
		if validateErr != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": validateErr.Error()})
			return
		}

		err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("table is not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()

		order.OrderId = order.ID.Hex()

		result, insertErr := orderCollection.InsertOne(ctx, order)
		if insertErr != nil {
			msg := fmt.Sprintf("chould not create order")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var order models.Order
		var table models.Table

		var orderId = c.Param("order_id")
		var filter = bson.M{"order_id": orderId}

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if order.TableId != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("mesage: table not found")
				c.JSON(http.StatusInternalServerError, gin.H{"msg": msg})
				return
			}
			updateObj = append(updateObj, bson.E{"table_id", order.TableId})
		}

		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"update_at", order.UpdatedAt})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating order"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}
