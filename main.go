package main

import (
	"github.com/gin-gonic/gin"
	"ijah-store/handler"
	"ijah-store/pkg"
)

func main() {
	r := gin.Default()

	pkg.InitDatabase()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/order-restock", handler.CreateRestockOrder)
		v1.POST("/receive-restock", handler.ReceiveRestock)
		v1.POST("/purchase", handler.CreatePurchaseOrder)
		v1.GET("/stock-info", handler.GetStockInfo)
		v1.GET("/restock-order-log", handler.GetRestockOrderLog)
		v1.GET("/purchase-order-log", handler.GetPurchaseOrderLog)
		v1.GET("/item-inventory", handler.GetItemInventoryReport)
		v1.GET("/sales", handler.GetSalesReport)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
