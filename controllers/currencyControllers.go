package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"BitcoinBot/constants"
	"BitcoinBot/types"
	"BitcoinBot/utils"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func GetPrice(w http.ResponseWriter, r *http.Request) {
	var price, ticker string
	var err error
	vars := mux.Vars(r)
	currency := vars["currency"]
	key := r.FormValue("unit")

	for k, val := range constants.AcceptedTokens {
		if val == strings.ToLower(currency) {
			ticker = k
			price, err = utils.GetPrice(ticker, key)
			break
		}
	}

	if err != nil {
		utils.DisplayError(w, err, "Failed to get price of the token", 500)
	}
	tokenPrice := types.TokenPrice{
		Ticker: ticker,
		Price:  price,
		Unit:   currency,
	}

	if j, err := json.Marshal(tokenPrice); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.WriteHeader(j)
	} else {
		utils.DisplayError(
			w,
			errors.Wrap(err, "Failed to parse json"),
			"Unexpected error",
			500,
		)
	}
}

func GetPrices(w http.ResponseWriter, r *http.Request) {
	// TODO
}
