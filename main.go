package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
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
					results := searchText(message.Text)
					messages := sendingMessages(results)

					if _, err = bot.ReplyMessage(event.ReplyToken, messages...).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func searchText(searchWord string) (results []string) {
	file, err := os.Open("./sasakiazusa.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	for i := 1; sc.Scan(); i++ {
		if err := sc.Err(); err != nil {
			log.Fatal(err)
			break
		}

		text := sc.Text()
		if strings.Contains(text, searchWord) {
			results = append(results, text)
		}
	}

	return results
}

func sendingMessages(lines []string) (messages []linebot.SendingMessage) {
	const maxMessageSize = 5
	messages = make([]linebot.SendingMessage, maxMessageSize)

	for i := 0; i < len(lines); i++ {
		messages = append([]linebot.SendingMessage{linebot.NewTextMessage(lines[i])}, messages...)
	}

	messageLength := maxMessageSize
	if maxMessageSize > len(lines) {
		messageLength = len(lines)
	}
	return messages[:messageLength]
}
