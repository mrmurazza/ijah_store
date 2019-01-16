package util

import "time"

const layout_date = "2006/01/02"

func ParseDateFromDefault(date string) (time.Time, error){
	return time.Parse("2006/01/02", date)
}
