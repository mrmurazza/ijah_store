package purchase

import (
	"database/sql"
	"time"
	"app/util"
)

type PurchaseOrder struct {
	id int
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

func GetAllOrders() []PurchaseOrder {
	query := "SELECT created_at, sku, item_name, quantity, price, order_id, notes FROM purchase_orders"
	// converting list of string to args
	rows, err := util.Database.Query(query)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var orders []PurchaseOrder
	for rows.Next() {
		order := PurchaseOrder{}

		err = rows.Scan(&order.CreatedAt, &order.SKU, &order.ItemName, &order.Quantity, &order.Price, &order.OrderId, &order.Notes)
		if err != nil {
			println("Exec err:", err.Error())
		}

		orders = append(orders, order)
	}
	return orders
}
