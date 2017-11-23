package constants

var AcceptedCurrencies = []string{
	"AUD",
	"USD",
}

var AcceptedTokens = map[string]string{
	"BCC": "bitconnect",
	"BCH": "bitcoin-cash",
	"BTC": "bitcoin",
	"DASH": "dash",
	"DOGE": "dogecoin",
	"ETC": "ethereum-classic",
	"ETH": "ethereum",
	"LTC": "litecoin",
	"UDST": "tether",
	"XRP": "ripple",
	"XMR": "monero",
}

const CoinMarketCapBaseURl = "https://api.coinmarketcap.com/v1/ticker/"