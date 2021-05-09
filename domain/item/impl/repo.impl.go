package impl

import (
	"database/sql"
	"ijah-store/domain/item"
	"strings"
)

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) item.Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) Persist(i item.Item) {
	statement, _ := r.db.Prepare("INSERT INTO items (sku, name, stock) VALUES (?, ?, ?)")
	statement.Exec(i.SKU, i.Name, i.Stock)
	statement.Close()
}

func (r *repo) UpdateStock(i item.Item) {
	statement, _ := r.db.Prepare("UPDATE items set stock = ? where sku = ?")
	statement.Exec(i.Stock, i.SKU)
	statement.Close()
}

func (r *repo) IsExist(i item.Item) bool {
	rows, _ := r.db.Query("SELECT count(*) FROM restock_orders where invoice_id = ?", i.SKU)
	var counter int
	rows.Next()
	rows.Scan(&counter)
	defer rows.Close()

	return counter > 0
}

func (r *repo) GetItem(sku string) item.Item {
	row := r.db.QueryRow("SELECT name, stock FROM items where sku = ?", sku)
	var (
		name string
		stock int
	)

	err := row.Scan(&name, &stock)
	if err != nil {
		println("Exec err:", err.Error())
	}

	return item.Item{
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

func (r *repo) GetItems(skuList []string) []*item.Item {
	query := "SELECT sku, name, stock FROM list where sku in (?" + strings.Repeat(",?", len(skuList)-1) + ")"
	args := convertStringListToInterface(skuList)
	rows, err := r.db.Query(query, args...)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var list []*item.Item
	for rows.Next() {
		i := item.Item{}

		err = rows.Scan(&i.SKU, &i.Name, &i.Stock)
		if err != nil {
			println("Exec err:", err.Error())
		}

		list = append(list, &i)
	}
	return list
}

func (r *repo) GetAllItems() []*item.Item {
	query := "SELECT sku, name, stock FROM list"
	rows, err := r.db.Query(query)
	defer rows.Close()

	if err != nil {
		println("Exec err:", err.Error())
	}

	var list []*item.Item
	for rows.Next() {
		i := item.Item{}

		err = rows.Scan(&i.SKU, &i.Name, &i.Stock)
		if err != nil {
			println("Exec err:", err.Error())
		}

		list = append(list, &i)
	}
	return list
}
