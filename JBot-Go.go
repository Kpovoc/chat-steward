package main

import (
	//"log"
	// "./IrcBot"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	//"os"
	//"os/signal"
	//"syscall"
	"github.com/Kpovoc/JBot-Go/adapter/discordbot"
)

type MainConf struct {
	Discord discordbot.DiscordConf
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

	errMsg, err := discordbot.CreateAndStartSession(conf.Discord)
	if errMsg != "" || err != nil {
		fmt.Printf("An Error has occured:\nMsg: %s\nErr: %s\n", errMsg, err)
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
