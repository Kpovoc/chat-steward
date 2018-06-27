package show

import (
	"gitlab.com/Kpovoc/chat-steward/core/message"
	"gitlab.com/Kpovoc/chat-steward/plugin/titles"
)

func StartShow(m *message.Message, msgContent string) string {
	author := m.Sender.DiscordUserName
	if "" == author {
		author = m.Sender.IrcID
	}

	titles.StartShow(msgContent)
	return "Started \"" + msgContent + "\""
}