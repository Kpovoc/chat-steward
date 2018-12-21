package plugin

import (
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/message"
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/response"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin/info"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin/show"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin/titles"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin/eightball"
  "math/rand"
)

type PluginConf struct {
  Info []info.InfoConfItem
}

func GetPluginResponse(pluginName string, msgContent string, m *message.Message) *response.Response {

  switch pluginName {
  case "8ball":
    fallthrough
  case "eightball":
    return eightball.Plugin(msgContent, randNum)
  case "suggest":
    return titles.Plugin(m, msgContent)
  case "show":
    return show.Plugin(m.Sender, msgContent)
  case "info":
    return info.Plugin(msgContent)
  default:
    return nil
  }
}

func randNum(x int) int {
  return rand.Intn(x)
}
