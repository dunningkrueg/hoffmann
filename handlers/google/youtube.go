package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type YouTubeSearchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
			} `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
}

func HandleYouTube(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if googleAPIKey == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ YouTube search is not configured. Please ask the bot administrator to set up Google API credentials.",
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
			Content: "🔍 Searching YouTube for: " + query,
		},
	})
	if err != nil {
		return
	}

	results, err := searchYouTube(query)
	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Error searching YouTube: %s", err.Error()))
		return
	}

	if len(results.Items) == 0 {
		editResponse(s, i, fmt.Sprintf("❌ No videos found for: %s", query))
		return
	}

	response := formatYouTubeResults(query, results)
	editResponse(s, i, response)
}

func searchYouTube(query string) (*YouTubeSearchResponse, error) {
	searchURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?key=%s&part=snippet&type=video&maxResults=5&q=%s",
		googleAPIKey,
		url.QueryEscape(query),
	)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("YouTube API returned status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result YouTubeSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func formatYouTubeResults(query string, results *YouTubeSearchResponse) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("🎥 **YouTube Search Results for: %s**\n\n", query))

	for i, item := range results.Items {
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.ID.VideoID)
		description := item.Snippet.Description
		if len(description) > 100 {
			description = description[:97] + "..."
		}

		sb.WriteString(fmt.Sprintf("**%d. [%s](%s)**\n", i+1, item.Snippet.Title, videoURL))
		sb.WriteString(fmt.Sprintf("Channel: %s\n", item.Snippet.ChannelTitle))
		sb.WriteString(fmt.Sprintf("Published: <t:%d:R>\n", item.Snippet.PublishedAt.Unix()))
		sb.WriteString(fmt.Sprintf("Description: %s\n\n", description))
	}

	return sb.String()
} 