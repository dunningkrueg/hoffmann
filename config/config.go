package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token                string
	GuildID              string
	SpotifyClientID      string
	SpotifyClientKey     string
	GoogleAPIKey         string
	GoogleSearchEngineID string
	YouTubeAPIKey        string
}

func LoadConfig() *Config {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Token:                os.Getenv("BOT_TOKEN"),
		GuildID:              os.Getenv("GUILD_ID"),
		SpotifyClientID:      os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientKey:     os.Getenv("SPOTIFY_CLIENT_SECRET"),
		GoogleAPIKey:         os.Getenv("GOOGLE_API_KEY"),
		GoogleSearchEngineID: os.Getenv("GOOGLE_SEARCH_ENGINE_ID"),
		YouTubeAPIKey:        os.Getenv("YOUTUBE_API_KEY"),
	}
}
