package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func Run() {
	// Create session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// Add event handler
	discord.AddHandler(newMessage)

	// Open Session
	discord.Open()
	defer discord.Close()

	// Keep bot running until OS cancellation (e.g. ctrl + C)
	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Prevent bot from responding to itself
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Response to user
	switch {
	case strings.Contains(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Sup widdit")
	case strings.Contains(message.Content, "!bye"):
		discord.ChannelMessageSend(message.ChannelID, "Later, homie")
	}
}
