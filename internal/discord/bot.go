package discord

import (
	"log"

	"discord-bot/config"
	"discord-bot/handlers"
	"discord-bot/handlers/administrator"
	"discord-bot/handlers/games"
	googleHandler "discord-bot/handlers/google"
	"discord-bot/handlers/meme"
	spotifyHandler "discord-bot/handlers/spotify"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
	Config  *config.Config
}

func NewBot(cfg *config.Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	err = spotifyHandler.InitSpotify(cfg.SpotifyClientID, cfg.SpotifyClientKey)
	if err != nil {
		log.Printf("Warning: Could not initialize Spotify client: %v", err)
	}

	googleHandler.InitGoogle(cfg.GoogleAPIKey, cfg.GoogleSearchEngineID)
	if cfg.GoogleAPIKey == "" || cfg.GoogleSearchEngineID == "" {
		log.Printf("Warning: Google Search API credentials not configured")
	}

	return &Bot{
		Session: session,
		Config:  cfg,
	}, nil
}

func (b *Bot) Start() error {
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Name {
		case "wiki":
			handlers.HandleWiki(s, i)
		case "clearmsg":
			handlers.HandleClearMsg(s, i)
		case "urban":
			handlers.HandleUrban(s, i)
		case "spotify":
			spotifyHandler.HandleSpotify(s, i)
		case "spotifyuser":
			spotifyHandler.HandleSpotifyUser(s, i)
		case "google":
			googleHandler.HandleGoogle(s, i)
		case "translate":
			googleHandler.HandleTranslate(s, i)
		case "ban":
			administrator.HandleBan(s, i)
		case "unban":
			administrator.HandleUnban(s, i)
		case "kick":
			administrator.HandleKick(s, i)
		case "mute":
			administrator.HandleMute(s, i)
		case "unmute":
			administrator.HandleUnmute(s, i)
		case "youtube":
			googleHandler.HandleYouTube(s, i)
		case "coinflip":
			games.HandleCoinFlip(s, i)
		case "meme":
			meme.HandleMeme(s, i)
		}
	})

	err := b.Session.Open()
	if err != nil {
		return err
	}

	RegisterCommands(b.Session, b.Config.GuildID)
	log.Println("Bot is running...")
	return nil
}
