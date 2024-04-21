package bot

import (
	"errors"
	"strings"
)

const HELP_PROMPT = "!h"
const MUSIC_PROMPT = "!m"
const RECOMMEND_PROMPT = "!r"

func handlePrompt(discord DiscordSession) {
	switch prompt := discord.Prompt; prompt {
	case HELP_PROMPT:
		discord.Session.ChannelMessageSend(discord.Message.ChannelID, "Commands:\n!r <song/artist/album> - returns ten song playlist based on user input.\n")
	case MUSIC_PROMPT:
		err := validateTokens(discord.Message.Content, 3)
		if err != nil {
			discord.Session.ChannelMessageSend(discord.Message.ChannelID, "Invalid music prompt: expected '!m <song-name> <artist>'")
			break
		}
	case RECOMMEND_PROMPT:
		handleRecommendation(discord)

	default:
		discord.Session.ChannelMessageSend(discord.Message.ChannelID, "I'm sorry, I don't understand that command.")
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
