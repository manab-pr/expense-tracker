package utils

import (
	"expanse-tracker/models"

	"github.com/go-playground/validator/v10"
)

func ValidateAccountType(fl validator.FieldLevel) bool {
	account_type := fl.Field().String()

	return account_type == string(models.Bank) || account_type == string(models.Wallet)
}
func ValidateBankName(fl validator.FieldLevel) bool {
	bank_name := fl.Field().String()

	return bank_name == string(models.SBI) || bank_name == string(models.PNB) || bank_name == string(models.CITY_BANK) || bank_name == string(models.PAYTM) || bank_name == string(models.YES_BANK)
}
func ValidateCategory(fl validator.FieldLevel) bool {
	category := fl.Field().Interface().(models.Category)

	return category == models.Food || category == models.Shopping || category == models.Insurance || category == models.Rent || category == models.Subscription
}
