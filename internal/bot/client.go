package bot

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	ChatID         int
	Text           string
	InvoicePayload string
	update         Update
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
	SuccessfulPayment SuccessfulPayment `json:"successful_payment"`
}

type Chat struct {
	FirstName string `json:"first_name,omitempty"`
	ID        int    `json:"id"`
	Type      string `json:"type,omitempty"`
	UserName  string `json:"username,omitempty"`
}

func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(r.Body).Decode(&update)
	ChatID = update.Message.Chat.ID
	Text = update.Message.Text
	InvoicePayload = update.PreCheckoutQuery.InvoicePayload

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	b := Bot{
		Token: os.Getenv("TOKEN"),
	}

	if InvoicePayload == "payment-payload" {
		b.AnswerPreCheckoutQuery(update.PreCheckoutQuery.ID, true, "")
	}
	b.ProcessMessage(Text)
	// log.Println("SuccessfulPayment Payload: ", update.Message.SuccessfulPayment.InvoicePayload)
	// log.Printf("\nChat ID = %d\nText = %s", update.Message.Chat.ID, update.Message.Text)
}
