package plugin

import (
	"github.com/Kpovoc/JBot-Go/core/message"
	"github.com/Kpovoc/JBot-Go/plugin/suggesttitle"
)

func GetPluginResponse(pluginName string, msgContent string, m *message.Message) string {
	response := ""

	switch pluginName {
	case "8ball":
		fallthrough
	case "eightball":
		response = EightBall(msgContent)
	case "suggest":
		response = suggesttitle.Plugin(m, msgContent)
	}

	return response
}