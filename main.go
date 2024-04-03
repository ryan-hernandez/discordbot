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
		log.Fatal("Error loading .env file")
		log.Fatal(err)
	}

	config := bot.BotConfig{
		Token:               os.Getenv("BOT_TOKEN"),
		SpotifyClientId:     os.Getenv("SPOTIFY_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_SECRET"),
	}

	bot.Run(config)
}
