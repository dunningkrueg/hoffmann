package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token            string
	GuildID          string
	SpotifyClientID  string
	SpotifyClientKey string
}

func LoadConfig() *Config {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Token:            os.Getenv("BOT_TOKEN"),
		GuildID:          os.Getenv("GUILD_ID"),
		SpotifyClientID:  os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientKey: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}
}
