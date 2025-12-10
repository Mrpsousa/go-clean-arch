package entity

import (
	"time"
)

type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
	Updated_at  time.Time
	Created_at  time.Time
}
