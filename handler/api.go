package handler

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"ijah-store/domain/item"
	"ijah-store/domain/purchase"
	"ijah-store/domain/request"
	"ijah-store/domain/restock"
	"ijah-store/util"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ApiHandler struct {
	itemSvc item.Service
	purchaseSvc purchase.Service
	restockSvc restock.Service
}

func NewApiHandler(itemSvc item.Service, purchaseSvc purchase.Service, restockSvc restock.Service) *ApiHandler {
	return &ApiHandler{
		itemSvc: itemSvc,
		purchaseSvc: purchaseSvc,
		restockSvc: restockSvc,
	}
}

func (h *ApiHandler) CreateRestockOrder(c *gin.Context) {
	var req request.RestockOrderRequest
	err := c.ShouldBindJSON(&req)

	var dateReceived time.Time
	if req.DateReceived != "" {
		dateReceived, err = util.ParseDateFromDefault(req.DateReceived)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restockOrder, errorMsg := h.restockSvc.SaveRestockOrder(req)
	if errorMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	h.itemSvc.CreateItemIfNotAny(req)
	h.restockSvc.SaveRestockReception(restockOrder, dateReceived, req.QuantityReceived)

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}

func (h *ApiHandler) ReceiveRestock(c *gin.Context) {
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
	if request.Quantity == 0 || request.DateReceived == "" || request.InvoiceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input ada yang kosong, tolong cek kembali"})
	}

	h.restockSvc.ReceiveRestock(request.InvoiceId, request.Quantity, dateReceived)

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}

func (h *ApiHandler) CreatePurchaseOrder(c *gin.Context) {
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

	itemMap := h.itemSvc.GetRequestedItemMap(requestedSkuList)
	errorMsg := h.purchaseSvc.CheckAvailability(req.Items, itemMap)

	if errorMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	h.purchaseSvc.HandlePurchase(req, itemMap)

	c.JSON(http.StatusOK, gin.H{"message" : "sukses"})
}

func (h *ApiHandler) GetStockInfo(c *gin.Context) {
	items := h.itemSvc.GetAllItems()

	c.JSON(http.StatusOK, items)
}

func (h *ApiHandler) GetRestockOrderLog(c *gin.Context) {
	restockLogData := h.restockSvc.GetAllRestockLog()

	c.JSON(http.StatusOK, restockLogData)
}

func (h *ApiHandler) GetPurchaseOrderLog(c *gin.Context) {
	orders := h.purchaseSvc.GetAllOrders()

	c.JSON(http.StatusOK, orders)
}

func (h *ApiHandler) GetItemInventoryReport(c *gin.Context) {
	//prep data
	items := h.itemSvc.GetAllItems()
	stockInfoMap := h.restockSvc.GetItemStockInfoMap()

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

func (h *ApiHandler) GetSalesReport(c *gin.Context) {
	purchasedOrders := h.purchaseSvc.GetAllOrders()
	itemStockInfoMap := h.restockSvc.GetItemStockInfoMap()

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
