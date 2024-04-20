package spotifyclient

import (
	"context"
	"fmt"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func getSongWithInput(input string) (string, error) {
	clientID := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")
	client := getClient(clientID, clientSecret)

	songLink, err := searchSong(client, input)
	if err != nil {
		fmt.Println("Error searching for song:", err)
		return "", err
	}

	fmt.Println("Spotify link to the first song found:", songLink)
	return songLink, nil
}

func getClient(clientID, clientSecret string) spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}

	httpClient := config.Client(context.Background())
	return spotify.NewClient(httpClient)
}

func searchSong(client spotify.Client, query string) (string, error) {
	result, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		return "", err
	}

	if len(result.Tracks.Tracks) == 0 {
		return "", fmt.Errorf("no tracks found for query: %s", query)
	}

	return result.Tracks.Tracks[0].Endpoint, nil
}
