// hookurl = https://0b99-196-188-126-6.ngrok-free.app/telegram-message

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/eyoatam/asrat_bot/internal/bot"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/telegram-webhook", bot.WebHookHandler)
	// print(bot.ChatID)

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	b := bot.Bot{
		Token:      os.Getenv("TOKEN"),
		WebHookUrl: os.Getenv("WEBHOOKURL"),
	}

	b.Connect()
	port := 4000
	if os.Getenv("PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}

	s := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("127.0.0.1:%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server started on ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
