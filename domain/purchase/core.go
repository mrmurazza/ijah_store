package purchase

import (
	"ijah-store/domain/item"
	"ijah-store/domain/request"
)

func CheckAvailability(itemDetails []request.ItemDetail, itemMap map[string]item.Item) string {
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
func HandlePurchase(req request.PurchaseOrderRequest, itemMap map[string]item.Item) {
	for _, itemDetail := range req.Items {
		product := itemMap[itemDetail.SKU]
		purchaseOrder := PurchaseOrder{
			OrderId: req.OrderId,
			SKU: product.SKU,
			Quantity: itemDetail.Quantity,
			ItemName: product.Name,
			Price: itemDetail.Price,
			Notes: req.Notes,
		}
		purchaseOrder.Persist()

		product.Stock -= itemDetail.Quantity
		product.UpdateStock()
	}
}
