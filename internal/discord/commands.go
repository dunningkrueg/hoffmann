package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "wiki",
		Description: "Search Wikipedia for information",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "What do you want to search for?",
				Required:    true,
			},
		},
	},
	{
		Name:                     "clearmsg",
		Description:              "Delete multiple messages (Admin only)",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Number of messages to delete (max 10000)",
				Required:    true,
				MinValue:    &[]float64{1}[0],
				MaxValue:    10000,
			},
		},
	},
	{
		Name:        "urban",
		Description: "Search Urban Dictionary for a word or phrase",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "term",
				Description: "Word or phrase to look up",
				Required:    true,
			},
		},
	},
	{
		Name:        "spotify",
		Description: "Search for song information on Spotify",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "song",
				Description: "Name of the song to search",
				Required:    true,
			},
		},
	},
	{
		Name:        "spotifyuser",
		Description: "Get Spotify user profile and top tracks",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "Spotify username to look up",
				Required:    true,
			},
		},
	},
	{
		Name:        "google",
		Description: "Search Google for information",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "What do you want to search for?",
				Required:    true,
			},
		},
	},
	{
		Name:        "translate",
		Description: "Translate text to another language",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "language",
				Description: "Target language",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "English", Value: "english"},
					{Name: "Indonesian", Value: "indonesian"},
					{Name: "Japanese", Value: "japanese"},
					{Name: "Korean", Value: "korean"},
					{Name: "Chinese", Value: "chinese"},
					{Name: "Spanish", Value: "spanish"},
					{Name: "French", Value: "french"},
					{Name: "German", Value: "german"},
					{Name: "Italian", Value: "italian"},
					{Name: "Russian", Value: "russian"},
					{Name: "Arabic", Value: "arabic"},
					{Name: "Hindi", Value: "hindi"},
					{Name: "Thai", Value: "thai"},
					{Name: "Vietnamese", Value: "vietnamese"},
					{Name: "Malay", Value: "malay"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "Text to translate",
				Required:    true,
			},
		},
	},
	{
		Name:                     "ban",
		Description:              "Ban a user from the server (Admin only)",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User to ban",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "Reason for the ban",
				Required:    false,
			},
		},
	},
	{
		Name:                     "unban",
		Description:              "Unban a user from the server (Admin only)",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "userid",
				Description: "User ID to unban",
				Required:    true,
			},
		},
	},
	{
		Name:                     "kick",
		Description:              "Kick a user from the server (Admin only)",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User to kick",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "Reason for the kick",
				Required:    false,
			},
		},
	},
	{
		Name:                     "mute",
		Description:              "Timeout a user for a specified duration (Admin only)",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User to mute",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "duration",
				Description: "Duration (e.g., 10m, 1h, 7d)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "Reason for the mute",
				Required:    false,
			},
		},
	},
	{
		Name:                     "unmute",
		Description:              "Remove timeout from a user (Admin only)",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User to unmute",
				Required:    true,
			},
		},
	},
	{
		Name:        "youtube",
		Description: "Search for videos on YouTube",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "What do you want to search for?",
				Required:    true,
			},
		},
	},
	{
		Name:        "coinflip",
		Description: "Play a game of heads or tails",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "choice",
				Description: "Choose heads or tails",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Heads",
						Value: "heads",
					},
					{
						Name:  "Tails",
						Value: "tails",
					},
				},
			},
		},
	},
	{
		Name:        "meme",
		Description: "Get a random meme from r/memes",
	},
	{
		Name:        "afk",
		Description: "Set your status as AFK",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "Reason for being AFK",
				Required:    false,
			},
		},
	},
	{
		Name:        "encrypt",
		Description: "Encrypt text to base64",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "Text to encrypt",
				Required:    true,
			},
		},
	},
	{
		Name:        "decrypt",
		Description: "Decrypt base64 to text",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "base64",
				Description: "Base64 string to decrypt",
				Required:    true,
			},
		},
	},
	{
		Name:        "mytopsongs",
		Description: "Get public playlists for a Spotify username",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "Spotify username (found in Account settings)",
				Required:    true,
			},
		},
	},
}

func RegisterCommands(s *discordgo.Session, guildID string) {
	if s.State.User == nil {
		log.Fatal("Bot user not found in session")
	}

	for _, cmd := range Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
		if err != nil {
			log.Printf("Error creating command %v: %v", cmd.Name, err)
		}
	}
}
