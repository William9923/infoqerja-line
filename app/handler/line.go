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
				HandleIncomingMessage(h.bot, event.ReplyToken, message.Text)
			}
		}
		if event.Type == linebot.EventTypeFollow {
			// add welcome handler
		}

		if event.Type == linebot.EventTypePostback {
			data := event.Postback.Data

			if data == "DATE" {
				log.Printf("Successful getting data : (%v)", *&event.Postback.Params.Date)
			}
			HandleIncomingMessage(h.bot, event.ReplyToken, "!show")
		}
	}
}
