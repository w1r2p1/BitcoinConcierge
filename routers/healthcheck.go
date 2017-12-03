package routers

import (
	"BitcoinBot/controllers"
	"github.com/gorilla/mux"
)

func SetHealthCheckRoute(router *mux.Router) *mux.Router {
	router.HandleFunc("/healthcheck", controllers.HealthcheckHandler).Methods("GET")
	return router
}
