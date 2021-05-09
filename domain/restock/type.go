package restock

import (
	"ijah-store/domain/request"
	"ijah-store/domain/response"
	"time"
)

type Repository interface {
	PersistReception(r *Reception)
	CountReceivedStock(restockOrderId int) int
	GetAllReceptions() []Reception

	PersistOrder(o *Order) int
	UpdateOrderStatus(o *Order)
	GetOrderByInvoiceId(invoiceId string) Order
	GetAllOrders() []Order
}

type Service interface {
	SaveRestockOrder(req request.RestockOrderRequest) (Order, string)
	SaveRestockReception(restockOrder Order, dateReceived time.Time, quantity int)
	ReceiveRestock(invoiceId string, quantity int, dateReceived time.Time)
	GetAllRestockLog() []response.RestockOrderResponse
	GetItemStockInfoMap() map[string]StockInfo
}
