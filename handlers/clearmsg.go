package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleClearMsg(s *discordgo.Session, i *discordgo.InteractionCreate) {
	member := i.Member
	if member == nil {
		editResponse(s, i, "❌ Error: Cannot verify permissions in DMs")
		return
	}

	hasAdmin := false
	for _, role := range member.Roles {
		roleObj, err := s.State.Role(i.GuildID, role)
		if err != nil {
			continue
		}
		if roleObj.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
			hasAdmin = true
			break
		}
	}

	if !hasAdmin {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Error: You need Administrator permission to use this command",
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
				Content: "Please provide the number of messages to delete",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	amount := int(options[0].IntValue())
	if amount <= 0 || amount > 10000 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Please provide a number between 1 and 10000",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🗑️ Deleting %d messages...", amount),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return
	}

	var allMessageIDs []string
	lastID := ""
	remainingAmount := amount + 1

	for remainingAmount > 0 {
		fetchAmount := 100
		if remainingAmount < 100 {
			fetchAmount = remainingAmount
		}

		messages, err := s.ChannelMessages(i.ChannelID, fetchAmount, lastID, "", "")
		if err != nil {
			editResponse(s, i, "❌ Error: Could not fetch messages. Please try a smaller number.")
			return
		}

		if len(messages) == 0 {
			break
		}

		for _, msg := range messages {
			allMessageIDs = append(allMessageIDs, msg.ID)
			lastID = msg.ID
		}

		remainingAmount -= len(messages)
	}

	if len(allMessageIDs) == 0 {
		editResponse(s, i, "❌ No messages found to delete")
		return
	}

	for start := 0; start < len(allMessageIDs); start += 100 {
		end := start + 100
		if end > len(allMessageIDs) {
			end = len(allMessageIDs)
		}

		batch := allMessageIDs[start:end]
		err = s.ChannelMessagesBulkDelete(i.ChannelID, batch)
		if err != nil {
			editResponse(s, i, "❌ Error: Could not delete some messages. Make sure they are not older than 14 days")
			return
		}
	}

	successMsg := fmt.Sprintf("✅ Successfully deleted %d messages", len(allMessageIDs)-1)
	editResponse(s, i, successMsg)
}
