package meme

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

type MemeResponse struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Postlink string `json:"postLink"`
	Ups      int    `json:"ups"`
}

func HandleMeme(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔍 Fetching a fresh meme...",
		},
	})
	if err != nil {
		return
	}

	meme, err := getRandomMeme()
	if err != nil {
		editResponse(s, i, "❌ Failed to fetch meme: "+err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: meme.Title,
		Image: &discordgo.MessageEmbedImage{
			URL: meme.URL,
		},
		Color: 0xff4500,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("👍 %d • Click title to view on Reddit", meme.Ups),
		},
		URL: meme.Postlink,
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}

func getRandomMeme() (*MemeResponse, error) {
	resp, err := client.Get("https://meme-api.com/gimme/memes")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meme: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var meme MemeResponse
	if err := json.Unmarshal(body, &meme); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &meme, nil
}

func editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
