package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

const HELP_PROMPT = "!h"
const MUSIC_PROMPT = "!m"
const RECOMMEND_PROMPT = "!r"

type BotConfig struct {
	Token     string
	OpenAIKey string
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
	// Prevent bot from responding to itself and anything that isn't a command
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
			discord.ChannelMessageSend(message.ChannelID, "Invalid music prompt: expected '!m <song-name> <artist>'")
			break
		}
	case RECOMMEND_PROMPT:
		err := handleRecommendation(discord, message)
		if err != nil {
			log.Fatal(err)
		}
	default:
		discord.ChannelMessageSend(message.ChannelID, "I'm sorry, I don't understand that command")
	}
}

func handleRecommendation(discord *discordgo.Session, message *discordgo.MessageCreate) error {
	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		discord.ChannelMessageSend(message.ChannelID, "Error with llm")
		return err
	}

	subject := message.Content
	result, err := getSongWithInput(subject)
	if err != nil {
		discord.ChannelMessageSend(message.ChannelID, "Error generating spotify link")
		return err
	}
	discord.ChannelMessageSend(message.ChannelID, "Creating a playlist based on the following song:\n"+result)

	modelPrompt :=
		"You are to act as the perfect music suggestor. Please, generate an un-numbered list of ten songs by distinct artists based on the following Spotify song link:" + result +
			". You are to listen to the song and analyze it to make a perfect playlist suggestion. Make absolute certain that each song is on Spotify having at least 5,000 plays." +
			" You will provide a link to search for the song on spotify with the following format: https://open.spotify.com/search/artist:{ARTIST}%20track:{TRACK} - {TRACK}, {ARTIST}." +
			" This will be used for the title of the song in the format of [title](link). Only provide the list, do not provide any flavor text."
	response, err := llms.GenerateFromSinglePrompt(ctx, llm, modelPrompt,
		llms.WithTemperature(0.9),
	)
	if err != nil {
		discord.ChannelMessageSend(message.ChannelID, "Error generating response")
		return err
	}

	discord.ChannelMessageSend(message.ChannelID, response)
	return nil
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

func getSongWithInput(input string) (string, error) {
	client := getClient()
	songLink, err := searchSong(client, input)
	if err != nil {
		fmt.Println("Error searching for song:", err)
		return "", err
	}

	return songLink, nil
}

func getClient() spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	httpClient := config.Client(context.Background())
	return spotify.NewClient(httpClient)
}

func searchSong(client spotify.Client, query string) (string, error) {
	result, err := client.Search(
		query,
		spotify.SearchTypeTrack|spotify.SearchTypeArtist|spotify.SearchTypeAlbum,
	)

	if err != nil {
		return "", fmt.Errorf("spotify client error: %s", err)
	}

	if len(result.Tracks.Tracks) == 0 {
		return "", fmt.Errorf("no tracks found using query: %s", query)
	}

	tracks := result.Tracks.Tracks
	mostPopularTrack := result.Tracks.Tracks[0]
	for i := 1; i < len(tracks); i++ {
		highestPopularity := mostPopularTrack.Popularity
		currentTrack := tracks[i]

		if highestPopularity < currentTrack.Popularity {
			mostPopularTrack = currentTrack
		}
	}

	return mostPopularTrack.SimpleTrack.ExternalURLs["spotify"], nil
}
