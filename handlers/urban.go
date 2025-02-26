package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type UrbanResponse struct {
	List []struct {
		Definition string   `json:"definition"`
		Permalink  string   `json:"permalink"`
		ThumbsUp   int      `json:"thumbs_up"`
		ThumbsDown int      `json:"thumbs_down"`
		Author     string   `json:"author"`
		Word       string   `json:"word"`
		Example    string   `json:"example"`
		Tags       []string `json:"tags"`
	} `json:"list"`
}

func HandleUrban(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a word to search!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	query := options[0].StringValue()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔍 Searching Urban Dictionary for: " + query,
		},
	})
	if err != nil {
		return
	}

	apiURL := fmt.Sprintf("https://api.urbandictionary.com/v0/define?term=%s", url.QueryEscape(query))
	resp, err := http.Get(apiURL)
	if err != nil {
		editResponse(s, i, "❌ Error: Could not connect to Urban Dictionary")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		editResponse(s, i, "❌ Error: Could not read response from Urban Dictionary")
		return
	}

	var urbanResp UrbanResponse
	err = json.Unmarshal(body, &urbanResp)
	if err != nil {
		editResponse(s, i, "❌ Error: Could not parse Urban Dictionary response")
		return
	}

	if len(urbanResp.List) == 0 {
		editResponse(s, i, fmt.Sprintf("❌ No definitions found for: %s", query))
		return
	}

	def := urbanResp.List[0]

	definition := strings.ReplaceAll(def.Definition, "[", "")
	definition = strings.ReplaceAll(definition, "]", "")
	example := strings.ReplaceAll(def.Example, "[", "")
	example = strings.ReplaceAll(example, "]", "")

	response := fmt.Sprintf("📚 **%s**\n\n**Definition:**\n%s\n\n**Example:**\n%s\n\n👍 %d | 👎 %d\n\n🔗 %s",
		def.Word,
		truncateText(definition, 1000),
		truncateText(example, 500),
		def.ThumbsUp,
		def.ThumbsDown,
		def.Permalink)

	editResponse(s, i, response)
}
