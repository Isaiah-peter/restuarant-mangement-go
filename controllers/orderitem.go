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
	"go.mongodb.org/mongo-driver/mongo/options"
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
		orderId := c.Param("order_id")
		orderItems, err := ItemByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occur while listing orderitem by order Id"})
			return
		}

		c.JSON(http.StatusOK, orderItems)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		orderItemId := c.Param("order_item_id")
		var orderItem models.OrderItem

		err := orderItemCollection.FindOne(ctx, bson.M{"orderItem_id": orderItemId}).Decode(&orderItem)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching orderitem"})
			return
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemTobeInserted := []interface{}{}
		order.TableId = orderItemPack.TableId
		orderId := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.OrderItem {
			orderItem.OrderId = orderId

			validarteErr := validate.Struct(orderItem)
			if validarteErr != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": validarteErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.OrderItemId = orderItem.ID.Hex()
			var num = toFixed(*&orderItem.UnitPrice, 2)
			orderItem.UnitPrice = num
			orderItemTobeInserted = append(orderItemTobeInserted, orderItem)
		}

		insertOrderItems, err := orderItemCollection.InsertOne(ctx, orderItemTobeInserted)

		if err != nil {
			log.Fatal(err)
		}
		defer cancel()

		c.JSON(http.StatusOK, insertOrderItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var orderItems models.OrderItem

		orderItemId := c.Param("order_irem_id")

		filter := bson.M{"order_item_id": orderItemId}

		if err := c.BindJSON(&orderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		var updateObj primitive.D

		if orderItems.UnitPrice != 0 {
			updateObj = append(updateObj, bson.E{"unit_price", &orderItems.UnitPrice})
		}

		if orderItems.FoodId != nil {
			updateObj = append(updateObj, bson.E{"food_id", orderItems.FoodId})
		}

		if orderItems.Quantity != "" {
			updateObj = append(updateObj, bson.E{"quantity", &orderItems.Quantity})
		}

		orderItems.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItems.UpdatedAt})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}

func ItemByOrder(id string) (orderitems []primitive.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"form", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unWindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unWindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrys", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unWindTableStage := bson.D{{"$unwind", bson.D{{"path", "table"}, {"preserveNullAndEmptyArrays", true}}}}

	projectStage := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"amount", "$food.price"},
			{"total_count", 1},
			{"food_name", "$food.name"},
			{"food_image", "$food.food_image"},
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
			{"price", "$food.price"},
			{"quantity", 1},
		}}}

	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}, {"patment_due", bson.D{{"$sum", "$amount"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"order_items", bson.D{{"$push", "$$ROOT"}}}}}}}}

	projectStage2 := bson.D{
		{"id", 0},
		{"payment_due", 1},
		{"total_count", 1},
		{"table_number", "$_.id.table_number"},
		{"order_items", 1},
	}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unWindStage,
		lookupOrderStage,
		unWindOrderStage,
		lookupTableStage,
		unWindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})

	if err != nil {
		panic(err)
	}

	if err = result.All(ctx, &orderitems); err != nil {
		panic(err)
	}

	defer cancel()

	return orderitems, err

}
