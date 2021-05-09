package impl

import (
	"database/sql"
	"ijah-store/domain/purchase"
	"ijah-store/pkg"
)

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) purchase.Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) Persist(order *purchase.Order) {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("INSERT INTO purchase_orders (order_id, sku, item_name, quantity, price, notes) VALUES (?, ?, ?, ?, ?, ?)")
	statement.Exec(order.OrderId, order.SKU, order.ItemName, order.Quantity, order.Price, order.Notes)
}

func (r *repo) GetAllOrders() []purchase.Order {
	query := "SELECT created_at, sku, item_name, quantity, price, order_id, notes FROM purchase_orders"
	// converting list of string to args
	rows, err := pkg.Database.Query(query)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var orders []purchase.Order
	for rows.Next() {
		order := purchase.Order{}

		err = rows.Scan(&order.CreatedAt, &order.SKU, &order.ItemName, &order.Quantity, &order.Price, &order.OrderId, &order.Notes)
		if err != nil {
			println("Exec err:", err.Error())
		}

		orders = append(orders, order)
	}
	return orders
}
