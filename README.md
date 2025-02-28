# Hoffmann Discord Bot

> ⚠️ **Note: This project is currently under active development**

A powerful Discord bot built in Go that combines various features including Spotify integration, Wikipedia search, Urban Dictionary lookup, message management capabilities, and fun games.

## Features

- **Spotify Integration**
  - `/spotify` - Search for song information on Spotify
  - `/spotifyuser` - Get user profile and top tracks (requires authentication - in development)

- **Information Commands**
  - `/wiki` - Search Wikipedia articles
  - `/urban` - Look up terms on Urban Dictionary
  - `/google` - Search Google for information
  - `/translate` - Translate text to different languages
  - `/youtube` - Search for videos on YouTube

- **Games & Fun**
  - `/coinflip` - Play a game of heads or tails with the bot
    - Choose 'heads' or 'tails'
    - Watch the coin flip animation
    - Win or lose with fun emoji responses!
  - `/meme` - Get a random meme from r/memes
    - Fresh memes from Reddit's r/memes
    - Includes title, image, and upvote count
    - Click to view original post

- **Moderation Tools**
  - `/clearmsg` - Bulk delete messages (Admin only)
  - `/mute` - Timeout a user for a specified duration (Admin only)
  - `/unmute` - Remove timeout from a user (Admin only)
  - `/kick` - Kick a user from the server (Admin only)
  - `/ban` - Ban a user from the server (Admin only)
  - `/unban` - Unban a user from the server (Admin only)

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Discord Bot Token
- Spotify Developer Credentials
- Google API Key and Custom Search Engine ID

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

# Google API Credentials
GOOGLE_API_KEY=your_google_api_key_here           # Get from Google Cloud Console
GOOGLE_SEARCH_ENGINE_ID=your_custom_search_engine_id_here  # Create at cse.google.com
```

### Installation
1. Clone the repository
2. Install dependencies:
```bash
go mod tidy
```
3. Run the bot:
```bash
go run cmd/bot/main.go
```

## Command Usage

### Translation
The `/translate` command supports multiple languages:
- english
- indonesian
- japanese
- korean
- chinese
- spanish
- french
- german
- italian
- russian
- arabic
- hindi
- thai
- vietnamese
- malay

Example usage:
```
/translate english Selamat pagi
/translate japanese Hello world
/translate spanish How are you?
```

### Games
#### Coinflip
Play a game of heads or tails with the bot:
```
/coinflip heads  - Bet on heads
/coinflip tails  - Bet on tails
```
The bot will:
1. Flip a virtual coin
2. Show a brief animation
3. Reveal if you won or lost with fun emoji responses

## Development Status

This bot is currently in active development. Some features may be incomplete or subject to change. Current development focus:
- Enhancing Spotify integration
- Adding user authentication for Spotify features
- Implementing additional games and commands
- Improving error handling and user feedback

## Dependencies

- [discordgo](https://github.com/bwmarrin/discordgo) - Discord API wrapper
- [spotify-go](https://github.com/zmb3/spotify) - Spotify API client
- [go-wiki](https://github.com/trietmn/go-wiki) - Wikipedia API wrapper
- [godotenv](https://github.com/joho/godotenv) - Environment configuration
- [go-googletrans](https://github.com/Conight/go-googletrans) - Google Translate API wrapper
- [google-api-go-client](https://github.com/googleapis/google-api-go-client) - Official Google APIs Client Library

## Contributing

This project is under development and contributions are welcome. Feel free to:
- Report bugs
- Suggest new features
- Submit pull requests

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

