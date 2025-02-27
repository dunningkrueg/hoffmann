package games

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleCoinFlip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please choose heads or tails!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	userChoice := options[0].StringValue()
	if userChoice != "heads" && userChoice != "tails" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please choose either 'heads' or 'tails'!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🎲 Flipping the coin...",
		},
	})

	rand.Seed(time.Now().UnixNano())
	outcomes := []string{"heads", "tails"}
	result := outcomes[rand.Intn(len(outcomes))]

	time.Sleep(2 * time.Second)

	var response string
	if result == userChoice {
		response = fmt.Sprintf("🎉 The coin landed on **%s**!\nCongratulations, you won! 🏆", result)
	} else {
		response = fmt.Sprintf("💫 The coin landed on **%s**!\nBetter luck next time! 🍀", result)
	}

	editResponse(s, i, response)
}

func editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
