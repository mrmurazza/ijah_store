package request

type (
	RestockOrderRequest struct {
		SKU string `json:"sku"`
		ItemName string `json:"itemName"`
		InvoiceId string `json:"invoiceId"`
		Price int32 `json:"price"`
		Quantity int `json:"quantity"`
		DateReceived string `json:"dateReceived"`
		QuantityReceived int `json:"quantityReceived"`
	}

	RestockReceiptRequest struct {
		InvoiceId string `json:"invoiceId"`
		DateReceived string `json:"dateReceived"`
		Quantity int `json:"quantity"`
	}

	PurchaseOrderRequest struct {
		OrderId string `json:"orderId"`
		SKU string `json:"sku"`
		ItemName string `json:"itemName"`
		Quantity int `json:"quantity"`
		Price int32 `json:"price"`
		Notes string `json:"notes"`
	}
)
