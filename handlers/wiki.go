package handlers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	gowiki "github.com/trietmn/go-wiki"
)

func HandleWiki(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a search query!",
			},
		})
		return
	}

	query := options[0].StringValue()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔍 Searching Wikipedia for: " + query,
		},
	})
	if err != nil {
		return
	}

	results, suggestions, err := gowiki.Search(query, 1, false)
	if err != nil || len(results) == 0 {
		if len(suggestions) > 0 {
			editResponse(s, i, fmt.Sprintf("❌ No results found for: %s\nDid you mean: %s?", query, suggestions[0]))
		} else {
			editResponse(s, i, "❌ Error: Could not find any results for: "+query)
		}
		return
	}

	title := results[0]
	safeTitle := url.QueryEscape(title)
	wikiURL := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", safeTitle)

	summary, err := gowiki.Summary(title, 1000, -1, false, false)
	if err != nil {
		editResponse(s, i, "❌ Error: Could not get summary for: "+title)
		return
	}

	response := fmt.Sprintf("📚 **%s**\n\n%s\n\n🔗 Read more: %s",
		title,
		truncateText(summary, 1500),
		wikiURL)

	editResponse(s, i, response)
}

func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	truncated := text[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace != -1 {
		truncated = truncated[:lastSpace]
	}
	return truncated + "..."
}

func editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
