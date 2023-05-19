package routes

import (
	"example/controllers"

	"github.com/julienschmidt/httprouter"
)

func RegisterRoutes(router *httprouter.Router) {
	router.GET("/", controllers.Hello)
	router.GET("/spots", controllers.GetSpots)
}
