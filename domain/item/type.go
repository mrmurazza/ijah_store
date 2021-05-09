package item

import "ijah-store/domain/request"

type Repository interface {
	Persist(i Item)
	UpdateStock(i Item)

	IsExist(i Item) bool
	GetItem(sku string) Item
	GetItems(skuList []string) []*Item
	GetAllItems() []*Item
}

type Service interface {
	CreateItemIfNotAny(request request.RestockOrderRequest)
	UpdateItemStock(sku string, quantity int)
	GetRequestedItemMap(requestedSkuList []string) map[string]*Item
	GetAllItems() []*Item

}
