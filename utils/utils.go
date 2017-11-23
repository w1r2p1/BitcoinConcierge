package utils

import (
	"strings"

	"BitcoinBot/constants"
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

func GetTokenAndCurrency(args []string) (string, string){
	var token, currency string
	token = args[0]
	if len(args) == 1 {
		currency = "AUD"
	} else {
		currency = args[1]
	}
	return token, currency
}