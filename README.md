# BitcoinConcierge
**`WIP`**

A cryptocurrency helper that can be integrated with Google Assistant, Stride, Slack ... to track prices and volatility. Users can transfer coins on chat channels.

## Usage

How to run this bot on Slack:
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

## Roadmap
* Add unit and integration tests
* Integrate with Google Assitant
* Transfer coins on chat channels
* User deep learning to make predictions about future prices
