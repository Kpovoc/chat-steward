package core

import (
	"strings"
	"gitlab.com/Kpovoc/chat-steward/plugin"
	"gitlab.com/Kpovoc/chat-steward/core/message"
)

func GenerateResponse(m *message.Message) string {
	content := m.Content
	if content[0] != '!' {
		return ""
	}

	pluginName, msgContent := parseContent(content)
	if pluginName == "" {
		return ""
	}

	return plugin.GetPluginResponse(pluginName, msgContent, m)
}

func parseContent(content string) (string, string) {
	content = strings.TrimSpace(content)
	args := strings.SplitN(content[1:], " ", 2)
	pluginName := args[0]
	msgContent := ""

	if len(args) > 1 {
		msgContent = args[1]
	}

	return pluginName, msgContent
}

