package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eyoatam/asrat_bot/internal/bot"
	"github.com/joho/godotenv"
)

func main() {
	http.HandleFunc("/telegram-webhook", bot.WebHookHandler)

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	b := bot.NewBot(os.Getenv("TOKEN"), os.Getenv("WEBHOOKURL"))
	b.Connect()

	port := "4000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	s := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server started on ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
