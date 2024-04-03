package bot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const HELP_PROMPT = "!help"
const MUSIC_PROMPT = "!music"
const GOODBYE_PROMPT = "!bye"

type BotConfig struct {
	Token               string
	SpotifyClientId     string
	SpotifyClientSecret string
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func Run(config BotConfig) {
	// Create session
	discord, err := discordgo.New("Bot " + config.Token)
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
	if message.Author.ID == discord.State.User.ID ||
		message.Content[0] != '!' {
		return
	}

	// Respond to user
	switch prompt := strings.Split(message.Content, " ")[0]; prompt {
	case HELP_PROMPT:
		discord.ChannelMessageSend(message.ChannelID, "Sup widdit")
	case MUSIC_PROMPT:
		err := validateTokens(message.Content, 3)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Invalid music prompt: expected '!music <song-name> <artist>'")
			break
		}
		discord.ChannelMessageSend(message.ChannelID, "Cool tunes")
	case GOODBYE_PROMPT:
		discord.ChannelMessageSend(message.ChannelID, "Later, homie")
	default:
		discord.ChannelMessageSend(message.ChannelID, "I'm sorry, I don't understand that command")
	}
}

func validateTokens(message string, expected int) error {
	// Check if prompt tokens match expected count
	words := strings.Fields(message)
	count := len(words)
	if count != expected {
		err := errors.New("token error: expected {{.expected}} but was {{.count}}")
		return err
	}

	return nil
}
