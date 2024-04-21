package bot

import (
	"context"
	"fmt"
	"os"

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

	spotifyResponse, err := searchSongWithInput(discord)
	discord.CheckError(discord.Message.ID, err)

	songPreview := spotifyResponse.Track.SimpleTrack.ExternalURLs["spotify"]
	discord.Session.ChannelMessageSend(discord.Message.ChannelID, "Creating a playlist based on the following song:\n"+songPreview)

	modelPrompt :=
		"You are to act as the perfect music suggestor. Please, generate an un-numbered list of ten songs by distinct artists based on the following Spotify song link:" + songPreview +
			". You are to listen to the song and analyze it to make a perfect playlist suggestion. Make absolute certain that each song is on Spotify having at least 5,000 plays" +
			". Make sure that the songs in the list were originally released within the five years of the following date: " + spotifyResponse.Track.Album.ReleaseDate +
			". Make absolute certain that the songs you recommend are from a very similar genre to the following: " + spotifyResponse.Genres +
			". Under absolutely no circumstances are you to recommend an Ed Sheeran song. Do not recommend any Ed Sheeran songs, whatsoever." +
			". You will provide a link to search for the song on spotify with the following format: https://open.spotify.com/search/artist:{ARTIST}%20track:{TRACK} - {TRACK}, {ARTIST}" +
			". This will be used for the title of the song in the format of [title](link). Only provide the list, do not provide any flavor text."
	modelResponse, err := llms.GenerateFromSinglePrompt(
		ctx, llm, modelPrompt,
		llms.WithModel(os.Getenv("GPT_MODEL")),
		llms.WithTemperature(0.9),
	)
	discord.CheckError(discord.Message.ChannelID, err)
	discord.Session.ChannelMessageSend(discord.Message.ChannelID, modelResponse)
}
