package administrator

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleUnmute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasAdminPermission(s, i) {
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 1 {
		respondWithError(s, i, "Please mention a user to unmute")
		return
	}

	userID := options[0].UserValue(nil).ID

	targetUser, err := s.User(userID)
	if err != nil {
		respondWithError(s, i, "Failed to fetch user information")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("⏳ Unmuting %s#%s...", targetUser.Username, targetUser.Discriminator),
		},
	})

	err = s.GuildMemberTimeout(i.GuildID, userID, nil)
	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Failed to unmute user: %s", err.Error()))
		return
	}

	response := fmt.Sprintf("✅ Successfully unmuted **%s#%s**", 
		targetUser.Username, 
		targetUser.Discriminator,
	)
	
	editResponse(s, i, response)
} 