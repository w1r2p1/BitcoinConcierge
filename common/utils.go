package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"BitcoinBot/constants"
	"BitcoinBot/models"
	"github.com/pkg/errors"
)

func IsAcceptedTicker(ticker string) bool {
	ticker = strings.ToUpper(ticker)
	if _, ok := constants.AcceptedTokens[ticker]; ok {
		return true
	}
	return false
}

func IsAcceptedToken(name string) (string, bool) {
	for k, val := range constants.AcceptedTokens {
		if val == strings.ToLower(name) {
			return k, true
		}
	}
	return "", false
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
	var token models.Tokens
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

func GetTokenPrices(currency string) ([]models.TokenPrice, error) {
	var tokens models.Tokens
	var tokenPrices []models.TokenPrice
	url := fmt.Sprint(constants.CoinMarketCapBaseURl, "/?convert=", currency)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to contact coinmarket server")
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&tokens)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read the response body")
	}

	for _, token := range tokens {
		if IsAcceptedTicker(token.Symbol) {
			tokenPrice := models.TokenPrice{
				Ticker: token.Symbol,
				Price: selectPriceField(currency, token),
				Unit: currency,
			}
			tokenPrices = append(tokenPrices, tokenPrice)
		}
	}
	return tokenPrices, nil
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

func selectPriceField(currency string, token models.TokenInfo) string {
	if currency == "AUD" || currency == "aud" {
		return token.PriceAud
	}
	return token.PriceUsd
}
