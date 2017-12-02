package routers

import (
	"BitcoinBot/controllers"
	"github.com/gorilla/mux"
)

func SetCurrencyRoutes(router *mux.Router) *mux.Router {
	router.Path("/price").Queries("unit", "{unit}").HandlerFunc(controllers.GetPrices).Methods("GET")
	router.Path("/price/{currency}").Queries("unit", "{unit}").HandlerFunc(controllers.GetPrice).Methods("GET")
	return router
}
