package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// represents the request sent to the Telegram API for sending messages
type SendMessageRequest struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type SendMessageButtonRequest struct {
	ChatID      int                 `json:"chat_id"`
	Text        string              `json:"text"`
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}

type SendMessageInlineRequest struct {
	ChatID      int                  `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

type SendInvoiceRequest struct {
	ChatID        int            `json:"chat_id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Payload       string         `json:"payload"`
	ProviderToken string         `json:"provider_token"`
	Currency      string         `json:"currency"`
	Prices        []LabeledPrice `json:"prices"`
}

type AnswerPreCheckoutQueryRequest struct {
	PreCheckQueryID string `json:"pre_checkout_query_id"`
	Ok              bool   `json:"ok"`
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

func (b *Bot) ProcessMessage(message string) {
	switch message {
	case "/start", "Help":
		b.handleStartOrHelpCommand()
	case "Pay Asrat":
		b.SendMessage("Please enter the amount", ReplyKeyboardMarkup{})
	case "":
		if update.Message.SuccessfulPayment.InvoicePayload == "payment-payload" {
			b.SendMessage("Thank you for paying!", ReplyKeyboardMarkup{})
		}
	default:
		b.handleDefault(message)
	}
}

func (b *Bot) handleStartOrHelpCommand() {
	data, err := os.ReadFile("welcome.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	msg := string(data)

	b.SendMessage(msg, ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{{Text: "Pay Asrat"}},
			{{Text: "Help"}},
		},
	})
}

func (b *Bot) handleDefault(message string) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env")
	}
	if amount, err := strconv.Atoi(message); err == nil && amount >= 56 {
		b.handlePayment(amount)
		return
	}
	b.SendMessage("Invalid Information, please enter a valid amount", ReplyKeyboardMarkup{})
}

func (b *Bot) handlePayment(amount int) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendInvoice", b.Token)

	invoice := SendInvoiceRequest{
		ChatID:        ChatID,
		Title:         "Asrat",
		Description:   "Payment For Asrat",
		Payload:       "payment-payload",
		ProviderToken: os.Getenv("PROVIDER_TOKEN"),
		Currency:      "ETB",
		Prices: []LabeledPrice{
			{Label: "Sub Total", Amount: amount * 100},
		},
	}

	jsonInvoice, err := json.Marshal(invoice)
	if err != nil {
		log.Fatal("Error encoding JSON of Invoice Request: ", err)
	}

	invoiceReader := bytes.NewReader(jsonInvoice)
	res, err := http.Post(apiURL, "application/json", invoiceReader)

	if err != nil {
		log.Fatal("Failed to send Invoice: ", err)
	}
	defer res.Body.Close()
}

func (b *Bot) AnswerPreCheckoutQuery(preCheckoutQueryID string, ok bool, errorMessage string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/answerPreCheckoutQuery", b.Token)

	response := AnswerPreCheckoutQueryRequest{
		PreCheckQueryID: preCheckoutQueryID,
		Ok:              ok,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		log.Println("Error encoding JSON for answerPreCheckoutQuery response:", err)
		return
	}

	responseReader := bytes.NewReader(responseData)
	res, err := http.Post(apiURL, "application/json", responseReader)
	if err != nil {
		log.Println("Failed to send answerPreCheckoutQuery response:", err)
		return
	}
	defer res.Body.Close()
}

func (b *Bot) SendMessage(text string, markup ReplyKeyboardMarkup) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.Token)

	var message interface{}

	if markup.Keyboard == nil {
		message = SendMessageRequest{
			ChatID: ChatID,
			Text:   text,
		}
	} else {
		message = SendMessageButtonRequest{
			ChatID:      ChatID,
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

func (b *Bot) SendMessageInline(chatID int, text string, markup InlineKeyboardMarkup) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.Token)
	var message interface{}

	if markup.InlineKeyboard == nil {
		message = SendMessageRequest{
			ChatID: chatID,
			Text:   text,
		}
	} else {
		message = SendMessageInlineRequest{
			ChatID:      chatID,
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
