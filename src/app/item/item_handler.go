package item

import (
	"app/request"
)

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
