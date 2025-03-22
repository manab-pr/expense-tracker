package routes

import (
	controller "expanse-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func AccountRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/v1/api/account/add-account", controller.AddAccount())
	incomingRoutes.POST("/v1/api/account/delete-account/:account_id", controller.DeleteAccount())
}
