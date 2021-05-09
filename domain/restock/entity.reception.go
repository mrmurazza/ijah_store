package restock

import (
	"ijah-store/pkg"
	"time"
)

type Reception struct {
	Id int
	RestockOrderId int
	DateReceived time.Time
	Quantity int
}

func (r Reception) Persist() {
	statement, _ := pkg.Database.Prepare("INSERT INTO restock_receptions (restock_order_id, quantity, date_received) VALUES (?, ?, ?)")
	statement.Exec(r.RestockOrderId, r.Quantity, r.DateReceived.Format(time.RFC3339))
	statement.Close()
}

func CountReceivedStock(restockOrderId int) int {
	row := pkg.Database.QueryRow("SELECT sum(quantity) FROM restock_receptions where restock_order_id = ? group by restock_order_id", restockOrderId)
	var totalQuantity int
	err := row.Scan(&totalQuantity)
	if err != nil {
		println("Exec err:", err.Error())
	}

	return totalQuantity
}

func GetAllReceptions() []Reception {
	query := "SELECT restock_order_id, date_received, quantity FROM restock_receptions"
	// converting list of string to args
	rows, err := pkg.Database.Query(query)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var receptions []Reception
	for rows.Next() {
		reception := Reception{}

		err = rows.Scan(&reception.RestockOrderId, &reception.DateReceived, &reception.Quantity)
		if err != nil {
			println("Exec err:", err.Error())
		}

		receptions = append(receptions, reception)
	}
	return receptions
}
