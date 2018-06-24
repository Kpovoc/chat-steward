package main

import (
	"gitlab.com/Kpovoc/JBot-Go/adapter/discordbot"
	"gitlab.com/Kpovoc/JBot-Go/adapter/ircbot"
	"flag"
	"time"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	"gitlab.com/Kpovoc/JBot-Go/plugin"
	"log"
	"gitlab.com/Kpovoc/JBot-Go/web"
)

type MainConf struct {
	Discord discordbot.DiscordConf
	IRC ircbot.IrcConf
}

func main() {
	rd := flag.String(
		"resources",
		"resources",
		"directory containing resources needed for application")
	udd := flag.String(
		"userdata",
		"data",
		"directory containing config file and database")
	flag.Parse()
	resourceDir := *rd
	userDataDir := *udd

	rand.Seed(time.Now().UTC().UnixNano()) // Go ahead and seed rand for plugins that need it.
	data, err := ioutil.ReadFile(userDataDir + "/main-conf.json")
	if err != nil {
		panic(err)
	}

	var conf MainConf
	if err = json.Unmarshal(data, &conf); err != nil {
		panic(err)
	}

	plugin.InitPlugins(resourceDir)

	// TODO: Implement with buffer channel so some parts can fail while others remain
	fatalChan := make(chan error)

	go ircbot.Start(conf.IRC, fatalChan)
	go web.Start(resourceDir + "/web")

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
