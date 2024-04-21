package bot

import (
	"context"
	spotifyapi "discord-bot/spotify"
	"fmt"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func (discord DiscordSession) CheckError(channelID string, err error) {
	if err != nil {
		errorMessage := fmt.Sprintf("Error: %s", err.Error())
		discord.Session.ChannelMessageSend(channelID, errorMessage)
		err = nil
	}
}

func handleRecommendation(discord DiscordSession) {
	ctx := context.Background()
	llm, err := openai.New()
	discord.CheckError(discord.Message.ChannelID, err)

	// Get Song
	spotifyResponse, err := spotifyapi.SearchSongWithInput(discord.Message.Content)
	discord.CheckError(discord.Message.ID, err)

	songPreview := spotifyResponse.Track.SimpleTrack.ExternalURLs["spotify"]
	discord.Session.ChannelMessageSend(discord.Message.ChannelID, "Creating a playlist based on the following song:\n"+songPreview)

	modelPrompt :=
		"You are to act as the perfect music suggestor. Please, generate an un-numbered list of ten songs by distinct artists based on the following Spotify song link:" + songPreview +
			". You are to listen to the song and analyze it to make a perfect playlist suggestion. Make absolute certain that each song is on Spotify having at least 5,000 plays" +
			". Make sure that the songs in the list were originally released within the five years of the following date: " + spotifyResponse.Track.Album.ReleaseDate +
			". Make absolute certain that the songs you recommend are from a very similar genre to the following: " + spotifyResponse.Genres +
			". Under absolutely no circumstances are you to recommend an Ed Sheeran song. Do not recommend any Ed Sheeran songs, whatsoever." +
			". You will provide a response consisting of a single string of artists and tracks in the following format: {ARTIST}%%{TRACK},{ARTIST}%%{TRACK},{ARTIST}%%{TRACK}.. etc" +
			". Make absolute certain that there are no line breaks in this string." +
			". Only provide this list and no other flavor text."
	modelResponse, err := llms.GenerateFromSinglePrompt(
		ctx, llm, modelPrompt,
		llms.WithModel(os.Getenv("GPT_MODEL")),
		llms.WithTemperature(0.9),
	)

	songs := strings.Split(modelResponse, ",")
	discord.CheckError(discord.Message.ChannelID, err)

	responses := spotifyapi.BuildTrackList(songs)
	playlist, err := spotifyapi.CreatePlaylist(responses)
	discord.CheckError(discord.Message.ChannelID, err)

	discord.Session.ChannelMessageSend(discord.Message.ChannelID, playlist.Playlist.ExternalURLs["spotify"])
}
