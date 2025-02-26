# Hoffmann Discord Bot

> ⚠️ **Note: This project is currently under active development**

A powerful Discord bot built in Go that combines various features including Spotify integration, Wikipedia search, Urban Dictionary lookup, and message management capabilities.

## Features

- **Spotify Integration**
  - `/spotify` - Search for song information on Spotify
  - `/spotifyuser` - Get user profile and top tracks (requires authentication - in development)

- **Information Commands**
  - `/wiki` - Search Wikipedia articles
  - `/urban` - Look up terms on Urban Dictionary

- **Moderation Tools**
  - `/clearmsg` - Bulk delete messages (Admin only)

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Discord Bot Token
- Spotify Developer Credentials

### Environment Setup
1. Copy `.env.example` to `config/.env`:
```bash
cp .env.example config/.env
```

2. Update the environment variables in `config/.env`:
```env
# Discord Configuration
BOT_TOKEN=your_discord_bot_token_here        # Get from Discord Developer Portal
GUILD_ID=your_discord_server_id_here         # Your Discord Server ID

# Spotify API Credentials
SPOTIFY_CLIENT_ID=your_spotify_client_id_here     # Get from Spotify Developer Dashboard
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret_here
```

### Installation
1. Clone the repository
2. Install dependencies:
```bash
go mod download
```
3. Run the bot:
```bash
go run cmd/bot/main.go
```

## Development Status

This bot is currently in active development. Some features may be incomplete or subject to change. Current development focus:
- Enhancing Spotify integration
- Adding user authentication for Spotify features
- Implementing additional commands
- Improving error handling and user feedback

## Dependencies

- [discordgo](https://github.com/bwmarrin/discordgo) - Discord API wrapper
- [spotify-go](https://github.com/zmb3/spotify) - Spotify API client
- [go-wiki](https://github.com/trietmn/go-wiki) - Wikipedia API wrapper
- [godotenv](https://github.com/joho/godotenv) - Environment configuration

## Contributing

This project is under development and contributions are welcome. Feel free to:
- Report bugs
- Suggest new features
- Submit pull requests

## License

This project is licensed under the MIT License - see the LICENSE file for details.

