package item

import (
	"ijah-store/pkg"
	"strings"
	"time"
)

type Item struct {
	SKU string
	Name string
	Stock int
	createdAt time.Time
}

func (item Item) Persist() {
	statement, _ := pkg.Database.Prepare("INSERT INTO items (sku, name, stock) VALUES (?, ?, ?)")
	statement.Exec(item.SKU, item.Name, item.Stock)
	statement.Close()
}

func (item Item) UpdateStock() {
	statement, _ := pkg.Database.Prepare("UPDATE items set stock = ? where sku = ?")
	statement.Exec(item.Stock, item.SKU)
	statement.Close()
}

func (item Item) IsExist() bool {
	rows, _ := pkg.Database.Query("SELECT count(*) FROM restock_orders where invoice_id = ?", item.SKU)
	var counter int
	rows.Next()
	rows.Scan(&counter)
	defer rows.Close()

	return counter > 0
}

func GetItem(sku string) Item {
	row := pkg.Database.QueryRow("SELECT name, stock FROM items where sku = ?", sku)
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

func convertStringListToInterface(list []string) []interface{} {
	args := []interface{}{}
	for _, el := range list {
		args = append(args, el)
	}

	return args
}

func GetItems(skuList []string) []Item {
	query := "SELECT sku, name, stock FROM items where sku in (?" + strings.Repeat(",?", len(skuList)-1) + ")"
	args := convertStringListToInterface(skuList)
	rows, err := pkg.Database.Query(query, args...)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var items []Item
	for rows.Next() {
		item := Item{}

		err = rows.Scan(&item.SKU, &item.Name, &item.Stock)
		if err != nil {
			println("Exec err:", err.Error())
		}

		items = append(items, item)
	}
	return items
}

func GetAllItems() []Item {
	query := "SELECT sku, name, stock FROM items"
	rows, err := pkg.Database.Query(query)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var items []Item
	for rows.Next() {
		item := Item{}

		err = rows.Scan(&item.SKU, &item.Name, &item.Stock)
		if err != nil {
			println("Exec err:", err.Error())
		}

		items = append(items, item)
	}
	return items
}

