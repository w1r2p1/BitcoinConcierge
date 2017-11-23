package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"BitcoinBot/constants"
	"BitcoinBot/utils"
	"BitcoinBot/types"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

func getPrice(ticker, currency string) (string, error){
	var token types.TokenInfo
	url := fmt.Sprint(constants.CoinMarketCapBaseURl, ticker)

	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "Failed to contact coinmarket server")
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read the response body")
	}

	if currency == "AUD" {
		return token.PriceAud, nil
	}
	return token.PriceUsd, nil
}

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var token, currency string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	args := strings.Split(text, " ")

	if len(args) > 2 {
		rtm.SendMessage(rtm.NewOutgoingMessage("Enter the ticker and currency only", msg.Channel))
	} else if len(args) < 1 {
		rtm.SendMessage(rtm.NewOutgoingMessage("Hey I am your BTC assistant", msg.Channel))
	} else {
		token, currency = utils.GetTokenAndCurrency(args)
	}

	if !utils.IsAcceptedToken(token) {
		rtm.SendMessage(rtm.NewOutgoingMessage("Your cryptocurrency is not supported", msg.Channel))
	} else if !utils.IsAcceptedCurrency(currency) {
		rtm.SendMessage(rtm.NewOutgoingMessage("Your currency is not supported. Please select USD or AUD", msg.Channel))
	} else if utils.IsAcceptedToken(text) && utils.IsAcceptedCurrency(currency) {
		price, err := getPrice(token, currency)
		if err != nil {
			fmt.Printf("%+v\n", err)
			rtm.SendMessage(rtm.NewOutgoingMessage("Internal Server Error", msg.Channel))
			return
		}
		rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprint(token, " ", price), msg.Channel))
	} else {
		rtm.SendMessage(rtm.NewOutgoingMessage("What you are trying to do is not supported", msg.Channel))
	}
}

func main() {
	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)
				fmt.Println("Bot ID:", ev.Info.User.ID)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					respond(rtm, ev, prefix)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}
