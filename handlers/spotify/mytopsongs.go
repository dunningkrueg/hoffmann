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
	}

	MyTopSongsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mytopsongs": HandleMyTopSongs,
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
		editResponse(s, i, fmt.Sprintf("ℹ️ User has no public playlists."))
		return
	}

	displayName := "User"
	if len(playlists.Playlists) > 0 {
		if playlists.Playlists[0].Owner.DisplayName != "" {
			displayName = playlists.Playlists[0].Owner.DisplayName
		}
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("📋 **Public Playlists for %s**\n\n", displayName))

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
