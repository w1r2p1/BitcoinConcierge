package models

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