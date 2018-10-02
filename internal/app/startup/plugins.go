package startup

import (
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin/info"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin/titles"
)

func InitializePlugins(resourceDir string) {
  titles.Init()
  titles.WebInit(resourceDir + "/web")

  // TODO: Move to top of function when plugin confs get implemented.
  pluginConf := &conf.Plugins
  if pluginConf == nil {
    return
  }

  info.Init(pluginConf.Info)
}
