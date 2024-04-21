package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DiscordSession struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Prompt  string
}

func Run() {
	// Create session
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	LogFatalError(err)

	// Add event handler
	discord.AddHandler(handleMessage)

	// Open Session
	discord.Open()
	defer discord.Close()

	// Run until OS cancellation (e.g. ctrl + C)
	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func handleMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore bot messages and non-commands
	if message.Author.ID == session.State.User.ID ||
		message.Content[0] != '!' {
		return
	}

	handlePrompt(DiscordSession{
		Session: session,
		Message: message,
		Prompt:  strings.Split(message.Content, " ")[0],
	})
}

func LogFatalError(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}
