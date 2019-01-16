package restock

import (
	"app/request"
	"time"
	"app/item"
	"app/response"
)

func validateRestockReq(req request.RestockOrderRequest) string {
	if req.SKU == "" || req.Price == 0 || req.Quantity == 0 || req.ItemName == "" {
		return "data input wajib ada yang kosong"
	}

	if req.InvoiceId == "" && req.QuantityReceived != req.Quantity {
		return "Ketika kwitansi hilang, semua barang harus diterima saat ini juga"
	}

	return ""
}

func GenerateRestockOrder(req request.RestockOrderRequest) (RestockOrder, string) {
	errorMsg := validateRestockReq(req)

	if errorMsg != "" {
		return RestockOrder{}, errorMsg
	}

	orderStatus := "pending"
	if req.Quantity == req.QuantityReceived {
		orderStatus = "finish"
	}

	if req.InvoiceId == "" {
		req.InvoiceId = "(Hilang)"
	}

	return RestockOrder{
		InvoiceId: req.InvoiceId,
		Quantity:  req.Quantity,
		Price:     req.Price,
		SKU:       req.SKU,
		Status:    orderStatus,
	}, ""
}

func SaveRestockReception(restockOrder RestockOrder, dateReceived time.Time, quantity int) {
	if quantity == 0 {
		return
	}

	restockReception := RestockReception{
		RestockOrderId: restockOrder.Id,
		DateReceived:   dateReceived,
		Quantity:       quantity,
	}
	restockReception.Persist()

	item.UpdateItemStock(restockOrder.SKU, quantity)
}

func ValidateRequest(request request.RestockReceiptRequest, restockOrder RestockOrder, totalQuantity int) string {
	var errorMsg string

	if request.Quantity == 0 || request.DateReceived == "" || request.InvoiceId == "" {
		errorMsg = "input ada yang kosong, tolong cek kembali"
	}

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
	if request.Quantity + existingQuantity == order.Quantity {
		order.Status = "finish"
		order.UpdateStatus()
	}
}

func getReceiptMapAndTotalReceived(receptions []RestockReception) (map[int][]response.ReceiptDetail, map[int]int) {
	receiptDetailMap := make(map[int][]response.ReceiptDetail)
	receivedTotalMap := make(map[int]int)

	for _,reception := range receptions {
		receiptDetailMap[reception.RestockOrderId] = append(receiptDetailMap[reception.RestockOrderId], response.ReceiptDetail{
			ReceivedAt: reception.DateReceived.Format(time.RFC850),
			Quantity: reception.Quantity,
		})

		receivedTotalMap[reception.RestockOrderId] += reception.Quantity
	}

	return receiptDetailMap, receivedTotalMap
}

func GetAllRestockLog() []response.RestockOrderResponse {
	orders := GetAllOrders()
	receptions := GetAllReceptions()

	receiptDetailMap, receivedStockMap := getReceiptMapAndTotalReceived(receptions)
	requestedSkuList := make([]string, len(orders))
	for i, order := range orders {
		requestedSkuList[i] = order.SKU
	}
	itemMap := item.GetRequestedItemMap(requestedSkuList)

	var responses []response.RestockOrderResponse
	for _, order := range orders {
		responses = append(responses, response.RestockOrderResponse{
			SKU:           order.SKU,
			ItemName:      itemMap[order.SKU].Name,
			InvoiceId:     order.InvoiceId,
			Price:         order.Price,
			Quantity:      order.Quantity,
			TotalPrice:    order.Price * int32(order.Quantity),
			ReceivedQuantityTotal: receivedStockMap[order.Id],
			ReceiptDetail: receiptDetailMap[order.Id],
		})
	}

	return responses
}