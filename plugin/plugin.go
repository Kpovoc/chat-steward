package plugin

import (
	"gitlab.com/Kpovoc/chat-steward/core/message"
	"gitlab.com/Kpovoc/chat-steward/plugin/titles"
	"gitlab.com/Kpovoc/chat-steward/plugin/show"
)

func GetPluginResponse(pluginName string, msgContent string, m *message.Message) string {
	response := ""

	switch pluginName {
	case "8ball":
		fallthrough
	case "eightball":
		response = EightBall(msgContent)
	case "suggest":
		response = titles.Plugin(m, msgContent)
	case "start_show":
		response = show.StartShow(m, msgContent)
	}

	return response
}

func InitPlugins(resourceDir string) {
	titles.WebInit(resourceDir + "/web")
}