package main

import (
	"github.com/gin-gonic/gin"
	itemImpl "ijah-store/domain/item/impl"
	purchaseImpl "ijah-store/domain/purchase/impl"
	restockImpl "ijah-store/domain/restock/impl"
	"ijah-store/handler"
	"ijah-store/pkg"
)

func main() {
	r := gin.Default()

	pkg.InitDatabase()

	// init service & repo
	itemRepo := itemImpl.NewRepo(pkg.Database)
	itemSvc := itemImpl.NewService(itemRepo)

	purchaseRepo := purchaseImpl.NewRepo(pkg.Database)
	purchaeSvc := purchaseImpl.NewService(purchaseRepo)

	restockRepo := restockImpl.NewRepo(pkg.Database)
	restockSvc := restockImpl.NewService(restockRepo)

	// init handler
	apiHandler := handler.NewApiHandler(itemSvc, purchaeSvc, restockSvc)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/order-restock", apiHandler.CreateRestockOrder)
		v1.POST("/receive-restock",  apiHandler.ReceiveRestock)
		v1.POST("/purchase",  apiHandler.CreatePurchaseOrder)
		v1.GET("/stock-info",  apiHandler.GetStockInfo)
		v1.GET("/restock-order-log",  apiHandler.GetRestockOrderLog)
		v1.GET("/purchase-order-log",  apiHandler.GetPurchaseOrderLog)
		v1.GET("/item-inventory",  apiHandler.GetItemInventoryReport)
		v1.GET("/sales",  apiHandler.GetSalesReport)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
