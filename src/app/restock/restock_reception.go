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

func GetAllReceptions() []RestockReception {
	query := "SELECT restock_order_id, date_received, quantity FROM restock_receptions"
	// converting list of string to args
	rows, err := util.Database.Query(query)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var receptions []RestockReception
	for rows.Next() {
		reception := RestockReception{}

		err = rows.Scan(&reception.RestockOrderId, &reception.DateReceived, &reception.Quantity)
		if err != nil {
			println("Exec err:", err.Error())
		}

		receptions = append(receptions, reception)
	}
	return receptions
}
