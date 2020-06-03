package handler

import (
	iqc "infoqerja-line/app/config"
	iql "infoqerja-line/app/line"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

// LineBotHandler will handle all line's callback request
type LineBotHandler struct {
	config iqc.Config
	bot    iql.BotClient
}

// BuildLineBotHandler returns LineBotHandler struct
func BuildLineBotHandler(config iqc.Config, bot iql.BotClient) *LineBotHandler {
	return &LineBotHandler{
		config: config,
		bot:    bot,
	}
}

// Callback will handle the callback from line
func (h LineBotHandler) Callback(w http.ResponseWriter, r *http.Request) {
	events, err := h.bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				if IsValidCommand(message.Text) {
					command := CheckCommand(message.Text)
					switch command {
					case "help":
						if err = h.HelpCommand(event.ReplyToken); err != nil {
							log.Print(err)
						}
					default:
						if err = h.InvalidCommand(event.ReplyToken); err != nil {
							log.Print(err)
						}
					}
				} else {
					if event.Source.Type == linebot.EventSourceTypeUser {
						if err = h.UnknownHandler(event.ReplyToken); err != nil {
							log.Print(err)
						}
					}
				}
			}
		}

		if event.Type == linebot.EventTypeFollow {
			if err = h.WelcomeHandler(event.ReplyToken); err != nil {
				log.Print(err)
			}
		}
	}
}
