package administrator

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleUnban(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasAdminPermission(s, i) {
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 1 {
		respondWithError(s, i, "Please provide a user ID to unban")
		return
	}

	userID := options[0].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("⏳ Unbanning user with ID %s...", userID),
		},
	})

	err := s.GuildBanDelete(i.GuildID, userID)
	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Failed to unban user: %s", err.Error()))
		return
	}

	response := fmt.Sprintf("✅ Successfully unbanned user with ID **%s**", userID)
	editResponse(s, i, response)
} 