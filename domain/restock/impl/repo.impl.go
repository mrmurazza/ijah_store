package impl

import (
	"database/sql"
	"ijah-store/domain/restock"
	"ijah-store/pkg"
	"time"
)

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) restock.Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) PersistReception(rc *restock.Reception) {
	statement, _ := pkg.Database.Prepare("INSERT INTO restock_receptions (restock_order_id, quantity, date_received) VALUES (?, ?, ?)")
	statement.Exec(rc.RestockOrderId, rc.Quantity, rc.DateReceived.Format(time.RFC3339))
	statement.Close()
}

func (r *repo) CountReceivedStock(restockOrderId int) int {
	row := pkg.Database.QueryRow("SELECT sum(quantity) FROM restock_receptions where restock_order_id = ? group by restock_order_id", restockOrderId)
	var totalQuantity int
	err := row.Scan(&totalQuantity)
	if err != nil {
		println("Exec err:", err.Error())
	}

	return totalQuantity
}

func (r *repo) GetAllReceptions() []restock.Reception {
	query := "SELECT restock_order_id, date_received, quantity FROM restock_receptions"
	// converting list of string to args
	rows, err := pkg.Database.Query(query)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var receptions []restock.Reception
	for rows.Next() {
		reception := restock.Reception{}

		err = rows.Scan(&reception.RestockOrderId, &reception.DateReceived, &reception.Quantity)
		if err != nil {
			println("Exec err:", err.Error())
		}

		receptions = append(receptions, reception)
	}
	return receptions
}

func (r *repo) PersistOrder(o *restock.Order) int {
	statement, _ := pkg.Database.Prepare("INSERT INTO restock_orders (invoice_id, quantity, price, sku, status) VALUES (?, ?, ?, ?, ?)")
	res, err := statement.Exec(o.InvoiceId, o.Quantity, o.Price, o.SKU, o.Status)
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

func (r *repo) UpdateOrderStatus(o *restock.Order) {
	statement, _ := pkg.Database.Prepare("UPDATE restock_orders set status = ? where id = ?")
	statement.Exec(o.Status, o.Id)
	statement.Close()
}

func (r *repo) GetOrderByInvoiceId(invoiceId string) restock.Order {
	rows := pkg.Database.QueryRow("SELECT id, status, quantity, sku, price FROM restock_orders where invoice_id = ?", invoiceId)
	var (
		id, quantity int
		price        int32
		status, sku  string
	)

	err := rows.Scan(&id, &status, &quantity, &sku, &price)
	if err != nil {
		println("Exec err:", err.Error())
	}

	return restock.Order{
		Id:       id,
		Status:   status,
		Quantity: quantity,
		Price:    price,
		SKU:      sku,
	}
}

func (r *repo) GetAllOrders() []restock.Order {
	query := "SELECT id, sku, quantity, price, invoice_id FROM restock_orders"
	// converting list of string to args
	rows, err := pkg.Database.Query(query)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var orders []restock.Order
	for rows.Next() {
		order := restock.Order{}

		err = rows.Scan(&order.Id, &order.SKU, &order.Quantity, &order.Price, &order.InvoiceId)
		if err != nil {
			println("Exec err:", err.Error())
		}

		orders = append(orders, order)
	}
	return orders
}
