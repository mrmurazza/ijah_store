package app

import (
	"github.com/gin-gonic/gin"
	"app/request"
	"net/http"
	"app/model"
	"time"
)

func CreateItemIfNotAny(request request.RestockOrderRequest) {
	item := model.Item{
		SKU: request.SKU,
		Name: request.ItemName,
	}

	if !item.IsExist() {
		item.Persist()
	}
}

func CreateRestockOrder(c *gin.Context) {
	var request request.RestockOrderRequest
	err := c.ShouldBindJSON(&request)
	t, err := time.Parse("2006/01/02", request.DateReceived)
	if err != nil {
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

	CreateItemIfNotAny(request)

	if request.QuantityReceived > 0 && request.Quantity != request.QuantityReceived {
		restockReception := model.RestockReception{
			RestockOrderId: restockOrderId,
			DateReceived:   t,
			Quantity:       request.QuantityReceived,
		}
		restockReception.Persist()
	}

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}

func ReceiveRestock(c *gin.Context) {
	var request request.RestockReceiptRequest
	c.BindJSON(&request)
	t, err := time.Parse("2006/01/02", request.DateReceived)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.InvoiceId == "" || request.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input invalid"})
		return
	}

	restockOrder := model.GetByInvoiceId(request.InvoiceId)

	var errorMsg string
	if restockOrder.Status == "finish" {
		errorMsg = "permintaan restock untuk kwitansi ini sudah terpenuhi semua"
	}

	if restockOrder.Quantity < request.Quantity {
		errorMsg = "input quantity tidak valid"
	}

	totalReceivedQuantity := model.CountReceivedStock(restockOrder.Id)
	if restockOrder.Quantity - totalReceivedQuantity < request.Quantity {
		errorMsg = "barang yang diterima lebih banyak daripada jumlah sisa permintaan dari kwitansi ini"
	}

	if errorMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	restockReception := model.RestockReception{
		RestockOrderId: restockOrder.Id,
		DateReceived: t,
		Quantity: request.Quantity,
	}
	restockReception.Persist()

	if request.Quantity + totalReceivedQuantity == restockOrder.Quantity {
		restockOrder.Status = "finish"
		restockOrder.UpdateStatus()
	}

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
