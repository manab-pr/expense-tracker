package routes

import (
	controllers "expanse-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/v1/api/users/getuser/:user_id", controllers.GetUSer())
	incomingRoutes.POST("/v1/api/users/signup", controllers.SignUp())
	incomingRoutes.POST("/v1/api/users/signin", controllers.Login())
	incomingRoutes.POST("/v1/api/users/signout", controllers.Login())

}
