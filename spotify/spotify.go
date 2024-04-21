package spotifyapi

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
	Track       spotify.FullTrack
	SimpleTrack spotify.SimpleTrack
	Genres      string
	Playlist    spotify.FullPlaylist
}

func SearchSongWithInput(input string) (SpotifyResponse, error) {
	client := getClient()
	track, err := client.searchSpotify(input, spotify.SearchTypeTrack|spotify.SearchTypeArtist|spotify.SearchTypeAlbum)
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

func BuildTrackList(songs []string) []SpotifyResponse {
	client := getClient()

	tracks := []SpotifyResponse{}
	for _, entry := range songs {
		info := strings.Split(entry, "%%")
		artist := info[0]
		song := info[1]

		query := fmt.Sprintf("%s %s", song, artist)
		track, err := client.searchSpotify(query, spotify.SearchTypeTrack|spotify.SearchTypeArtist|spotify.SearchTypeAlbum)
		if err != nil {
			continue
		}

		tracks = append(tracks, SpotifyResponse{SimpleTrack: track.SimpleTrack})
	}

	return tracks
}

func CreatePlaylist(responses []SpotifyResponse) (SpotifyResponse, error) {
	client := getClient()

	user, err := client.Client.CurrentUser()
	if err != nil {
		return SpotifyResponse{},
			fmt.Errorf("error retrieving current user: %v", err)
	}

	playlist, err := client.Client.CreatePlaylistForUser(user.ID, "Temporary Playlist", "New playlist created using the Spotify API", true)
	if err != nil {
		return SpotifyResponse{},
			fmt.Errorf("error creating playlist: %v", err)
	}

	var trackIDs []spotify.ID
	for _, response := range responses {
		trackID := spotify.ID(response.Track.ID)
		trackIDs = append(trackIDs, trackID)
	}

	_, err = client.Client.AddTracksToPlaylist(playlist.ID, trackIDs...)
	if err != nil {
		return SpotifyResponse{},
			fmt.Errorf("error adding tracks to playlist: %v", err)
	}

	return SpotifyResponse{Playlist: *playlist}, err
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

func (spotifyClient SpotifyClient) searchSpotify(query string, searchType spotify.SearchType) (spotify.FullTrack, error) {
	searchResult, err := spotifyClient.Client.Search(query, searchType)
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
