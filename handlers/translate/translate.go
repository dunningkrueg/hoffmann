package translate

import (
	"fmt"
	"strings"

	googletrans "github.com/Conight/go-googletrans"
	"github.com/bwmarrin/discordgo"
)

var trans = googletrans.New()

var languageCodes = map[string]string{
	"english":    "en",
	"indonesian": "id",
	"japanese":   "ja",
	"korean":     "ko",
	"chinese":    "zh-cn",
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

func HandleTranslate(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	result, err := trans.Translate(text, "auto", langCode)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &[]string{fmt.Sprintf("❌ Error translating text: %s", err.Error())}[0],
		})
		return
	}

	var sourceLangName string
	for lang, code := range languageCodes {
		if code == result.Src {
			sourceLangName = lang
			break
		}
	}
	if sourceLangName == "" {
		sourceLangName = result.Src
	}

	response := fmt.Sprintf("🌐 **Translation**\n\n"+
		"**From:** %s\n"+
		"**To:** %s\n\n"+
		"**Original:**\n%s\n\n"+
		"**Translation:**\n%s",
		strings.Title(sourceLangName),
		strings.Title(targetLang),
		text,
		result.Text)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &response,
	})
}
