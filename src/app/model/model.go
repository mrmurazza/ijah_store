package model

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"database/sql"
)

type (
	Item struct {
		sku string
		name string
		stock int
		createdAt timestamp.Timestamp
	}

	RestockOrder struct {
		id int
		invoiceId string
		quantity int
		price int32
		sku string
		status string
		createdAt timestamp.Timestamp
	}

	RestockReception struct {
		id int
		restockOrderId int
		dateReceived timestamp.Timestamp
		quantity int
	}

	PurchaseOrder struct {
		id int
		orderId string
		createdAt timestamp.Timestamp
		sku string
		itemName string
		quantity int
		price int32
		notes string
	}
)

func (item Item) Persist() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("INSERT INTO items (sku, name, stock) VALUES (?, ?, ?)")
	statement.Exec(item.sku, item.name, 0)
}

func (item Item) UpdateStock() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("UPDATE items set stock = ? where sku = ?")
	statement.Exec(item.stock, item.sku)
}

func (order RestockOrder) Persist() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("INSERT INTO restock_orders (invoice_id, quantity, price, sku) VALUES (?, ?, ?, ?)")
	statement.Exec(order.invoiceId, order.quantity, order.price, order.sku)
}

func (order RestockOrder) UpdateStatus() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("UPDATE restock_orders set status = ? where invoice_id = ?")
	statement.Exec(order.status, order.invoiceId)
}

func (order RestockReception) Persist() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("INSERT INTO restock_receptions (restock_order_id, quantity, date_received) VALUES (?, ?, ?)")
	statement.Exec(order.restockOrderId, order.quantity, order.dateReceived)
}

func (order PurchaseOrder) Persist() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("INSERT INTO purchase_orders (order_id, sku, item_name, quantity, price, notes) VALUES (?, ?, ?, ?, ?, ?)")
	statement.Exec(order.orderId, order.sku, order.itemName, order.quantity, order.price, order.notes)
}
