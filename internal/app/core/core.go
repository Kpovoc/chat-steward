package core

import (
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/response"
  "strings"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin"
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/message"
)

func GenerateResponse(m *message.Message) *response.Response {
  content := m.Content
  if content[0] != '!' {
    return nil
  }

  pluginName, msgContent := parseContent(content)
  if pluginName == "" {
    return nil
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

