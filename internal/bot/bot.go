package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Bot struct {
	Token      string
	WebHookUrl string
}

type TelegramResponse struct {
	Ok          bool        `json:"ok,omitempty"`
	Result      interface{} `json:"result,omitempty"`
	Description string      `json:"description,omitempty"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

type SendMessageRequest struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type SendMessageWithButtonRequest struct {
	ChatID      int         `json:"chat_id"`
	Text        string      `json:"text"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
}

func (b *Bot) Connect() {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", b.Token)

	body := map[string]string{
		"url": b.WebHookUrl,
	}

	data, err := json.Marshal(body)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	dataReader := bytes.NewReader(data)

	res, err := http.Post(apiURL, "application/json", dataReader)

	if err != nil {
		log.Fatal("POST request error: ", err)
	}

	defer res.Body.Close()

	var resp TelegramResponse
	json.NewDecoder(res.Body).Decode(&resp)
	if resp.Ok {
		fmt.Println("Bot connected successfully!")
	}
}

func (b *Bot) ProcessMessage(chatid int, message string) {
	switch message {
	case "/start":
		msg := "Welcome to this bot!!"
		b.SendMessage(chatid, msg, ReplyMarkup{})
		// b.SendMessage(chatid, msg, ReplyMarkup{
		// 	Keyboard: [][]KeyboardButton{
		// 		{{Text: "Button"}, {Text: "Button1"}},
		// 		{{Text: "Button2"}},
		// 		{{Text: "Button3"}, {Text: "Button4"}},
		// 	},
		// })
	}
}

func (b *Bot) SendMessage(chatid int, text string, markup ReplyMarkup) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.Token)

	var message interface{}

	if markup.Keyboard == nil {
		message = SendMessageRequest{
			ChatID: chatid,
			Text:   text,
		}
	} else {
		message = SendMessageWithButtonRequest{
			ChatID:      chatid,
			Text:        text,
			ReplyMarkup: markup,
		}
	}

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Error encoding JSON: ", err)
	}

	msgReader := bytes.NewReader(jsonMsg)
	res, err := http.Post(apiURL, "application/json", msgReader)
	if err != nil {
		log.Fatal("Failed to send message: ", err)
	}

	defer res.Body.Close()

	var response TelegramResponse
	json.NewDecoder(res.Body).Decode(&response)
	fmt.Println(response.Description)
	if response.Ok {
		fmt.Println("Message sent successfully")
	}
}
