package response

type (
	ReceiptDetail struct {
		ReceivedAt string `json:"receivedAt"`
		Quantity int `json:"quantityReceived"`
	}

	RestockOrderResponse struct {
		SKU string `json:"sku"`
		ItemName string `json:"itemName"`
		InvoiceId string `json:"invoiceId"`
		Price int32 `json:"price"`
		Quantity int `json:"quantity"`
		TotalPrice int32 `json:"totalPrice"`
		ReceivedQuantityTotal int `json:"receivedQuantityTotal"`
		ReceiptDetail []ReceiptDetail `json:"receptions"`
	}
)