package purchase

import (
	"time"
)

type Order struct {
	id int
	OrderId string
	CreatedAt time.Time
	SKU string
	ItemName string
	Quantity int
	Price int32
	Notes string
}
