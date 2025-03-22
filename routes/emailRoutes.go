package routes

import (
	handler "expanse-tracker/handlers"

	"github.com/gin-gonic/gin"
)

func EmailRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/v1/api/email/send-otp", handler.SendOTPHandler())
	incomingRoutes.POST("/v1/api/email/verify-otp", handler.VerifyOTPHandler())
	incomingRoutes.POST("/v1/api/email/forgot-password", handler.ForgotPasswordHandler())
	incomingRoutes.POST("/v1/api/email/reset-password", handler.ResetPasswordHandler())

}
