package bot

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	paymentPayload = "payment-payload"
)

type PreCheckoutQuery struct {
	ID             string `json:"id"`
	From           From   `json:"from"`
	TotalAmount    int    `json:"total_amount"`
	InvoicePayload string `json:"invoice_payload"`
}
type From struct {
	ID    int  `json:"id"`
	IsBot bool `json:"is_bot"`
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
	SuccessfulPayment SuccessfulPayment `json:"successful_payment"`
}

type Chat struct {
	FirstName string `json:"first_name,omitempty"`
	ID        int    `json:"id"`
	Type      string `json:"type,omitempty"`
	UserName  string `json:"username,omitempty"`
}

func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	update := Update{}
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text
	preCheckoutInvoicePayload := update.PreCheckoutQuery.InvoicePayload

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}
	b := NewBot(os.Getenv("TOKEN"), os.Getenv("WebHookURL"))
	b.ProcessMessage(update)

	if preCheckoutInvoicePayload == paymentPayload {
		b.AnswerPreCheckoutQuery(update.PreCheckoutQuery.ID, true, "")
	}
	log.Printf("Chat ID = %d\nText = %s", chatID, text)
}
