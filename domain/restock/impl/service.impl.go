package impl

import (
	"ijah-store/domain/item"
	"ijah-store/domain/request"
	"ijah-store/domain/response"
	"ijah-store/domain/restock"
	"time"
)

type service struct {
	repo    restock.Repository
	itemSvc item.Service
}

func NewService(repo restock.Repository) restock.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) validateRestockReq(req request.RestockOrderRequest) string {
	if req.SKU == "" || req.Price == 0 || req.Quantity == 0 || req.ItemName == "" {
		return "data input wajib ada yang kosong"
	}

	if req.InvoiceId == "" && req.QuantityReceived != req.Quantity {
		return "Ketika kwitansi hilang, semua barang harus diterima saat ini juga"
	}

	return ""
}

func (s *service) SaveRestockOrder(req request.RestockOrderRequest) (restock.Order, string) {
	errorMsg := s.validateRestockReq(req)

	if errorMsg != "" {
		return restock.Order{}, errorMsg
	}

	orderStatus := "pending"
	if req.Quantity == req.QuantityReceived {
		orderStatus = "finish"
	}

	if req.InvoiceId == "" {
		req.InvoiceId = "(Hilang)"
	}

	o := restock.Order{
		InvoiceId: req.InvoiceId,
		Quantity:  req.Quantity,
		Price:     req.Price,
		SKU:       req.SKU,
		Status:    orderStatus,
	}
	s.repo.PersistOrder(&o)

	return o, ""
}

func (s *service) ReceiveRestock(invoiceId string, quantity int, dateReceived time.Time) {
	order := s.repo.GetOrderByInvoiceId(invoiceId)
	totalReceivedQuantity := s.repo.CountReceivedStock(order.Id)

	errorMsg := s.validateRequest(order, totalReceivedQuantity, quantity)
	if errorMsg != "" {
		return
	}

	s.SaveRestockReception(order, dateReceived, quantity)
	if quantity+totalReceivedQuantity == order.Quantity {
		order.Status = "finish"
		s.repo.UpdateOrderStatus(&order)
	}}

func (s *service) SaveRestockReception(restockOrder restock.Order, dateReceived time.Time, quantity int) {
	if quantity == 0 {
		return
	}

	restockReception := restock.Reception{
		RestockOrderId: restockOrder.Id,
		DateReceived:   dateReceived,
		Quantity:       quantity,
	}
	s.repo.PersistReception(&restockReception)

	s.itemSvc.UpdateItemStock(restockOrder.SKU, quantity)
}

func (s *service) validateRequest(restockOrder restock.Order, totalQuantity, quantity int  ) string {
	var errorMsg string

	if restockOrder.Status == "finish" {
		errorMsg = "permintaan restock untuk kwitansi ini sudah terpenuhi semua"
	}

	if restockOrder.Quantity < quantity {
		errorMsg = "input quantity tidak valid"
	}

	if restockOrder.Quantity-totalQuantity < quantity {
		errorMsg = "barang yang diterima lebih banyak daripada jumlah sisa permintaan dari kwitansi ini"
	}

	return errorMsg
}

func (s *service) getReceiptMapAndTotalReceived(receptions []restock.Reception) (map[int][]response.ReceiptDetail, map[int]int) {
	receiptDetailMap := make(map[int][]response.ReceiptDetail)
	receivedTotalMap := make(map[int]int)

	for _, reception := range receptions {
		receiptDetailMap[reception.RestockOrderId] = append(receiptDetailMap[reception.RestockOrderId], response.ReceiptDetail{
			ReceivedAt: reception.DateReceived.Format(time.RFC850),
			Quantity:   reception.Quantity,
		})

		receivedTotalMap[reception.RestockOrderId] += reception.Quantity
	}

	return receiptDetailMap, receivedTotalMap
}

func (s *service) GetAllRestockLog() []response.RestockOrderResponse {
	orders := s.repo.GetAllOrders()
	receptions := s.repo.GetAllReceptions()

	receiptDetailMap, receivedStockMap := s.getReceiptMapAndTotalReceived(receptions)
	requestedSkuList := make([]string, len(orders))
	for i, order := range orders {
		requestedSkuList[i] = order.SKU
	}
	itemMap := s.itemSvc.GetRequestedItemMap(requestedSkuList)

	var responses []response.RestockOrderResponse
	for _, order := range orders {
		responses = append(responses, response.RestockOrderResponse{
			SKU:                   order.SKU,
			ItemName:              itemMap[order.SKU].Name,
			InvoiceId:             order.InvoiceId,
			Price:                 order.Price,
			Quantity:              order.Quantity,
			TotalPrice:            order.Price * int32(order.Quantity),
			ReceivedQuantityTotal: receivedStockMap[order.Id],
			ReceiptDetail:         receiptDetailMap[order.Id],
			Status:                order.Status,
		})
	}

	return responses
}

func (s *service) GetItemStockInfoMap() map[string]restock.StockInfo {
	restockLogs := s.GetAllRestockLog()

	stockInfoMap := make(map[string]restock.StockInfo)
	for _, log := range restockLogs {
		stockInfo := stockInfoMap[log.SKU]
		stockInfo.TotalPurchasedPrice += log.TotalPrice
		stockInfo.TotalQuantity += log.Quantity
		stockInfoMap[log.SKU] = stockInfo
	}

	return stockInfoMap
}
