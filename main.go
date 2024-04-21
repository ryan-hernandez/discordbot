package main

import (
	"discord-bot/bot"
	"log"

	godotenv "github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		log.Fatal(err)
	}

	bot.Run()
}
