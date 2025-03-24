package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var (
	instagramPattern   = regexp.MustCompile(`(?i)https?://(www\.)?(instagram\.com|instagr\.am)/(?:p|reel)/([a-zA-Z0-9_-]+)/?`)
	twitterPattern     = regexp.MustCompile(`(?i)https?://(www\.)?(twitter\.com|x\.com)/([a-zA-Z0-9_]+)/status/([0-9]+)/?`)
	tiktokPattern      = regexp.MustCompile(`(?i)https?://(www\.)?(tiktok\.com)/@([a-zA-Z0-9_\.]+)/video/([0-9]+)/?`)
	twitterBearerToken string
)

func InitTwitterCredentials(bearerToken string) {
	twitterBearerToken = bearerToken
}

func HandleAutoEmbed(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
		return
	}

	if instagramPattern.MatchString(m.Content) {
		embedInstagram(s, m)
	} else if twitterPattern.MatchString(m.Content) {
		embedTwitter(s, m)
	} else if tiktokPattern.MatchString(m.Content) {
		embedTiktok(s, m)
	}
}

func embedInstagram(s *discordgo.Session, m *discordgo.MessageCreate) {
	matches := instagramPattern.FindStringSubmatch(m.Content)
	if len(matches) < 4 {
		return
	}

	postID := matches[3]

	apiURL := fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1", postID)

	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return
	}

	embed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: "Instagram Post",
		URL:   fmt.Sprintf("https://www.instagram.com/p/%s/", postID),
		Color: 0xE1306C,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Instagram",
			IconURL: "https://www.instagram.com/static/images/ico/favicon-192.png/68d99ba29cc8.png",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func embedTwitter(s *discordgo.Session, m *discordgo.MessageCreate) {
	matches := twitterPattern.FindStringSubmatch(m.Content)
	if len(matches) < 5 {
		return
	}

	username := matches[3]
	tweetID := matches[4]

	apiURL := fmt.Sprintf("https://api.twitter.com/2/tweets/%s?expansions=author_id,attachments.media_keys&user.fields=name,profile_image_url&media.fields=url,preview_image_url,type", tweetID)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Bearer "+twitterBearerToken)

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return
	}

	embed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: fmt.Sprintf("Tweet by @%s", username),
		URL:   fmt.Sprintf("https://twitter.com/%s/status/%s", username, tweetID),
		Color: 0x1DA1F2,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Twitter/X",
			IconURL: "https://abs.twimg.com/responsive-web/web/icon-default.604e2486a34a2f6e.png",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func embedTiktok(s *discordgo.Session, m *discordgo.MessageCreate) {
	matches := tiktokPattern.FindStringSubmatch(m.Content)
	if len(matches) < 5 {
		return
	}

	username := matches[3]
	videoID := matches[4]

	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       fmt.Sprintf("TikTok by @%s", username),
		URL:         fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", username, videoID),
		Color:       0x000000,
		Description: "Click to watch TikTok video",
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "TikTok",
			IconURL: "https://sf16-scmcdn-sg.ibytedtos.com/goofy/tiktok/web/node/_next/static/images/logo-black-3d833cb1aca8c5b4033f0b316c1b0ac4.svg",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func InitAutoEmbed(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		HandleAutoEmbed(s, m)
	})
}
