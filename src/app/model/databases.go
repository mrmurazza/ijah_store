package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() {
	database, _ := sql.Open("sqlite3", "ijah_store.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS items (" +
		"sku varchar(20) primary key, " +
		"name varchar(50) not null, " +
		"stock integer not null default 0, " +
		"created_at datetime default current_timestamp);")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO items (sku, name, stock) VALUES (?, ?, ?)")
	statement.Exec("SSI-D00791015-LL-BWH", "Zalekia Plain Casual Blouse (L,Broken White)", 154)

	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS restock_orders ( " +
		"id integer auto increment, " +
		"invoice_id varchar(20) default '(Hilang)', " +
		"quantity integer not null," +
		"price biginteger not null, " +
		"SKU varchar(20) not null, " +
		"status varchar(10) default 'pending', " +
		"created_at datetime default current_timestamp);")
	statement.Exec()

	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS restock_receptions ( " +
		"id integer auto increment, " +
		"restock_order_id integer, " +
		"quantity integer not null," +
		"created_at datetime default current_timestamp);")
	statement.Exec()

	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS purchase_orders ( " +
		"id integer auto increment, " +
		"order_id varchar(20), " +
		"quantity integer not null," +
		"price biginteger not null, " +
		"SKU varchar(20) not null, " +
		"item_name varchar(50) not null, " +
		"notes text, " +
		"created_at datetime default current_timestamp);")
	statement.Exec()
}
