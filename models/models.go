package models

import (
	"fmt"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

type Tokens []TokenInfo

type TokenInfo struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUsd         string `json:"price_usd"`
	PriceBtc         string `json:"price_btc"`
	Two4HVolumeUsd   string `json:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	MaxSupply        string `json:"max_supply"`
	PercentChange1H  string `json:"percent_change_1h"`
	PercentChange24H string `json:"percent_change_24h"`
	PercentChange7D  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
	PriceAud         string `json:"price_aud"`
	Two4HVolumeAud   string `json:"24h_volume_aud"`
	MarketCapAud     string `json:"market_cap_aud"`
}

type TokenPrice struct {
	Ticker string `json:"ticker"`
	Price  string `json:"price"`
	Unit   string `json:"unit"`
}

func (t TokenPrice) Description() (string, error) {
	roundedPrice, err := roundUp(t)
	if err != nil {
		return "", err
	}
	return t.Ticker + " $" + roundedPrice + " " + t.Unit, nil
}

func (t TokenPrice) RoundUpPrice() (string, error) {
	return roundUp(t)
}

func roundUp(t TokenPrice) (string, error) {
	price, err := strconv.ParseFloat(t.Price, 64)
	if err != nil {
		return "", errors.Wrap(err, "Failed to convert string to float64")
	}
	roundedPrice := round(price, 0.5, 2)
	return fmt.Sprintf("%v", roundedPrice), nil
}

func round(val float64, roundOn float64, places int) float64 {

	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)

	var round float64
	if val > 0 {
		if div >= roundOn {
			round = math.Ceil(digit)
		} else {
			round = math.Floor(digit)
		}
	} else {
		if div >= roundOn {
			round = math.Floor(digit)
		} else {
			round = math.Ceil(digit)
		}
	}

	return round / pow
}
