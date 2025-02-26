package discord

import (
	"log"

	"discord-bot/config"
	"discord-bot/handlers"
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
