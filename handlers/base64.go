package handlers

import (
	"encoding/base64"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	Base64Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "encrypt",
			Description: "Encrypt text to base64",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "Text to encrypt",
					Required:    true,
				},
			},
		},
		{
			Name:        "decrypt",
			Description: "Decrypt base64 to text",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "base64",
					Description: "Base64 string to decrypt",
					Required:    true,
				},
			},
		},
	}

	Base64Handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"encrypt": handleEncrypt,
		"decrypt": handleDecrypt,
	}
)

func handleEncrypt(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	text := options[0].StringValue()

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString([]byte(text))

	response := fmt.Sprintf("🔒 Encrypted text:\n```%s```", encoded)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func handleDecrypt(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	encodedText := options[0].StringValue()

	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Error: Invalid base64 string",
			},
		})
		return
	}

	response := fmt.Sprintf("🔓 Decrypted text:\n```%s```", string(decoded))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
