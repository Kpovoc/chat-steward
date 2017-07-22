package core

import (
	"strings"
	"github.com/Kpovoc/JBot-Go/plugin"
)

func GenerateResponse(m *Message) string {
	content := m.Content
	if content[0] != '!' {
		return ""
	}

	pluginName, args := parseContent(content)
	if pluginName == "" {
		return ""
	}

	return plugin.GetPluginResponse(pluginName, args)
}

func parseContent(content string) (string, []string) {
	content = strings.TrimSpace(content)
	args := strings.Split(content[1:], " ")
	pluginName := args[0]
	return pluginName, args[1:]
}

