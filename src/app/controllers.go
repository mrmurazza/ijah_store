package app

import (
	"github.com/gin-gonic/gin"
	"app/request"
)

func CreateRestockOrder(c *gin.Context) {
	var request request.RestockOrderRequest
	c.BindJSON(&request)
}
func ReceiveRestock(c *gin.Context) {
	var request request.RestockReceiptRequest
	c.BindJSON(&request)
}
func CreatePurchaseOrder(c *gin.Context) {
	var request request.PurchaseOrderRequest
	c.BindJSON(&request)
}
func GetStockInfo(c *gin.Context) {
}
func GetRestockOrderLog(c *gin.Context) {
}
func GetPurchaseOrderLog(c *gin.Context) {
}
func GetItemInventoryReport(c *gin.Context) {
}
func GetSalesReport(c *gin.Context) {
}
