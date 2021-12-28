package routes

import (
	"golang-management-restaurant/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(route *gin.Engine) {
	route.GET("/invoices", controllers.CreateInvoice())
	route.GET("/invoice/:invoice_id", controllers.GetInvoice())
	route.POST("/invoices", controllers.CreateInvoice())
	route.POST("/invoice/:invoice_id", controllers.UpdateInvoice())
}
