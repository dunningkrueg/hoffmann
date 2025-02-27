package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	googleAPIKey         string
	googleSearchEngineID string
)

type GoogleSearchResult struct {
	Items []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		DisplayLink string `json:"displayLink"`
		Snippet     string `json:"snippet"`
	} `json:"items"`
	SearchInformation struct {
		TotalResults        string  `json:"totalResults"`
		SearchTime          float64 `json:"searchTime"`
		FormattedSearchTime string  `json:"formattedSearchTime"`
	} `json:"searchInformation"`
}

func InitGoogle(apiKey, searchEngineID string) {
	googleAPIKey = apiKey
	googleSearchEngineID = searchEngineID
}

func HandleGoogle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if googleAPIKey == "" || googleSearchEngineID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Google Search is not configured. Please ask the bot administrator to set up Google API credentials.",
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
				Content: "Please provide a search query!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	query := options[0].StringValue()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔍 Searching Google for: " + query,
		},
	})
	if err != nil {
		return
	}

	results, err := searchGoogle(query)
	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Error searching Google: %s", err.Error()))
		return
	}

	if len(results.Items) == 0 {
		editResponse(s, i, fmt.Sprintf("❌ No results found for: %s", query))
		return
	}

	response := formatSearchResults(query, results)
	editResponse(s, i, response)
}

func searchGoogle(query string) (*GoogleSearchResult, error) {
	searchURL := fmt.Sprintf(
		"https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s",
		googleAPIKey,
		googleSearchEngineID,
		url.QueryEscape(query),
	)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result GoogleSearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func formatSearchResults(query string, results *GoogleSearchResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("🔍 **Google Search Results for: %s**\n", query))
	sb.WriteString(fmt.Sprintf("Found approximately %s results in %s seconds\n\n",
		results.SearchInformation.TotalResults,
		results.SearchInformation.FormattedSearchTime))

	const maxResults = 5
	resultCount := len(results.Items)
	if resultCount > maxResults {
		resultCount = maxResults
	}

	for i := 0; i < resultCount; i++ {
		item := results.Items[i]
		sb.WriteString(fmt.Sprintf("**%d. [%s](%s)**\n", i+1, item.Title, item.Link))
		sb.WriteString(fmt.Sprintf("%s\n", item.Snippet))
		sb.WriteString(fmt.Sprintf("Source: %s\n\n", item.DisplayLink))
	}

	return sb.String()
}

func editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
