package main

import (
	//"log"
	// "./IrcBot"

	"fmt"
	//"os"
	//"os/signal"
	//"syscall"
	"github.com/Kpovoc/JBot-Go/adapter/discordbot"
	"math/rand"
	"time"
	"encoding/json"
	"io/ioutil"
)
type MainConf struct {
	Discord DiscordConf
}

type DiscordConf struct {
	BotToken string
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // Go ahead and seed rand for plugins that need it.
	data, err := ioutil.ReadFile("./resources/configs/main.conf")
	if err != nil {
		panic(err)
	}

	var conf MainConf
	if err = json.Unmarshal(data, &conf); err != nil {
		panic(err)
	}

	discordObj := conf.Discord
	token := discordObj.BotToken

	discordBot, err := discordbot.New(token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	err = discordBot.Run()
	if err != nil {
		fmt.Println("Error occured during Run: ", err)
		return
	}

	//
	//Session.Close()
	//japawig := IrcBot.NewJbIrcBot()
	////japawig.SetChannel("#unfilter")
	//err := japawig.Launch()
	//if err != nil {
	//	log.Printf("Exited with error: %s\n", err)
	//} else {
	//	log.Println("Exited with no problems.")
	//}
}
