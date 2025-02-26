package spotify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

var spotifyClient *spotify.Client

func InitSpotify(clientID, clientSecret string) error {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(context.Background(), token)
	spotifyClient = spotify.New(httpClient)
	return nil
}

func formatDuration(ms int) string {
	duration := time.Duration(ms) * time.Millisecond
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func HandleSpotify(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if spotifyClient == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Spotify client not initialized",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a song to search!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	query := options[0].StringValue()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔍 Searching Spotify for: " + query,
		},
	})
	if err != nil {
		return
	}

	results, err := spotifyClient.Search(context.Background(), query, spotify.SearchTypeTrack)
	if err != nil {
		editResponse(s, i, "❌ Error searching Spotify")
		return
	}

	if results.Tracks == nil || len(results.Tracks.Tracks) == 0 {
		editResponse(s, i, fmt.Sprintf("❌ No songs found for: %s", query))
		return
	}

	track := results.Tracks.Tracks[0]

	fullTrack, err := spotifyClient.GetTrack(context.Background(), track.ID)
	if err != nil {
		editResponse(s, i, "❌ Error getting track details")
		return
	}

	var artists []string
	for _, artist := range fullTrack.Artists {
		artists = append(artists, artist.Name)
	}

	album, err := spotifyClient.GetAlbum(context.Background(), fullTrack.Album.ID)
	if err != nil {
		editResponse(s, i, "❌ Error getting album details")
		return
	}

	response := fmt.Sprintf("🎵 **%s**\n\n"+
		"👤 **Artists:** %s\n"+
		"💿 **Album:** %s\n"+
		"⏱️ **Duration:** %s\n"+
		"📅 **Release Date:** %s\n"+
		"🎼 **Key Features:**\n"+
		"- Popularity: %d/100\n"+
		"- Track Number: %d\n"+
		"- Disc Number: %d\n\n"+
		"🔗 **Listen on Spotify:** %s",
		fullTrack.Name,
		strings.Join(artists, ", "),
		fullTrack.Album.Name,
		formatDuration(fullTrack.Duration),
		album.ReleaseDate,
		fullTrack.Popularity,
		fullTrack.TrackNumber,
		fullTrack.DiscNumber,
		fullTrack.ExternalURLs["spotify"])

	editResponse(s, i, response)
}

func editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
