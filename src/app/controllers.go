package app

import (
	"github.com/gin-gonic/gin"
	"app/request"
	"net/http"
	"app/item"
	"app/restock"
	"app/util"
	"time"
	"app/purchase"
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

	restockOrder.Persist()
	item.CreateItemIfNotAny(request)
	restock.SaveRestockReception(restockOrder, dateReceived, request.QuantityReceived)

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

	restock.SaveRestockReception(restockOrder, dateReceived, request.Quantity)
	restock.HandleStatusUpdate(request, totalReceivedQuantity, restockOrder)

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}

func CreatePurchaseOrder(c *gin.Context) {
	var req request.PurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// prep requested items data
	requestedSkuList := make([]string, len(req.Items))
	for i, itemDetail := range req.Items {
		requestedSkuList[i] = itemDetail.SKU
	}

	// get requested item data and more prep
	var itemMap = make(map[string]item.Item)
	items := item.GetItems(requestedSkuList[:])
	for _, product := range items {
		itemMap[product.SKU] = product
	}

	// check availability
	var errorMsg string
	for _, itemDetail := range req.Items {
		product := itemMap[itemDetail.SKU]

		switch {
		case product.SKU == "":
			errorMsg = "produk dengan SKU " + itemDetail.SKU + " tidak ditemukan"
		case product.Stock < itemDetail.Quantity:
			errorMsg = "produk dengan SKU " + product.SKU + " tidak memiliki stok yang cukup"
		default:
			errorMsg = ""
		}

		if errorMsg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
			return
		}
	}


	// save purchase order and reduce stock
	for _, itemDetail := range req.Items {
		product := itemMap[itemDetail.SKU]
		purchaseOrder := purchase.PurchaseOrder{
			OrderId: req.OrderId,
			SKU: product.SKU,
			Quantity: itemDetail.Quantity,
			ItemName: product.Name,
			Price: itemDetail.Price,
			Notes: req.Notes,
		}
		purchaseOrder.Persist()

		product.Stock -= itemDetail.Quantity
		product.UpdateStock()
	}

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
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
