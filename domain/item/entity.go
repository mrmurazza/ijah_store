package item

import (
	"time"
)

type Item struct {
	SKU string
	Name string
	Stock int
	createdAt time.Time
}

