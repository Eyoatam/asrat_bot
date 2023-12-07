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

const (
	apiBaseURL = "https://api.telegram.org/bot%s"
)

var (
	httpClient = &http.Client{}
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

func NewBot(token string, webhookurl string) *Bot {
	return &Bot{
		Token:      token,
		WebHookUrl: webhookurl,
	}
}
func (b *Bot) Connect() {
	apiURL := fmt.Sprintf(apiBaseURL+"/setWebhook", b.Token)
	body := map[string]string{"url": b.WebHookUrl}

	if err := b.postJSON(apiURL, body, "Bot connected successfully!"); err != nil {
		log.Println("Failed to connect:", err)
	}
}

func (b *Bot) postJSON(url string, data interface{}, successMessage string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	response, err := httpClient.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("POST request error: %w", err)
	}
	defer response.Body.Close()

	var resp TelegramResponse
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	if resp.Ok {
		fmt.Println(successMessage)
	}
	return nil
}

func (b *Bot) ProcessMessage(update Update) {
	message := update.Message.Text
	chatID := update.Message.Chat.ID
	switch message {
	case "/start", "/help", "Help":
		b.handleStartOrHelpCommand(chatID)
	case "Pay Asrat":
		b.SendMessage(chatID, "Please enter the amount", ReplyKeyboardMarkup{})
	case "":
		if update.Message.SuccessfulPayment.InvoicePayload == "payment-payload" {
			fmt.Println("The chatID is: ", chatID)
			b.SendMessage(chatID, "Thank you for paying!", ReplyKeyboardMarkup{})
		}
	default:
		b.handleDefault(chatID, message)
	}
}

func (b *Bot) handleStartOrHelpCommand(chatID int) {
	data, err := os.ReadFile("welcome.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	msg := string(data)

	b.SendMessage(chatID, msg, ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{{Text: "Pay Asrat"}},
			{{Text: "Help"}},
		},
	})
}

func (b *Bot) handleDefault(chatID int, message string) {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load .env")
	}

	if amount, err := strconv.Atoi(message); err == nil && amount >= 56 {
		b.handlePayment(chatID, amount)
		return
	}

	b.SendMessage(chatID, "Invalid information. Please check /help for more information.", ReplyKeyboardMarkup{})
}

func (b *Bot) handlePayment(chatID int, amount int) {
	apiURL := fmt.Sprintf(apiBaseURL+"/sendInvoice", b.Token)

	invoice := SendInvoiceRequest{
		ChatID:        chatID,
		Title:         "Asrat",
		Description:   "Payment For Asrat",
		Payload:       "payment-payload",
		ProviderToken: os.Getenv("PROVIDER_TOKEN"),
		Currency:      "ETB",
		Prices: []LabeledPrice{
			{Label: "Sub Total", Amount: amount * 100},
		},
	}

	if err := b.postJSON(apiURL, invoice, "Invoice sent successfully"); err != nil {
		log.Println("Failed to send Invoice:", err)
	}
}

func (b *Bot) SendMessage(chatID int, text string, markup ReplyKeyboardMarkup) {
	apiURL := fmt.Sprintf(apiBaseURL+"/sendMessage", b.Token)

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

	if err := b.postJSON(apiURL, message, "Message sent successfully"); err != nil {
		log.Println("Failed to send message:", err)
	}
}

func (b *Bot) SendMessageInline(chatID int, text string, markup InlineKeyboardMarkup) {
	apiURL := fmt.Sprintf(apiBaseURL+"/sendMessage", b.Token)
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

	if err := b.postJSON(apiURL, message, "Message sent successfully"); err != nil {
		log.Println("Failed to send message:", err)
	}
}

func (b *Bot) AnswerPreCheckoutQuery(preCheckoutQueryID string, ok bool, errorMessage string) {
	apiURL := fmt.Sprintf(apiBaseURL+"/answerPreCheckoutQuery", b.Token)

	response := AnswerPreCheckoutQueryRequest{
		PreCheckQueryID: preCheckoutQueryID,
		Ok:              ok,
	}

	if err := b.postJSON(apiURL, response, "Answer PreCheckoutQuery successful"); err != nil {
		log.Println("Failed to Answer PreCheckoutQuery:", err)
	}
}
