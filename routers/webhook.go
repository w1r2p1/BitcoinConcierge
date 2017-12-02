package routers

import (
	"BitcoinBot/controllers"
	"github.com/gorilla/mux"
)

func SetWebhookRoute(router *mux.Router) *mux.Router {
	router.HandleFunc("/webhook", controllers.SelectHandler).Methods("POST")
	return router
}
