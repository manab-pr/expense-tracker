package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AccountType string
type BankName string

const (
	Bank   AccountType = "Bank"
	Wallet AccountType = "Wallet"
)

const (
	SBI       BankName = "SBI"
	PNB       BankName = "PNB"
	CITY_BANK BankName = "CITY_BANK"
	YES_BANK  BankName = "YES_BANK"
	PAYTM     BankName = "PAYTM"
)

type Account struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"name"`
	Account_Type *AccountType       `json:"account_type" validate:"account_type"`
	Bank_Name    *BankName          `json:"bank_name,omitempty" validate:"bank_name"`
	Amount       *float64           `json:"amount"`
}
