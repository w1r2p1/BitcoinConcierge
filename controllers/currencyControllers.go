package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"BitcoinBot/common"
	"BitcoinBot/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func GetPrice(w http.ResponseWriter, r *http.Request) {
	var ticker string
	var isAcceptedCrypto bool
	vars := mux.Vars(r)
	cryptocurrency := vars["cryptocurrency"]
	queries := r.URL.Query()
	currency := queries.Get("unit")

	if currency == "" {
		currency = "AUD"
	}

	if ticker, isAcceptedCrypto = common.IsAcceptedToken(cryptocurrency); !isAcceptedCrypto {
		common.DisplayError(w, nil, "This cryptocurrency is not supported", 404)
	}

	if !common.IsAcceptedCurrency(currency) {
		common.DisplayError(w, nil, "This currency is not supported", 404)
	}

	price, err := common.GetPrice(ticker, strings.ToUpper(currency))
	if err != nil {
		common.DisplayError(w, err, "Failed to get price of the token", 500)
	}
	tokenPrice := models.TokenPrice{
		Ticker: ticker,
		Price:  price,
		Unit:   currency,
	}

	j, err := json.Marshal(tokenPrice)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	} else {
		common.DisplayError(
			w,
			errors.Wrap(err, "Failed to parse json"),
			"Failed to parse response from upstream server",
			500,
		)
	}
}

func GetPrices(w http.ResponseWriter, r *http.Request) {
	// TODO
}
