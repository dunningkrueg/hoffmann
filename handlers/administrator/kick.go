package administrator

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleKick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasAdminPermission(s, i) {
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 1 {
		respondWithError(s, i, "Please mention a user to kick")
		return
	}

	userID := options[0].UserValue(nil).ID

	var reason string
	if len(options) > 1 {
		reason = options[1].StringValue()
	} else {
		reason = "No reason provided"
	}

	targetUser, err := s.User(userID)
	if err != nil {
		respondWithError(s, i, "Failed to fetch user information")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("⏳ Kicking %s#%s...", targetUser.Username, targetUser.Discriminator),
		},
	})

	err = s.GuildMemberDeleteWithReason(i.GuildID, userID, reason)
	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Failed to kick user: %s", err.Error()))
		return
	}

	response := fmt.Sprintf("✅ Successfully kicked **%s#%s**\nReason: %s",
		targetUser.Username,
		targetUser.Discriminator,
		reason,
	)

	editResponse(s, i, response)
}
