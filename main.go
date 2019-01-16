package main

import (
	"github.com/gin-gonic/gin"
	"app/model"
	"app"
)


func main() {
	r := gin.Default()

	model.InitDatabase()

	api := r.Group("/api/v1")
	{
		api.POST("/order-restock", app.CreateRestockOrder)
		api.POST("/receive-restock", app.ReceiveRestock)
		api.POST("/purchase", app.CreatePurchaseOrder)
		api.GET("/stock-info", app.GetStockInfo)
		api.GET("/restock-order-log", app.GetRestockOrderLog)
		api.GET("/purchase-order-log", app.GetPurchaseOrderLog)
		api.GET("/item-inventory/:format", app.GetItemInventoryReport)
		api.DELETE("/sales/:format", app.GetSalesReport)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
