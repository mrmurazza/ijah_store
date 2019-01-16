package item

import (
	"time"
	"fmt"
	"app/util"
)

type Item struct {
	SKU string
	Name string
	Stock int
	createdAt time.Time
}

func (item Item) Persist() {
	statement, _ := util.Database.Prepare("INSERT INTO items (sku, name, stock) VALUES (?, ?, ?)")
	statement.Exec(item.SKU, item.Name, item.Stock)
	statement.Close()
}

func (item Item) UpdateStock() {
	statement, _ := util.Database.Prepare("UPDATE items set stock = ? where sku = ?")
	statement.Exec(item.Stock, item.SKU)
	statement.Close()
}

func (item Item) IsExist() bool {
	rows, _ := util.Database.Query("SELECT count(*) FROM restock_orders where invoice_id = ?", item.SKU)
	var counter int
	rows.Next()
	rows.Scan(&counter)
	defer rows.Close()

	fmt.Printf("%d",counter)
	return counter > 0
}

func GetItem(sku string) Item {
	rows, _ := util.Database.Query("SELECT name, stock FROM items where sku = ?", sku)
	var (
		name string
		stock int
	)
	rows.Next()
	rows.Scan(&name)
	rows.Scan(&stock)
	defer rows.Close()

	return Item{
		SKU: sku,
		Name: name,
		Stock: stock,
	}
}

