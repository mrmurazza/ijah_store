package restock

import (
	"ijah-store/pkg"
	"time"
)

type Order struct {
	Id int
	InvoiceId string
	Quantity int
	Price int32
	SKU string
	Status string
	CreatedAt time.Time
}


func (o *Order) Persist() int {
	statement, _ := pkg.Database.Prepare("INSERT INTO restock_orders (invoice_id, quantity, price, sku, status) VALUES (?, ?, ?, ?, ?)")
	res,err := statement.Exec(o.InvoiceId, o.Quantity, o.Price, o.SKU, o.Status)
	defer statement.Close()

	if err != nil {
		println("Exec err:", err.Error())
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			println("Error:", err.Error())
		}

		o.Id = int(id)
		return int(id)
	}
	return -1
}

func (o Order) UpdateStatus() {
	statement, _ := pkg.Database.Prepare("UPDATE restock_orders set status = ? where id = ?")
	statement.Exec(o.Status, o.Id)
	statement.Close()
}

func GetByInvoiceId(invoiceId string) Order {
	rows:= pkg.Database.QueryRow("SELECT id, status, quantity, sku, price FROM restock_orders where invoice_id = ?", invoiceId)
	var (
		id, quantity int
		price int32
		status, sku string
	)

	err := rows.Scan(&id, &status, &quantity, &sku, &price)
	if err != nil {
		println("Exec err:", err.Error())
	}

	return Order{
		Id: id,
		Status: status,
		Quantity: quantity,
		Price: price,
		SKU: sku,
	}
}

func GetAllOrders() []Order {
	query := "SELECT id, sku, quantity, price, invoice_id FROM restock_orders"
	// converting list of string to args
	rows, err := pkg.Database.Query(query)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var orders []Order
	for rows.Next() {
		order := Order{}

		err = rows.Scan(&order.Id, &order.SKU, &order.Quantity, &order.Price, &order.InvoiceId)
		if err != nil {
			println("Exec err:", err.Error())
		}

		orders = append(orders, order)
	}
	return orders
}

