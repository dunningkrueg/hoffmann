package administrator

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleBan(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasAdminPermission(s, i) {
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 1 {
		respondWithError(s, i, "Please mention a user to ban")
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
			Content: fmt.Sprintf("⏳ Banning user %s#%s...", targetUser.Username, targetUser.Discriminator),
		},
	})

	err = s.GuildBanCreateWithReason(i.GuildID, userID, reason, 0)
	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Failed to ban user: %s", err.Error()))
		return
	}

	response := fmt.Sprintf("✅ Successfully banned **%s#%s**\nReason: %s",
		targetUser.Username,
		targetUser.Discriminator,
		reason,
	)

	editResponse(s, i, response)
}

func hasAdminPermission(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	member := i.Member
	if member == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Error: Cannot verify permissions in DMs",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return false
	}

	for _, role := range member.Roles {
		roleObj, err := s.State.Role(i.GuildID, role)
		if err != nil {
			continue
		}
		if roleObj.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
			return true
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ Error: You need Administrator permission to use this command",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	return false
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, errMsg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ " + errMsg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
