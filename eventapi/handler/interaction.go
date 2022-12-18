package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dai65527/weather-slack-app/weather"
	"github.com/slack-go/slack"
)

type InteractivityHandler struct {
	SlackClient *slack.Client
}

func (h *InteractivityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var interaction slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &interaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(interaction.ActionCallback.BlockActions) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	action := interaction.ActionCallback.BlockActions[0]
	switch action.ActionID {
	case "select_city":
		weather, err := weather.GetWeather(action.SelectedOption.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _, _, err = h.SlackClient.SendMessage(
			"",
			slack.MsgOptionReplaceOriginal(interaction.ResponseURL),
			slack.MsgOptionBlocks(
				slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "どの都市の天気を調べますか？: " + action.SelectedOption.Text.Text,
					},
				},
				slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "```\n" + weather + "```",
					},
				},
			),
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
