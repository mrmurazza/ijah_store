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
	row := util.Database.QueryRow("SELECT sum(quantity) FROM restock_receptions where restock_order_id = ? group by restock_order_id", restockOrderId)
	var totalQuantity int
	err := row.Scan(&totalQuantity)
	if err != nil {
		println("Exec err:", err.Error())
	}

	return totalQuantity
}
