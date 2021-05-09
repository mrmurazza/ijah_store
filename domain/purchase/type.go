package purchase

import (
	"ijah-store/domain/item"
	"ijah-store/domain/request"
)

type Repository interface {
	Persist(*Order)
	GetAllOrders() []Order
}

type Service interface {
	CheckAvailability(itemDetails []request.ItemDetail, itemMap map[string]*item.Item) string
	HandlePurchase(req request.PurchaseOrderRequest, itemMap map[string]*item.Item)
	GetAllOrders() []Order
}
