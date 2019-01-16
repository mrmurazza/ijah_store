package item

import "app/request"

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
