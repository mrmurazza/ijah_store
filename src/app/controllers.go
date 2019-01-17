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
	"os"
	"encoding/csv"
	"strconv"
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

	itemMap := item.GetRequestedItemMap(requestedSkuList)
	errorMsg := purchase.CheckAvailability(req.Items, itemMap)

	if errorMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	purchase.HandlePurchase(req, itemMap)

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}

func GetStockInfo(c *gin.Context) {
	items := item.GetAllItems()

	c.JSON(http.StatusOK, items)
}

func GetRestockOrderLog(c *gin.Context) {
	restockLogData := restock.GetAllRestockLog()

	c.JSON(http.StatusOK, restockLogData)
}

func GetPurchaseOrderLog(c *gin.Context) {
	orders := purchase.GetAllOrders()

	c.JSON(http.StatusOK, orders)
}

func GetItemInventoryReport(c *gin.Context) {
	//prep data
	items := item.GetAllItems()
	stockInfoMap := restock.GetItemStockInfoMap()

	var (
		data [][]string
		totalStock int
		totalValue int64
	)

	for _, it := range items {
		stockInfo := stockInfoMap[it.SKU]
		currValue := int64(stockInfo.AVGPrice()) * int64(it.Stock)
		row := []string{
			it.SKU,
			it.Name,
			strconv.Itoa(it.Stock),
			strconv.FormatInt(stockInfo.AVGPrice(), 10),
			strconv.FormatInt(currValue, 10),
		}

		data = append(data, row)
		totalStock += it.Stock
		totalValue += currValue
	}

	file, err := os.Create("inventory_report_"+time.Now().Format("02_01_06")+".csv")
	defer file.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"LAPORAN NILAI BARANG"})
	err = writer.Write([]string{})
	err = writer.Write([]string{"Tanggal Cetak", time.Now().Format("02 January 2006")})
	err = writer.Write([]string{"Jumlah SKU", strconv.Itoa(len(items))})
	err = writer.Write([]string{"Jumlah Total Barang", strconv.Itoa(totalStock)})
	err = writer.Write([]string{"Total Nilai", strconv.FormatInt(totalValue, 10)})
	err = writer.Write([]string{})
	err = writer.Write([]string{"SKU", "Nama Item", "Jumlah", "Rata-Rata Harga Beli", "Total"})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}

func GetSalesReport(c *gin.Context) {
	purchasedOrders := purchase.GetAllOrders()
	itemStockInfoMap := restock.GetItemStockInfoMap()

	var (
		omzet, grossValue int64
		totalSold int
		data [][]string
	)

	for _, order := range purchasedOrders {
		stockInfo := itemStockInfoMap[order.SKU]
		totalPrice := int64(order.Price) * int64(order.Quantity)
		value := (int64(order.Price) - stockInfo.AVGPrice()) * int64(order.Quantity)
		row := []string{
			order.OrderId,
			order.CreatedAt.Format(time.RFC822),
			order.SKU,
			order.ItemName,
			strconv.Itoa(order.Quantity),
			strconv.FormatInt(int64(order.Price), 10),
			strconv.FormatInt(totalPrice, 10),
			strconv.FormatInt(stockInfo.AVGPrice(), 10),
			strconv.FormatInt(value, 10),
		}

		data = append(data, row)
		omzet += totalPrice
		grossValue += value
		totalSold += order.Quantity
	}

	file, err := os.Create("sales_report_"+time.Now().Format("02_01_06")+".csv")
	defer file.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"LAPORAN PENJUALAN"})
	err = writer.Write([]string{})
	err = writer.Write([]string{"Tanggal Cetak", time.Now().Format("02 January 2006")})
	err = writer.Write([]string{"Total Omzet", strconv.FormatInt(omzet, 10)})
	err = writer.Write([]string{"Total Laba Kotor", strconv.FormatInt(grossValue, 10)})
	err = writer.Write([]string{"Total Penjualan", strconv.Itoa(len(purchasedOrders))})
	err = writer.Write([]string{"Total Barang Terjual", strconv.Itoa(totalSold)})
	err = writer.Write([]string{})
	err = writer.Write([]string{"ID Pesanan", "Waktu", "SKU", "Nama Barang", "Jumlah", "Harga Jual", "Total", "Harga Beli", "Laba"})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}
