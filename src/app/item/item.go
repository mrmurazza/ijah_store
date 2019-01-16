package item

import (
	"time"
	"fmt"
	"app/util"
	"strings"
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
	row := util.Database.QueryRow("SELECT name, stock FROM items where sku = ?", sku)
	var (
		name string
		stock int
	)

	err := row.Scan(&name, &stock)
	if err != nil {
		println("Exec err:", err.Error())
	}

	return Item{
		SKU: sku,
		Name: name,
		Stock: stock,
	}
}

func GetItems(skus []string) []Item {
	query := "SELECT sku, name, stock FROM items where sku in (?" + strings.Repeat(",?", len(skus)-1) + ")"
	rows, _ := util.Database.Query(query, skus)
	var (
		sku, name string
		stock int
		items []Item
	)
	for rows.Next() {
		rows.Scan(&sku)
		rows.Scan(&name)
		rows.Scan(&stock)
		items = append(items, Item{
			SKU: sku,
			Name: name,
			Stock: stock,
		})
	}
	defer rows.Close()

	return items
}

