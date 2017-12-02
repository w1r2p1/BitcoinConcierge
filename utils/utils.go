package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"BitcoinBot/constants"
	"BitcoinBot/types"
	"github.com/pkg/errors"
)

func IsAcceptedToken(ticker string) bool {
	ticker = strings.ToUpper(ticker)
	if _, ok := constants.AcceptedTokens[ticker]; ok {
		return true
	}
	return false
}

func IsAcceptedCurrency(currency string) bool {
	currency = strings.ToUpper(currency)
	for _, val := range constants.AcceptedCurrencies {
		if currency == val {
			return true
		}
	}
	return false
}

func GetPrice(ticker, currency string) (string, error) {
	var token types.TokenInfo
	ticker = strings.ToUpper(ticker)
	tokenID := constants.AcceptedTokens[ticker]
	url := fmt.Sprint(constants.CoinMarketCapBaseURl, tokenID, "/?convert=", currency)

	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "Failed to contact coinmarket server")
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read the response body")
	}

	if currency == "AUD" || currency == "aud" {
		return token[0].PriceAud, nil
	}
	return token[0].PriceUsd, nil
}

func GetPrices(currency string) (types.TokenPrices, error) {
	// TODO
}

func GetTokenAndCurrency(args []string) (string, string) {
	var token, currency string
	token = args[0]
	if len(args) == 1 {
		currency = "AUD"
	} else {
		currency = args[1]
	}
	return token, currency
}
