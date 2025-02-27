package handlers

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type AFKInfo struct {
	Reason    string
	Timestamp time.Time
}

var (
	afkUsers = make(map[string]*AFKInfo)
	afkMutex sync.RWMutex
)

func HandleAFK(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var reason string
	if len(options) > 0 {
		reason = options[0].StringValue()
	} else {
		reason = "No reason provided"
	}

	userID := i.Member.User.ID
	afkMutex.Lock()
	afkUsers[userID] = &AFKInfo{
		Reason:    reason,
		Timestamp: time.Now(),
	}
	afkMutex.Unlock()

	response := fmt.Sprintf("✅ You are now AFK: %s", reason)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	afkMutex.RLock()
	defer afkMutex.RUnlock()

	for _, mention := range m.Mentions {
		if afkInfo, isAFK := afkUsers[mention.ID]; isAFK {
			timeSince := time.Since(afkInfo.Timestamp)
			timeStr := formatDuration(timeSince)

			response := fmt.Sprintf("⚠️ %s is AFK (%s ago)\nReason: %s",
				mention.Username,
				timeStr,
				afkInfo.Reason,
			)

			s.ChannelMessageSend(m.ChannelID, response)
		}
	}

	if afkInfo, isAFK := afkUsers[m.Author.ID]; isAFK {
		afkMutex.RUnlock()
		afkMutex.Lock()
		delete(afkUsers, m.Author.ID)
		afkMutex.Unlock()
		afkMutex.RLock()

		timeSince := time.Since(afkInfo.Timestamp)
		timeStr := formatDuration(timeSince)

		response := fmt.Sprintf("👋 Welcome back %s! You were AFK for %s",
			m.Author.Username,
			timeStr,
		)

		s.ChannelMessageSend(m.ChannelID, response)
	}
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		if hours == 1 {
			if minutes == 0 {
				return "1 hour"
			}
			return fmt.Sprintf("1 hour %d minutes", minutes)
		}
		if minutes == 0 {
			return fmt.Sprintf("%d hours", hours)
		}
		return fmt.Sprintf("%d hours %d minutes", hours, minutes)
	}

	if minutes == 0 {
		return "less than a minute"
	}
	if minutes == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%d minutes", minutes)
}
