package app

import (
	"github.com/gin-gonic/gin"
	"app/request"
	"net/http"
	"app/model"
)

func CreateRestockOrder(c *gin.Context) {
	var request request.RestockOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.InvoiceId == "" && request.QuantityReceived <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ketika kwitansi hilang barang harus diterima saat ini juga"})
		return
	}

	orderStatus := "pending"
	if request.Quantity == request.QuantityReceived {
		orderStatus = "finish"
	}

	restockOrder := model.RestockOrder{
		InvoiceId: request.InvoiceId,
		Quantity: request.Quantity,
		Price: request.Price,
		SKU: request.SKU,
		Status: orderStatus,
	}
	restockOrderId := restockOrder.Persist()

	if request.QuantityReceived > 0 {
		restockReception := model.RestockReception{
			RestockOrderId: restockOrderId,
			DateReceived:   request.DateReceived,
			Quantity:       request.QuantityReceived,
		}
		restockReception.Persist()
	}

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}
func ReceiveRestock(c *gin.Context) {
	var request request.RestockReceiptRequest
	c.BindJSON(&request)

	if request.InvoiceId == "" || request.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input invalid"})
		return
	}

	restockOrderId := model.GetIdByInvoiceId(request.InvoiceId)
	restockReception := model.RestockReception{
		RestockOrderId: restockOrderId,
		DateReceived: request.DateReceived,
		Quantity: request.Quantity,
	}
	restockReception.Persist()

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}
func CreatePurchaseOrder(c *gin.Context) {
	var request request.PurchaseOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
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
