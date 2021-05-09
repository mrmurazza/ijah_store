package restock

import (
	"time"
)

type Reception struct {
	Id int
	RestockOrderId int
	DateReceived time.Time
	Quantity int
}

type Order struct {
	Id int
	InvoiceId string
	Quantity int
	Price int32
	SKU string
	Status string
	CreatedAt time.Time
}

type StockInfo struct {
	TotalQuantity       int
	TotalPurchasedPrice int32
}

func (stockInfo StockInfo) AVGPrice() int64 {
	return int64(stockInfo.TotalPurchasedPrice) / int64(stockInfo.TotalQuantity)
}