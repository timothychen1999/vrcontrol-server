package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/timothychen1999/vrcontrol-server/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		log.Printf("Error: %v", err)
	}
	router := createRouter()
	router.Run()

}

func createRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	ws := router.Group("/ws")
	routes.SetClientWsRoutes(ws)
	simple := router.Group("/simple")
	routes.SetSimpleControlRoutes(simple)
	control := router.Group("/control")
	routes.SetControlRoute(control)
	return router
}
