package util

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func InitDatabase() {
	Database, _ = sql.Open("sqlite3", "ijah_store.db")
	statement, _ := Database.Prepare("CREATE TABLE IF NOT EXISTS items (" +
		"sku varchar(20) primary key, " +
		"name varchar(50) not null, " +
		"stock integer not null default 0, " +
		"created_at datetime default current_timestamp);")
	statement.Exec()

	statement, _ = Database.Prepare("CREATE TABLE IF NOT EXISTS restock_orders ( " +
		"id integer primary key, " +
		"invoice_id varchar(20) default '(Hilang)', " +
		"quantity integer not null," +
		"price biginteger not null, " +
		"SKU varchar(20) not null, " +
		"status varchar(10) default 'pending', " +
		"created_at datetime default current_timestamp);")
	statement.Exec()

	statement, _ = Database.Prepare("CREATE TABLE IF NOT EXISTS restock_receptions ( " +
		"id integer primary key, " +
		"restock_order_id integer, " +
		"quantity integer not null," +
		"date_received datetime default current_timestamp);")
	statement.Exec()

	statement, _ = Database.Prepare("CREATE TABLE IF NOT EXISTS purchase_orders ( " +
		"id integer primary key, " +
		"order_id varchar(20), " +
		"quantity integer not null," +
		"price biginteger not null, " +
		"SKU varchar(20) not null, " +
		"item_name varchar(50) not null, " +
		"notes text, " +
		"created_at datetime default current_timestamp);")
	statement.Exec()

	statement.Close()
}
