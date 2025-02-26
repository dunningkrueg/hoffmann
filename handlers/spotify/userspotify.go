package spotify

import (
	"github.com/bwmarrin/discordgo"
)

func HandleSpotifyUser(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				Content: "Please provide a Spotify username!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	username := options[0].StringValue()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔍 Fetching Spotify profile for: " + username,
		},
	})
	if err != nil {
		return
	}

	editResponse(s, i, "❌ Sorry, this feature requires Spotify user authentication. Please use /spotify command instead.")
	return
}
