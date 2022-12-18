package main

import (
	"net/http"
	"os"

	"github.com/dai65527/weather-slack-app/eventapi/handler"
	"github.com/dai65527/weather-slack-app/slackhandler"
	"github.com/slack-go/slack"
)

func main() {
	oauthToken := os.Getenv("SLACK_OAUTH_TOKEN")
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")

	api := slack.New(oauthToken)

	slackHandler := slackhandler.SlackHandler{
		Api: api,
	}

	http.Handle("/slack/events", &handler.EventHandler{
		SlackHandler:  &slackHandler,
		SigningSecret: signingSecret,
	})

	http.Handle("/slack/interaction", &handler.InteractivityHandler{
		SlackHandler: &slackHandler,
	})

	http.ListenAndServe(":8080", nil)
}
