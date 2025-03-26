package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	PingCommand = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Check your ping to Discord",
		},
	}

	PingHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": HandlePing,
	}
)

// HandlePing checks the ping time between a user's command and the bot receiving it
func HandlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Calculate the ping time based on the interaction creation vs current time
	// InteractionCreate.ID is a Snowflake which contains a timestamp
	snowflakeTime := SnowflakeTimestamp(i.ID)
	ping := time.Since(snowflakeTime).Milliseconds()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🏓 **Pong!**\nYour ping: **%dms**", ping),
		},
	})
}

// SnowflakeTimestamp extracts the timestamp from a Discord Snowflake ID
// Discord Snowflake format: https://discord.com/developers/docs/reference#snowflakes
func SnowflakeTimestamp(snowflake string) time.Time {

	var snowflakeID int64
	fmt.Sscanf(snowflake, "%d", &snowflakeID)

	discordEpoch := int64(1420070400000)

	timestamp := (snowflakeID >> 22) + discordEpoch

	return time.Unix(0, timestamp*1000000)
}
