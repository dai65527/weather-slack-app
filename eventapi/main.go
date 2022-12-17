package main

import (
	"net/http"
	"os"

	"github.com/dai65527/weather-slack-app/eventapi/handler"
	"github.com/slack-go/slack"
)

func main() {
	oauthToken := os.Getenv("SLACK_OAUTH_TOKEN")
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")

	slackClient := slack.New(oauthToken)

	http.Handle("/slack/events", &handler.EventHandler{
		SlackClient:   slackClient,
		OauthToken:    oauthToken,
		SigningSecret: signingSecret,
	})

	http.ListenAndServe(":8080", nil)
}
