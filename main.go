package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"BitcoinBot/routers"
	"BitcoinBot/utils"
	"github.com/nlopes/slack"
	"github.com/urfave/negroni"
)

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
	}

	if strings.ToUpper(args[0]) == "global" {
		// TODO
	}

	token, currency = utils.GetTokenAndCurrency(args)

	if !utils.IsAcceptedToken(token) {
		rtm.SendMessage(rtm.NewOutgoingMessage("Your cryptocurrency is not supported", msg.Channel))
		return
	}

	if !utils.IsAcceptedCurrency(currency) {
		rtm.SendMessage(rtm.NewOutgoingMessage("Your currency is not supported. Please select USD or AUD", msg.Channel))
		return
	}

	if (utils.IsAcceptedToken(token)) && (utils.IsAcceptedCurrency(currency)) {
		price, err := utils.GetPrice(token, currency)
		if err != nil {
			log.Printf("%+v\n", err)
			rtm.SendMessage(rtm.NewOutgoingMessage("Internal Server Error", msg.Channel))
			return
		}

		log.Println("price:", price)
		rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprint(strings.ToUpper(token), " $", price), msg.Channel))
		return
	}

	rtm.SendMessage(rtm.NewOutgoingMessage("What you are trying to do is not supported", msg.Channel))
}

func runSlackBot() {
	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	log.Println("Slack RTM is running")

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

func runWebServer() {
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	server := &http.Server{
		Addr:    ":8080",
		Handler: n,
	}
	log.Println("Server is listening on port 8080")
	server.ListenAndServe()
}

func main() {
	go runSlackBot()
	go runWebServer()
}
