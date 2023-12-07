package bot

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	ChatID int
	Text   string
	update Update
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
	update = Update{}
	json.NewDecoder(r.Body).Decode(&update)
	ChatID = update.Message.Chat.ID
	Text = update.Message.Text
	PreCheckoutInvoicePayload := update.PreCheckoutQuery.InvoicePayload

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	b := Bot{
		Token: os.Getenv("TOKEN"),
	}

	b.ProcessMessage(ChatID, Text)

	if PreCheckoutInvoicePayload == "payment-payload" {
		b.AnswerPreCheckoutQuery(update.PreCheckoutQuery.ID, true, "")
	}
	log.Printf("Chat ID = %d\nText = %s", update.Message.Chat.ID, update.Message.Text)
}
