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

func (b *Bot) ProcessMessage(chatID int, message string) {
	switch message {
	case "/start", "Help":
		b.handleStartOrHelpCommand(chatID)
	case "Pay Asrat":
		b.SendMessage(chatID, "Please enter the amount", ReplyKeyboardMarkup{})
	default:
		b.handleDefault(chatID)
	}
}

func (b *Bot) handleStartOrHelpCommand(chatID int) {
	data, err := os.ReadFile("welcome.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of file:")
	msg := string(data)
	fmt.Println()

	b.SendMessage(chatID, msg, ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{{Text: "Pay Asrat"}},
			{{Text: "Help"}},
		},
	})
}

func (b *Bot) handleDefault(chatID int) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env")
	}
	amount, err := strconv.Atoi(Text)
	fmt.Print(amount < 56)
	if err != nil || amount < 56 {
		fmt.Println(err)
		b.SendMessage(chatID, "Invalid Information, please enter a valid amount", ReplyKeyboardMarkup{})
		return
	}
	if ChatID == chatID {
		invoice := SendInvoiceRequest{
			ChatID:        chatID,
			Title:         "Test Payment",
			Description:   "Payment For Asrat",
			Payload:       "blah blah blah",
			ProviderToken: os.Getenv("PROVIDER_TOKEN"),
			Currency:      "ETB",
			Prices: []LabeledPrice{
				{Label: "some label for test", Amount: amount * 100},
			},
		}
		b.handlePayments(invoice)
	}
}

func (b *Bot) handlePayments(invoice SendInvoiceRequest) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendInvoice", b.Token)
	invoiceMessage := invoice

	jsonInvoice, err := json.Marshal(invoiceMessage)
	if err != nil {
		log.Fatal("Error encoding JSON of Invoice Request: ", err)
	}

	invoiceReader := bytes.NewReader(jsonInvoice)
	res, err := http.Post(apiURL, "application/json", invoiceReader)

	if err != nil {
		log.Fatal("Failed to send Invoice: ", err)
	}

	defer res.Body.Close()

	var response TelegramResponse
	json.NewDecoder(res.Body).Decode(&response)
	if response.Ok {
		fmt.Println(response.Result)
	}
}

func (b *Bot) SendMessage(chatID int, text string, markup ReplyKeyboardMarkup) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.Token)

	var message interface{}

	if markup.Keyboard == nil {
		message = SendMessageRequest{
			ChatID: chatID,
			Text:   text,
		}
	} else {
		message = SendMessageButtonRequest{
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
