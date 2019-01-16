package restock

import (
	"time"
	"app/util"
)

type RestockOrder struct {
	Id int
	InvoiceId string
	Quantity int
	Price int32
	SKU string
	Status string
	CreatedAt time.Time
}


func (order RestockOrder) Persist() int {
	statement, _ := util.Database.Prepare("INSERT INTO restock_orders (invoice_id, quantity, price, sku, status) VALUES (?, ?, ?, ?, ?)")
	res,err := statement.Exec(order.InvoiceId, order.Quantity, order.Price, order.SKU, order.Status)
	defer statement.Close()

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
	statement, _ := util.Database.Prepare("UPDATE restock_orders set status = ? where id = ?")
	statement.Exec(order.Status, order.Id)
	statement.Close()
}

func GetByInvoiceId(invoiceId string) RestockOrder {
	rows, _ := util.Database.Query("SELECT id, status, quantity FROM restock_orders where invoice_id = ?", invoiceId)
	var id, quantity int
	var status string
	rows.Next()
	rows.Scan(&id)
	rows.Scan(&quantity)
	rows.Scan(&status)
	defer rows.Close()

	return RestockOrder{
		Id: id,
		Status: status,
		Quantity: quantity,
	}
}

