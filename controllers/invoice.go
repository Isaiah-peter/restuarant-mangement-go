package controllers

import (
	"context"
	"fmt"
	"golang-management-restaurant/database"
	"golang-management-restaurant/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvoiceVeiwFormat struct {
	InvoiceId        string
	Payment_method   string
	OrderId          string
	Payment_statue   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_detail     interface{}
}

var InvoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		result, err := InvoiceCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error will trying to get all invoice"})
		}

		var allInVoice []bson.M

		if err := result.All(ctx, &allInVoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while getting data "})
		}

		defer cancel()

		c.JSON(http.StatusOK, allInVoice)

	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		invoiceId := c.Param("invoice_id")
		var invoice models.Invoice

		err := InvoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while getting invoice error"})
		}
		var invoiceView InvoiceVeiwFormat

		allOrderItems, err := ItemByOrder(invoiceView.OrderId)
		invoiceView.OrderId = invoice.OrderId
		invoiceView.Payment_due_date = invoice.PaymentDueDate
		invoiceView.Payment_statue = *&invoice.PaymentStatus
		invoiceView.Payment_method = "null"
		if invoice.PaymentMethod != nil {
			invoiceView.Payment_method = *invoice.PaymentMethod
		}
		invoiceView.InvoiceId = invoice.InvoiceId
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_detail = allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)

		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderId}).Decode(&order)
		defer cancel()

		if err != nil {
			msg := fmt.Sprintf("message: order not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		status := "PENDING"
		if invoice.PaymentStatus != nil {
			invoice.PaymentStatus = &status
		}
		invoice.CreatedAt = timetouse
		invoice.UpdatedAt = timetouse
		invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()

		invoice.InvoiceId = invoice.ID.Hex()

		validateErr := validate.Struct(invoice)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while yoy input your invoice"})
			return
		}

		result, err := InvoiceCollection.InsertOne(ctx, invoice)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while creating invoice "})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)

		var invoice models.Invoice
		var invoiceId = c.Param("invoice_id")
		var filter = bson.M{"invoice_id": invoiceId}
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"update_at", invoice.UpdatedAt})

		if invoice.PaymentMethod != nil {
			updateObj = append(updateObj, bson.E{"payment_method", invoice.PaymentMethod})
		}
		if invoice.PaymentStatus != nil {
			updateObj = append(updateObj, bson.E{"payment_status", invoice.PaymentStatus})
		}

		status := "PENDING"

		if invoice.PaymentStatus != nil {
			invoice.PaymentStatus = &status
			updateObj = append(updateObj, bson.E{"payment_status", invoice.PaymentStatus})
		}

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := InvoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating invioce"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}
