package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Budget struct {
	Id               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Category         Category           `json:"category" validate:"category"`
	Amount           *int64             `json:"amount"`        //default 0
	Recieve_Alert    *bool              `json:"receive_alert"` //default false
	Month            string             `json:"month,omitempty"`
	Alert_Percentage *int               `json:"alert_percentage"`
	Created_at       time.Time          `json:"created_at"`
	Updated_at       time.Time          `json:"updated_at"`
}
