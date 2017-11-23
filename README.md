# BitcoinBot
**`WIP`**

A bot in Slack that can update a channel about current price of cryptocurrencies. This can be done via slash commands or prices are notified to the channel hourly

How to run this bot:
```
brew install glide
glide install
SLACK_TOKEN="XXX" go run main.go
```

To run inside a docker container:
```$xslt
export SLACK_TOKEN="XXX"
bin/run.sh
```
