package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/timothychen1999/vrcontrol-server/controller"
)

func SetSimpleControlRoutes(router *gin.RouterGroup) {
	router.GET("/assignseq/:roomId/:clientId/:seq", controller.AssignSequence)
	router.GET("/forcemove/:roomId/:clientId/:dest", controller.ForceMove)
}
