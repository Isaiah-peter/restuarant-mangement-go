package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
