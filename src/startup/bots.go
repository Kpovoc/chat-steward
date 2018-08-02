package startup

import (
	"gitlab.com/Kpovoc/chat-steward/src/adapter/ircbot"
	"gitlab.com/Kpovoc/chat-steward/src/web"
	"log"
	)

func LaunchBots(resourceDir string) {
	// TODO: Implement with buffer channel so some parts can fail while others remain
	fatalChan := make(chan error)

	go ircbot.Start(conf.IRC, fatalChan)
	//go discordbot.CreateAndStartSession(conf.Discord)
	go web.Start(resourceDir + "/web", conf.WebSitePort)

	err := <- fatalChan
	if nil != err {
		log.Printf("Exited with error: %s\n", err)
	} else {
		log.Println("Exited with no problems.")
	}
}