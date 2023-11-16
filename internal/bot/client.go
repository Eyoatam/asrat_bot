package bot

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	ChatID             int
	Text               string
	PreCheckoutQueryID string
	InvoicePayload     string
)

type PreCheckoutQuery struct {
	ID             string `json:"id"`
	TotalAmount    int    `json:"total_amount"`
	InvoicePayload string `json:"invoice_payload"`
}
type Update struct {
	UpdateID         int              `json:"update_id"`
	Message          Messages         `json:"message,omitempty"`
	PreCheckoutQuery PreCheckoutQuery `json:"pre_checkout_query"`
}

type SuccessfulPayment struct {
	InvoicePayload string `json:"invoice_payload"`
}
type Messages struct {
	Message_Id        int               `json:"message_id"`
	Text              string            `json:"text,omitempty"`
	Chat              Chat              `json:"chat,omitempty"`
	SuccessfulPayment SuccessfulPayment `json:"successful_payment "`
}

type Chat struct {
	FirstName string `json:"first_name,omitempty"`
	ID        int    `json:"id"`
	Type      string `json:"type,omitempty"`
	UserName  string `json:"username,omitempty"`
}

func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	var update Update
	json.NewDecoder(r.Body).Decode(&update)
	ChatID = update.Message.Chat.ID
	Text = update.Message.Text
	PreCheckoutQueryID = update.PreCheckoutQuery.ID
	InvoicePayload = update.PreCheckoutQuery.InvoicePayload

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	b := Bot{
		Token: os.Getenv("TOKEN"),
	}

	if InvoicePayload == "pay-load" {
		b.AnswerPreCheckoutQuery(update.PreCheckoutQuery.ID, true, "")
	}
	b.ProcessMessage(ChatID, Text)
	log.Printf("\nChat ID = %d\nText = %s", update.Message.Chat.ID, update.Message.Text)
}
