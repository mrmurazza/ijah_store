package purchase

import (
	"database/sql"
	"time"
)

type PurchaseOrder struct {
	Id int
	OrderId string
	CreatedAt time.Time
	SKU string
	ItemName string
	Quantity int
	Price int32
	Notes string
}


func (order PurchaseOrder) Persist() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("INSERT INTO purchase_orders (order_id, sku, item_name, quantity, price, notes) VALUES (?, ?, ?, ?, ?, ?)")
	statement.Exec(order.OrderId, order.SKU, order.ItemName, order.Quantity, order.Price, order.Notes)
}
