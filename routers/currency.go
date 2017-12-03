package routers

import (
	"BitcoinBot/controllers"
	"github.com/gorilla/mux"
)

func SetCurrencyRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/price", controllers.GetPrices).Methods("GET")
	router.HandleFunc("/price/{cryptocurrency}", controllers.GetPrice).Methods("GET")
	return router
}
