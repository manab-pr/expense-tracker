package models

import (
	"time"
)

type Transaction struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Type        Type      `bson:"type" json:"type"`
	Amount      float64   `bson:"amount" json:"amount"`
	Description string    `bson:"description" json:"description"`
	CreatedDate time.Time `bson:"created_date" json:"created_date"`
}
