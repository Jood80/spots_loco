package routes

import (
	"example/controllers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router){
	r.HandleFunc("/", controllers.Hello)
	r.HandleFunc("/test", controllers.Test)
	r.HandleFunc("/spots", controllers.GetSpots)
}
