package show

import (
	"github.com/Kpovoc/JBot-Go/core/message"
	"github.com/Kpovoc/JBot-Go/plugin/titles"
)

func StartShow(m *message.Message, msgContent string) string {
	author := m.Sender.DiscordUserName
	if "" == author {
		author = m.Sender.IrcID
	}

	titles.StartShow(msgContent)
	return "Started \"" + msgContent + "\""
}