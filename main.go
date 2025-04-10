package main

import (
	"github.com/gin-gonic/gin"
	"github.com/timothychen1999/vrcontrol-server/routes"
)

func main() {

	router := createRouter()
	router.Run()

}

func createRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	ws := router.Group("/ws")
	routes.SetClientWsRoutes(ws)
	return router
}
