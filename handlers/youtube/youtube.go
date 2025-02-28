package youtube

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

var youtubeService *youtube.Service

type Video struct {
	Title    string
	ID       string
	Channel  string
	Duration string
	Views    string
}

func InitYouTube(apiKey string) error {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return fmt.Errorf("error creating YouTube client: %v", err)
	}
	youtubeService = service
	return nil
}

func HandleYouTube(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if youtubeService == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ YouTube API is not configured. Please ask the bot administrator to set up YouTube API credentials.",
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

	videos, err := searchYouTube(query)
	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Error searching YouTube: %s", err.Error()))
		return
	}

	if len(videos) == 0 {
		editResponse(s, i, fmt.Sprintf("❌ No videos found for: %s", query))
		return
	}

	response := formatResults(query, videos)
	editResponse(s, i, response)
}

func searchYouTube(query string) ([]Video, error) {
	searchCall := youtubeService.Search.List([]string{"snippet"}).
		Q(query).
		Type("video").
		MaxResults(5)

	searchResponse, err := searchCall.Do()
	if err != nil {
		return nil, fmt.Errorf("error performing search: %v", err)
	}

	var videoIDs []string
	for _, item := range searchResponse.Items {
		videoIDs = append(videoIDs, item.Id.VideoId)
	}

	videosCall := youtubeService.Videos.List([]string{"contentDetails", "statistics"}).
		Id(strings.Join(videoIDs, ","))

	videosResponse, err := videosCall.Do()
	if err != nil {
		return nil, fmt.Errorf("error getting video details: %v", err)
	}

	videos := make([]Video, 0, len(searchResponse.Items))
	for _, searchItem := range searchResponse.Items {
		var duration, views string
		for _, videoItem := range videosResponse.Items {
			if videoItem.Id == searchItem.Id.VideoId {
				duration = parseDuration(videoItem.ContentDetails.Duration)
				views = formatViews(videoItem.Statistics.ViewCount)
				break
			}
		}

		videos = append(videos, Video{
			Title:    searchItem.Snippet.Title,
			ID:       searchItem.Id.VideoId,
			Channel:  searchItem.Snippet.ChannelTitle,
			Duration: duration,
			Views:    views,
		})
	}

	return videos, nil
}

func parseDuration(duration string) string {
	duration = strings.TrimPrefix(duration, "PT")
	duration = strings.ToLower(duration)

	var hours, minutes, seconds int

	if strings.Contains(duration, "h") {
		parts := strings.Split(duration, "h")
		fmt.Sscanf(parts[0], "%d", &hours)
		duration = parts[1]
	}

	if strings.Contains(duration, "m") {
		parts := strings.Split(duration, "m")
		fmt.Sscanf(parts[0], "%d", &minutes)
		duration = parts[1]
	}

	if strings.Contains(duration, "s") {
		parts := strings.Split(duration, "s")
		fmt.Sscanf(parts[0], "%d", &seconds)
	}

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func formatViews(viewCount uint64) string {
	if viewCount >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(viewCount)/1000000)
	} else if viewCount >= 1000 {
		return fmt.Sprintf("%.1fK", float64(viewCount)/1000)
	}
	return fmt.Sprintf("%d", viewCount)
}

func formatResults(query string, videos []Video) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("🎥 **YouTube Search Results for: %s**\n\n", query))

	for i, video := range videos {
		sb.WriteString(fmt.Sprintf("**%d. [%s](https://youtu.be/%s)**\n", i+1, video.Title, video.ID))
		sb.WriteString(fmt.Sprintf("Channel: %s\n", video.Channel))
		sb.WriteString(fmt.Sprintf("Duration: %s | Views: %s\n\n", video.Duration, video.Views))
	}

	return sb.String()
}

func editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
