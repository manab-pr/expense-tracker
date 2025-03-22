package routes

import (
	controller "expanse-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func IncomeExpanseRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/v1/api/income-expanse/add-income", controller.AddIncome())
	incomingRoutes.POST("/v1/api/income-expanse/add-expense", controller.AddExpanse())
	incomingRoutes.GET("/v1/api/income-expanse/get-months-income", controller.GetTotalIncomeByMonth())
	incomingRoutes.GET("/v1/api/income-expanse/get-months-expense", controller.GetTotalExpanseByMonth())
	incomingRoutes.GET("/v1/api/income-expanse/get-transactions", controller.GetTransactions())

}
