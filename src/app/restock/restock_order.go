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

		order.Id = int(id)
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
	rows:= util.Database.QueryRow("SELECT id, status, quantity, sku, price FROM restock_orders where invoice_id = ?", invoiceId)
	var (
		id, quantity int
		price int32
		status, sku string
	)

	err := rows.Scan(&id, &status, &quantity, &sku, &price)
	if err != nil {
		println("Exec err:", err.Error())
	}

	return RestockOrder{
		Id: id,
		Status: status,
		Quantity: quantity,
		Price: price,
		SKU: sku,
	}
}

