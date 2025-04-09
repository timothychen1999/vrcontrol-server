package main

import (
	"github.com/gin-gonic/gin"

)

func main() {

	router := createRouter()
	router.Run()

}

func createRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	return router
}
