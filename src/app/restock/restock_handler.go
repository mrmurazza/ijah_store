package restock

import (
	"app/request"
	"time"
)

func GenerateRestockOrder(request request.RestockOrderRequest) (RestockOrder, string) {
	if request.InvoiceId == "" && request.QuantityReceived <= 0 {
		return RestockOrder{}, "Ketika kwitansi hilang, semua barang harus diterima saat ini juga"
	}

	orderStatus := "pending"
	if request.Quantity == request.QuantityReceived {
		orderStatus = "finish"
	}

	return RestockOrder{
		InvoiceId: request.InvoiceId,
		Quantity: request.Quantity,
		Price: request.Price,
		SKU: request.SKU,
		Status: orderStatus,
	}, ""
}

func SaveRestockReception(restockOrderId int, dateReceived time.Time, quantity int) {
	restockReception := RestockReception{
		RestockOrderId: restockOrderId,
		DateReceived:   dateReceived,
		Quantity:       quantity,
	}
	restockReception.Persist()
}

func ValidateRequest(request request.RestockReceiptRequest, restockOrder RestockOrder, totalQuantity int) string {
	var errorMsg string
	if restockOrder.Status == "finish" {
		errorMsg = "permintaan restock untuk kwitansi ini sudah terpenuhi semua"
	}

	if restockOrder.Quantity < request.Quantity {
		errorMsg = "input quantity tidak valid"
	}

	if restockOrder.Quantity - totalQuantity < request.Quantity {
		errorMsg = "barang yang diterima lebih banyak daripada jumlah sisa permintaan dari kwitansi ini"
	}

	return errorMsg
}

func HandleStatusUpdate(request request.RestockReceiptRequest, existingQuantity int, order RestockOrder){
	if request.Quantity +existingQuantity == order.Quantity {
		order.Status = "finish"
		order.UpdateStatus()
	}
}
