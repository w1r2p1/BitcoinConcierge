package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"BitcoinBot/constants"
	"BitcoinBot/types"
	"BitcoinBot/utils"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

func getPrice(ticker, currency string) (string, error) {
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

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var token, currency string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	args := strings.Split(text, " ")

	if len(args) > 2 {
		rtm.SendMessage(rtm.NewOutgoingMessage("Enter the ticker and currency only", msg.Channel))
		return
	} else {
		token, currency = utils.GetTokenAndCurrency(args)
	}

	if !utils.IsAcceptedToken(token) {
		rtm.SendMessage(rtm.NewOutgoingMessage("Your cryptocurrency is not supported", msg.Channel))
		return
	} else if !utils.IsAcceptedCurrency(currency) {
		rtm.SendMessage(rtm.NewOutgoingMessage("Your currency is not supported. Please select USD or AUD", msg.Channel))
		return
	} else if (utils.IsAcceptedToken(token)) && (utils.IsAcceptedCurrency(currency)) {
		price, err := getPrice(token, currency)
		if err != nil {
			log.Printf("%+v\n", err)
			rtm.SendMessage(rtm.NewOutgoingMessage("Internal Server Error", msg.Channel))
		} else {
			log.Println("price:", price)
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprint(strings.ToUpper(token), " $", price), msg.Channel))
		}
		return
	} else {
		rtm.SendMessage(rtm.NewOutgoingMessage("What you are trying to do is not supported", msg.Channel))
		return
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
				log.Println("Bot is connected")

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
