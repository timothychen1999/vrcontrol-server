package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/timothychen1999/vrcontrol-server/controller"
)

func SetClientWsRoutes(router *gin.RouterGroup) {
	router.GET("/client/:clientId", controller.ConnectToRoomSocket)
	router.GET("/control/:roomId", controller.ConnectToRoomControlSocket)
}
