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

var languageCodes = map[string]string{
	"english":    "en",
	"indonesian": "id",
	"japanese":   "ja",
	"korean":     "ko",
	"chinese":    "zh",
	"spanish":    "es",
	"french":     "fr",
	"german":     "de",
	"italian":    "it",
	"russian":    "ru",
	"arabic":     "ar",
	"hindi":      "hi",
	"thai":       "th",
	"vietnamese": "vi",
	"malay":      "ms",
}

type TranslateResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText         string `json:"translatedText"`
			DetectedSourceLanguage string `json:"detectedSourceLanguage,omitempty"`
		} `json:"translations"`
	} `json:"data"`
}

func HandleTranslate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if googleAPIKey == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Google Translate is not configured. Please ask the bot administrator to set up Google API credentials.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 2 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide both target language and text to translate!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	targetLang := strings.ToLower(options[0].StringValue())
	text := options[1].StringValue()

	langCode, exists := languageCodes[targetLang]
	if !exists {
		var supportedLangs []string
		for lang := range languageCodes {
			supportedLangs = append(supportedLangs, lang)
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("❌ Unsupported language. Supported languages are: %s", strings.Join(supportedLangs, ", ")),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔄 Translating...",
		},
	})
	if err != nil {
		return
	}

	translated, detectedLang, err := translateText(text, langCode)
	if err != nil {
		editResponse(s, i, fmt.Sprintf("❌ Error translating text: %s", err.Error()))
		return
	}

	var sourceLangName string
	for lang, code := range languageCodes {
		if code == detectedLang {
			sourceLangName = lang
			break
		}
	}
	if sourceLangName == "" {
		sourceLangName = detectedLang
	}

	response := fmt.Sprintf("🌐 **Translation**\n\n"+
		"**From:** %s\n"+
		"**To:** %s\n\n"+
		"**Original:**\n%s\n\n"+
		"**Translation:**\n%s",
		strings.Title(sourceLangName),
		strings.Title(targetLang),
		text,
		translated)

	editResponse(s, i, response)
}

func translateText(text, targetLang string) (string, string, error) {
	apiURL := fmt.Sprintf(
		"https://translation.googleapis.com/language/translate/v2?key=%s",
		googleAPIKey,
	)

	data := url.Values{}
	data.Set("q", text)
	data.Set("target", targetLang)

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return "", "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("translation API returned status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var result TranslateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Data.Translations) == 0 {
		return "", "", fmt.Errorf("no translation returned")
	}

	translation := result.Data.Translations[0]
	return translation.TranslatedText, translation.DetectedSourceLanguage, nil
}
