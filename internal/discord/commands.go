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
