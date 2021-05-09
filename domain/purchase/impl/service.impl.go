package impl

import (
	"ijah-store/domain/item"
	"ijah-store/domain/purchase"
	"ijah-store/domain/request"
)

type service struct {
	repo purchase.Repository
	itemSvc item.Service
}

func NewService(repo purchase.Repository) purchase.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CheckAvailability(itemDetails []request.ItemDetail, itemMap map[string]*item.Item) string {
	var errorMsg string
	for _, itemDetail := range itemDetails {
		product := itemMap[itemDetail.SKU]

		switch {
		case product.SKU == "":
			errorMsg = "produk dengan SKU " + itemDetail.SKU + " tidak ditemukan"
		case product.Stock < itemDetail.Quantity:
			errorMsg = "produk dengan SKU " + product.SKU + " tidak memiliki stok yang cukup"
		default:
			errorMsg = ""
		}

	}
	return errorMsg
}

// save purchase order and reduce stock
func (s *service) HandlePurchase(req request.PurchaseOrderRequest, itemMap map[string]*item.Item) {
	for _, itemDetail := range req.Items {
		product := itemMap[itemDetail.SKU]
		purchaseOrder := purchase.Order{
			OrderId: req.OrderId,
			SKU: product.SKU,
			Quantity: itemDetail.Quantity,
			ItemName: product.Name,
			Price: itemDetail.Price,
			Notes: req.Notes,
		}
		s.repo.Persist(&purchaseOrder)

		product.Stock -= itemDetail.Quantity
		s.itemSvc.UpdateItemStock(product.SKU, product.Stock)
	}
}

func (s *service) GetAllOrders() []purchase.Order {
	return s.repo.GetAllOrders()
}
