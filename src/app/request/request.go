package request

import "time"

type (
	RestockOrderRequest struct {
		sku string
		itemName string
		invoiceId string
		price int32
		quantity int
		dateReceived time.Time
		quantityReceived int
	}

	RestockReceiptRequest struct {
		invoiceId string
		dateReceived time.Time
		quantity int
	}

	PurchaseOrderRequest struct {
		orderId string
		sku string
		itemName string
		quantity int
		price int32
		notes string
	}
)
