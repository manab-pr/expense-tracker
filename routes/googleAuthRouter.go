package routes

import (
	controller "expanse-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func GoogleAuth(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/v1/api/google/g-login", controller.HandleGoogleLogin)
	incomingRoutes.GET("/v1/api/google/g-callback", controller.HandleGoogleCallback)
}
