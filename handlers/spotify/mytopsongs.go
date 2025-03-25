package spotify

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify/v2"
)

var (
	MyTopSongsCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "mytopsongs",
			Description: "Get public playlists for a Spotify username",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "Spotify username (found in Account settings)",
					Required:    true,
				},
			},
		},
		{
			Name:        "howtogetusertoken",
			Description: "Learn how to get your Spotify user token",
		},
	}

	MyTopSongsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mytopsongs":        HandleMyTopSongs,
		"howtogetusertoken": HandleHowToGetUserToken,
	}
)

func HandleMyTopSongs(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				Content: "Please provide your Spotify username (found in your Spotify account settings).",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	username := options[0].StringValue()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔍 Looking up Spotify profile for user: " + username,
		},
	})
	if err != nil {
		return
	}

	playlists, err := spotifyClient.GetPlaylistsForUser(context.Background(), username, spotify.Limit(10))

	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Error finding Spotify user: %s\n\nMake sure you're using your Spotify username (not display name). You can find it in Spotify Account settings.", username))
		return
	}

	if playlists.Total == 0 {
		editResponse(s, i, fmt.Sprintf("ℹ️ User **%s** has no public playlists.", username))
		return
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("📋 **Public Playlists for %s**\n\n", username))

	for i, playlist := range playlists.Playlists {
		trackCount := playlist.Tracks.Total
		response.WriteString(fmt.Sprintf("%d. **%s** - %d tracks\n   🔗 %s\n",
			i+1,
			playlist.Name,
			trackCount,
			playlist.ExternalURLs["spotify"]))
	}

	response.WriteString("\n⚠️ **Note:** Due to Spotify API limitations, we can only show public playlists. Personal top songs require user authentication.")

	editResponse(s, i, response.String())
}

func HandleHowToGetUserToken(s *discordgo.Session, i *discordgo.InteractionCreate) {
	instructions := `
**How to Get Your Spotify User Token**

1. **Visit the Spotify Developer Dashboard:**
   Go to https://developer.spotify.com/dashboard/

2. **Log in** with your Spotify account

3. **Create a new app:**
   - Click "Create An App"
   - Fill in the app name and description (anything will work)
   - Set the redirect URI to: http://localhost:8888/callback
   - Accept the terms and create the app

4. **Get your token:**
   - Go to https://spotify-token-generator.onrender.com/
   - Enter your Client ID and Client Secret from the developer dashboard
   - Click "Generate Token" and select these scopes:
     • user-read-private
     • user-read-email
     • user-top-read

5. **Copy your access token** and use it with the /mytopsongs command

⚠️ **Important:**
- Your token will expire after 1 hour
- Never share your token with anyone else
- This token gives access to your Spotify data
`

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: instructions,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
