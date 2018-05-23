package plugin

import (
	"github.com/Kpovoc/JBot-Go/core/message"
	"github.com/Kpovoc/JBot-Go/plugin/titles"
	"github.com/Kpovoc/JBot-Go/plugin/show"
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

func InitPlugins() {
	titles.WebInit()
}