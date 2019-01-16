package model

import (
	"database/sql"
	"time"
)

type (
	Item struct {
		sku string
		name string
		stock int
		createdAt time.Time
	}

	RestockOrder struct {
		Id int
		InvoiceId string
		Quantity int
		Price int32
		SKU string
		Status string
		CreatedAt time.Time
	}

	RestockReception struct {
		Id int
		RestockOrderId int
		DateReceived time.Time
		Quantity int
	}

	PurchaseOrder struct {
		id int
		orderId string
		createdAt time.Time
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

func (order RestockOrder) Persist() int {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("INSERT INTO restock_orders (invoice_id, quantity, price, sku) VALUES (?, ?, ?, ?)")
	res,err := statement.Exec(order.InvoiceId, order.Quantity, order.Price, order.SKU)
	if err != nil {
		println("Exec err:", err.Error())
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			println("Error:", err.Error())
		}
		return int(id)
	}
	return -1
}

func (order RestockOrder) UpdateStatus() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("UPDATE restock_orders set status = ? where invoice_id = ?")
	statement.Exec(order.Status, order.InvoiceId)
}

func GetIdByInvoiceId(invoiceId string) int {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	rows, _ := database.Query("SELECT id FROM restock_orders where invoice_id = ?", invoiceId)
	var id int
	rows.Next()
	rows.Scan(&id)
	return id
}

func (order RestockReception) Persist() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("INSERT INTO restock_receptions (restock_order_id, quantity, date_received) VALUES (?, ?, ?)")
	statement.Exec(order.RestockOrderId, order.Quantity, order.DateReceived)
}

func (order PurchaseOrder) Persist() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("INSERT INTO purchase_orders (order_id, sku, item_name, quantity, price, notes) VALUES (?, ?, ?, ?, ?, ?)")
	statement.Exec(order.orderId, order.sku, order.itemName, order.quantity, order.price, order.notes)
}
