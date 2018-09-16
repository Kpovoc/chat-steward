package user

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofrs/uuid"
)

type User struct{
	ID uuid.UUID  // v4
	JBID string
	DiscordID string
	DiscordUserName string
	IrcID string
	TwitchID string
	TelegramID string
	IsBotAdmin bool

	// Temp until the relationships are more clear
	DiscordUser *discordgo.User
}

func New(jbID string, discordID string, discordUserName string, ircID string,
		 twitchID string, telegramID string, isBotAdmin bool) *User {
		 	id, _ := uuid.NewV4()
		 	return &User{
		 		ID: id,
		 		JBID: jbID,
		 		DiscordID: discordID,
		 		DiscordUserName: discordUserName,
		 		IrcID: ircID,
		 		TwitchID: twitchID,
		 		TelegramID: telegramID,
		 		IsBotAdmin: isBotAdmin,
			}
}
