package plugin

import (
	"gitlab.com/Kpovoc/chat-steward/src/core/message"
	"gitlab.com/Kpovoc/chat-steward/src/plugin/titles"
	"gitlab.com/Kpovoc/chat-steward/src/plugin/eightball"
	"math/rand"
)

func GetPluginResponse(pluginName string, msgContent string, m *message.Message) string {
	response := ""

	switch pluginName {
	case "8ball":
		fallthrough
	case "eightball":
		response = eightball.Plugin(msgContent, randNum)
	case "suggest":
		response = titles.Plugin(m, msgContent)
	case "start_show":
		response = titles.StartShow(m, msgContent)
	}

	return response
}

func InitPlugins(resourceDir string) {
	titles.Init()
	titles.WebInit(resourceDir + "/web")
}

func randNum(x int) int {
	return rand.Intn(x)
}