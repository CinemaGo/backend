package routes

import (
	"cinemaGo/backend/api/handlers"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	router := gin.Default()

	router.GET("/", handlers.MainPage)

	return router
}
