package bot

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyClient struct {
	Client spotify.Client
}

type SpotifyResponse struct {
	Track  spotify.FullTrack
	Genres string
}

func searchSongWithInput(discord DiscordSession) (SpotifyResponse, error) {
	client := getClient()
	track, err := client.searchSpotify(discord.Message.Content)
	if err != nil {
		return SpotifyResponse{}, err
	}

	genres, err := client.getGenres(track)
	if err != nil {
		return SpotifyResponse{}, err
	}

	return SpotifyResponse{
		Track:  track,
		Genres: genres,
	}, nil
}

func getClient() SpotifyClient {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	return SpotifyClient{
		Client: spotify.NewClient(config.Client(context.Background()))}
}

func (spotifyClient SpotifyClient) searchSpotify(query string) (spotify.FullTrack, error) {
	searchResult, err := spotifyClient.Client.Search(
		query,
		spotify.SearchTypeTrack|spotify.SearchTypeArtist|spotify.SearchTypeAlbum,
	)
	if err != nil {
		return spotify.FullTrack{},
			fmt.Errorf("spotify client error: %s", err)
	}
	if len(searchResult.Tracks.Tracks) == 0 {
		return spotify.FullTrack{},
			fmt.Errorf("no tracks found using query: %s", query)
	}

	return searchResult.Tracks.Tracks[0], nil

}

func (spotifyClient SpotifyClient) getGenres(track spotify.FullTrack) (string, error) {
	artist, err := spotifyClient.Client.GetArtist(track.Artists[0].ID)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve artist details: %v", err)
	}

	return strings.Join(artist.Genres, ", "), nil
}
