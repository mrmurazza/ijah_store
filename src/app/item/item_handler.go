package item

import (
	"app/request"
	"app/restock"
)

type StockInfo struct {
	TotalQuantity       int
	TotalPurchasedPrice int32
}

func (stockInfo StockInfo) AVGPrice() int64 {
	return int64(stockInfo.TotalPurchasedPrice) / int64(stockInfo.TotalQuantity)
}

func CreateItemIfNotAny(request request.RestockOrderRequest) {
	item := Item{
		SKU: request.SKU,
		Name: request.ItemName,
		// stock will be updated separately
	}

	if !item.IsExist() {
		item.Persist()
	}
}

func UpdateItemStock(sku string, quantity int) {
	item := GetItem(sku)

	item.Stock += quantity
	item.UpdateStock()
}

func GetRequestedItemMap(requestedSkuList []string) map[string]Item{
	// get requested item data and more prep
	var itemMap = make(map[string]Item)
	items := GetItems(requestedSkuList[:])
	for _, product := range items {
		itemMap[product.SKU] = product
	}
	return itemMap
}

func GetItemStockInfoMap() map[string]StockInfo {
	restockLogs := restock.GetAllRestockLog()

	stockInfoMap := make(map[string]StockInfo)
	for _, log := range restockLogs {
		stockInfo := stockInfoMap[log.SKU]
		stockInfo.TotalPurchasedPrice += log.TotalPrice
		stockInfo.TotalQuantity += log.Quantity
		stockInfoMap[log.SKU] = stockInfo
	}

	return stockInfoMap
}
