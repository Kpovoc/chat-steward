package main

import (
	"github.com/Kpovoc/JBot-Go/adapter/discordbot"
	"github.com/Kpovoc/JBot-Go/adapter/ircbot"
	"math/rand"
	"time"
	"io/ioutil"
	"encoding/json"
	"log"
	"github.com/Kpovoc/JBot-Go/plugin"
	"github.com/Kpovoc/JBot-Go/web"
)

type MainConf struct {
	Discord discordbot.DiscordConf
	IRC ircbot.IrcConf
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // Go ahead and seed rand for plugins that need it.
	data, err := ioutil.ReadFile("./resources/configs/main-conf.json")
	if err != nil {
		panic(err)
	}

	var conf MainConf
	if err = json.Unmarshal(data, &conf); err != nil {
		panic(err)
	}

	plugin.InitPlugins()

	// TODO: Implement with buffer channel so some parts can fail while others remain
	fatalChan := make(chan error)

	go ircbot.Start(conf.IRC, fatalChan)
	go web.Start()

	err = <- fatalChan
	if nil != err {
		log.Printf("Exited with error: %s\n", err)
	} else {
		log.Println("Exited with no problems.")
	}

	//errMsg, err := discordbot.CreateAndStartSession(conf.Discord)
	//if errMsg != "" || err != nil {
	//	fmt.Printf("An Error has occured:\nMsg: %s\nErr: %s\n", errMsg, err)
	//	return
	//}
}
