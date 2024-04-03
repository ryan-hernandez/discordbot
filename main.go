package main

import (
	bot "discord-bot/bot"
	"log"
	"os"

	godotenv "github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")
	bot.BotToken = botToken
	bot.Run()
}
