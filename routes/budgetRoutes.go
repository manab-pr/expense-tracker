package routes

import (
	controller "expanse-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func BudgetRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/v1/api/budget/create-budget", controller.CreateBudget())
	incomingRoutes.POST("/v1/api/budget/update-budget", controller.UpdateBudget())
	incomingRoutes.POST("/v1/api/budget/delete-budget", controller.DeleteBudget())
	incomingRoutes.POST("/v1/api/budget/get-budget", controller.GetBudgetAmount())
	incomingRoutes.POST("/v1/api/budget/remaining-budget", controller.GetRemainingAmount())

}
