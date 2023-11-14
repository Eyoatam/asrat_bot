package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  Messages `json:"message,omitempty"`
}

type Messages struct {
	Message_Id int    `json:"message_id"`
	Text       string `json:"text,omitempty"`
	Chat       Chat   `json:"chat,omitempty"`
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
	ChatID := update.Message.Chat.ID
	Text := update.Message.Text

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	b := Bot{
		Token: os.Getenv("TOKEN"),
	}

	b.ProcessMessage(ChatID, Text)
	fmt.Printf("\nChat Id = %d\nText = %s", ChatID, Text)
}
