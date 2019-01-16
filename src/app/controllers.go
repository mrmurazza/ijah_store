package app

import (
	"github.com/gin-gonic/gin"
	"app/request"
	"net/http"
	"app/item"
	"app/restock"
	"app/util"
	"time"
)

func CreateRestockOrder(c *gin.Context) {
	var request request.RestockOrderRequest
	err := c.ShouldBindJSON(&request)

	var dateReceived time.Time
	if request.DateReceived != "" {
		dateReceived, err = util.ParseDateFromDefault(request.DateReceived)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restockOrder, errorMsg := restock.GenerateRestockOrder(request)
	if errorMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	restockOrderId := restockOrder.Persist()
	item.CreateItemIfNotAny(request)

	if request.QuantityReceived > 0 && request.Quantity != request.QuantityReceived {
		restock.SaveRestockReception(restockOrderId, dateReceived, request.QuantityReceived)
	}
	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}

func ReceiveRestock(c *gin.Context) {
	var request request.RestockReceiptRequest
	c.BindJSON(&request)
	dateReceived, err := util.ParseDateFromDefault(request.DateReceived)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.InvoiceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kwitansi harus ada"})
		return
	}

	restockOrder := restock.GetByInvoiceId(request.InvoiceId)
	totalReceivedQuantity := restock.CountReceivedStock(restockOrder.Id)

	errorMsg := restock.ValidateRequest(request, restockOrder, totalReceivedQuantity)
	if errorMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	restock.SaveRestockReception(restockOrder.Id, dateReceived, request.Quantity)
	restock.HandleStatusUpdate(request, totalReceivedQuantity, restockOrder)
	item.UpdateItemStock(restockOrder.SKU, request.Quantity)

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
