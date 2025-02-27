package administrator

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleMute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasAdminPermission(s, i) {
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 2 {
		respondWithError(s, i, "Please mention a user and provide a duration (e.g., 10m, 1h, 7d)")
		return
	}

	userID := options[0].UserValue(nil).ID
	durationStr := options[1].StringValue()

	duration, err := parseDuration(durationStr)
	if err != nil {
		respondWithError(s, i, "Invalid duration format. Please use a number followed by 'm' (minutes), 'h' (hours), or 'd' (days).")
		return
	}

	targetUser, err := s.User(userID)
	if err != nil {
		respondWithError(s, i, "Failed to fetch user information")
		return
	}

	var reason string
	if len(options) > 2 {
		reason = options[2].StringValue()
	} else {
		reason = "No reason provided"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("⏳ Muting %s#%s for %s...", targetUser.Username, targetUser.Discriminator, formatDuration(duration)),
		},
	})

	until := time.Now().Add(duration)
	err = s.GuildMemberTimeout(i.GuildID, userID, &until)
	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Failed to mute user: %s", err.Error()))
		return
	}

	response := fmt.Sprintf("✅ Successfully muted **%s#%s** for **%s**\nReason: %s\nWill be unmuted at: <t:%d:F>", 
		targetUser.Username, 
		targetUser.Discriminator,
		formatDuration(duration),
		reason,
		until.Unix(),
	)
	
	editResponse(s, i, response)
}

func parseDuration(input string) (time.Duration, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	
	if len(input) < 2 {
		return 0, fmt.Errorf("invalid format")
	}
	
	unit := input[len(input)-1:]
	valueStr := input[:len(input)-1]
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}
	
	switch unit {
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown time unit: %s", unit)
	}
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	
	parts := []string{}
	
	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}
	
	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}
	
	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}
	
	if len(parts) == 0 {
		return "0 minutes"
	}
	
	return strings.Join(parts, ", ")
} 