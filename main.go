package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"BitcoinBot/common"
	"BitcoinBot/routers"
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

	if strings.ToLower(args[0]) == "global" && len(args) == 2 {
		rtm.SendMessage(rtm.NewOutgoingMessage("global command can only be used only. Try `global`", msg.Channel))
		return
	}

	token, currency = common.GetTokenAndCurrency(args)

	if strings.ToLower(args[0]) == "global" {
		var tokenDescriptionArr []string
		tokenPrices, err := common.GetTokenPrices(strings.ToUpper(currency))
		if err != nil {
			log.Printf("%+v\n", err)
			rtm.SendMessage(rtm.NewOutgoingMessage("Internal Server Error", msg.Channel))
			return
		}

		for _, token := range tokenPrices {
			tokenDescription, err := token.Description()
			if err != nil {
				log.Printf("%+v\n", err)
				continue
			}
			tokenDescriptionArr = append(tokenDescriptionArr, tokenDescription)
		}

		response := strings.Join(tokenDescriptionArr, "\n")
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
		return
	}

	if !common.IsAcceptedTicker(token) {
		rtm.SendMessage(rtm.NewOutgoingMessage("Your cryptocurrency is not supported", msg.Channel))
		return
	}

	if !common.IsAcceptedCurrency(currency) {
		rtm.SendMessage(rtm.NewOutgoingMessage("Your currency is not supported. Please select USD or AUD", msg.Channel))
		return
	}

	if (common.IsAcceptedTicker(token)) && (common.IsAcceptedCurrency(currency)) {
		price, err := common.GetPrice(token, currency)
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

func runSlackBot(wg *sync.WaitGroup) {
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
	wg.Done()
}

func runWebServer(wg *sync.WaitGroup) {
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	server := &http.Server{
		Addr:    ":8080",
		Handler: n,
	}
	log.Println("Server is listening on port 8080")
	server.ListenAndServe()
	wg.Done()
}

func main() {
	log.Println("Running both webserver and slack bot")
	var wg sync.WaitGroup
	wg.Add(1)
	go runWebServer(&wg)
	wg.Add(1)
	go runSlackBot(&wg)
	wg.Wait()
}
