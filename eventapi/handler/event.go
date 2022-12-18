package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type EventHandler struct {
	SlackClient   *slack.Client
	SigningSecret string
}

func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// リクエストの検証
	sv, err := slack.NewSecretsVerifier(r.Header, h.SigningSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// eventをパース
	event, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// URLVerification eventをhandle（EventAPI有効化時に叩かれる）
	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	if event.Type == slackevents.CallbackEvent {
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		// SlackAppがメンションされた時に発火
		case *slackevents.AppMentionEvent:
			_, _, err := h.SlackClient.PostMessage(ev.Channel, slack.MsgOptionBlocks(
				slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: "どの都市の天気を調べますか？",
					},
					Accessory: &slack.Accessory{
						SelectElement: &slack.SelectBlockElement{
							ActionID: "select_city",
							Type:     slack.OptTypeStatic,
							Placeholder: &slack.TextBlockObject{
								Type: slack.PlainTextType,
								Text: "都市を選択",
							},
							Options: []*slack.OptionBlockObject{
								{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "東京"}, Value: "東京"},
								{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "ソウル"}, Value: "ソウル"},
								{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "北京"}, Value: "北京"},
								{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "シドニー"}, Value: "シドニー"},
								{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "パリ"}, Value: "パリ"},
								{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "ロンドン"}, Value: "ロンドン"},
								{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "ベルリン"}, Value: "ベルリン"},
								{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "ニューヨーク"}, Value: "ニューヨーク"},
								{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "ロサンゼルス"}, Value: "ロサンゼルス"},
							},
						},
					},
				},
			))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
