package restock

import (
	"time"
	"app/util"
)

type RestockReception struct {
	Id int
	RestockOrderId int
	DateReceived time.Time
	Quantity int
}

func (order RestockReception) Persist() {
	statement, _ := util.Database.Prepare("INSERT INTO restock_receptions (restock_order_id, quantity, date_received) VALUES (?, ?, ?)")
	statement.Exec(order.RestockOrderId, order.Quantity, order.DateReceived.Format(time.RFC3339))
	statement.Close()
}

func CountReceivedStock(restockOrderId int) int {
	rows, _ := util.Database.Query("SELECT sum(quantity) FROM restock_receptions where restock_order_id = ?", restockOrderId)
	var totalQuantity int
	rows.Next()
	rows.Scan(&totalQuantity)
	defer rows.Close()

	return totalQuantity
}
