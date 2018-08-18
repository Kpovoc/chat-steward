package startup

import (
  "gitlab.com/Kpovoc/chat-steward/src/adapter/ircbot"
  "log"
  "gitlab.com/Kpovoc/chat-steward/src/web"
  "gitlab.com/Kpovoc/chat-steward/src/adapter/discordbot"
)

func LaunchBots(resourceDir string) {
  // TODO: Implement with buffer channel so some parts can fail while others remain
  fatalChan := make(chan error)

  go discordbot.CreateAndStartSession(conf.Discord, fatalChan)
  go ircbot.Start(conf.IRC, fatalChan)
  go web.Start(resourceDir + "/web", conf.WebSitePort)

  err := <- fatalChan
  if nil != err {
    log.Printf("Exited with error: %s\n", err)
  } else {
    log.Println("Exited with no problems.")
  }
}
