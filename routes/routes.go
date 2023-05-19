package routes

import (
	"example/controllers"
	"example/middleware"

	"github.com/julienschmidt/httprouter"
)

func RegisterRoutes(router *httprouter.Router) {
	router.GET("/", controllers.Hello)
	router.GET("/spots", middleware.ValidateQueryParams(controllers.GetSpots))
}
