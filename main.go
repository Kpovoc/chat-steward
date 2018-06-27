package main

import (
	"gitlab.com/Kpovoc/chat-steward/adapter/discordbot"
	"gitlab.com/Kpovoc/chat-steward/adapter/ircbot"
	"flag"
	"time"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	"gitlab.com/Kpovoc/chat-steward/plugin"
	"log"
	"gitlab.com/Kpovoc/chat-steward/web"
	"fmt"
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

	// Go ahead and seed rand for plugins that need it.
	rand.Seed(time.Now().UTC().UnixNano())

	// Initialize the application config from file.
	conf := readInitialConfig(userDataDir)

	// Initialize application plugins
	plugin.InitPlugins(resourceDir)

	// TODO: Implement with buffer channel so some parts can fail while others remain
	fatalChan := make(chan error)

	go ircbot.Start(conf.IRC, fatalChan)
	go web.Start(resourceDir + "/web")

	err := <- fatalChan
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

func ExitWithError(errMsg string) {
	log.Fatalf("An error has occurred:\n  %s\n  Application will now exit.", errMsg)
}

func readInitialConfig(userDataDir string) MainConf {
	var conf MainConf

	data, err := ioutil.ReadFile(userDataDir + "/main-conf.json")
	if nil != err {
		// Either error with current main-conf.json file, or file does not exist.
		askUserToStartConfigWizard(userDataDir)

		data, err = ioutil.ReadFile(userDataDir + "/main-conf.json")
		if nil != err {
			ExitWithError(err.Error())
		}
	}

	if err = json.Unmarshal(data, &conf); err != nil {
		ExitWithError(err.Error())
	}

	return conf
}

func askUserToStartConfigWizard(userDataDir string) {
	userChoice := "y"
	fmt.Printf(
		"Could not read %s/main-conf.json. " +
			"Would you like to initiate the setup wizard? (y/n) [%s]: ", userDataDir, userChoice)
	fmt.Scanf("%s", &userChoice)

	for "y" != userChoice && "n" != userChoice {
		fmt.Printf("Unrecognized choice: %s\n", userChoice)
		fmt.Printf("Would you like to initiate the setup wizard? (y/n) [%s]: ", userChoice)
		fmt.Scanf("%s", &userChoice)
	}

	if ("n" == userChoice) {
		log.Fatalln(
			"ChatSteward can not operate without a valid main-conf.json. " +
				"Application will now exit.")
	}

	createConfig(userDataDir)
}

func createConfig(userDataDir string) {
	mainConfStart :=
		"{\n"
	discordConf :=
		"  \"Discord\": {\n" +
		"    \"BotToken\": \"\"\n" +
		"  },\n"
	ircConfStart :=
		"  \"IRC\": {\n" +
		"    \"Servers\": [\n"
	ircServerConfStart :=
		"      {\n"
	ircServerConfServer :=
		"        \"Server\": "
	ircServerConfPort :=
		"        \"Port\": "
	ircServerConfNick :=
		"        \"Nick\": "
	ircServerConfUser :=
		"        \"User\": "
	ircServerConfChannelsStart :=
		"        \"Channels\": [\n"
	ircServerConfChannelsEnd :=
		"        ]\n"
	ircServerConfEndComma :=
		"      },\n"
	ircServerConfEnd :=
		"      }\n"
	ircConfEnd :=
		"    ]\n" +
		"  }\n"
	mainConfEnd :=
		"}"

	out := mainConfStart +
		discordConf +
		ircConfStart +
		ircServerConfStart +
		ircServerConfServer

	server := ""
	port := ""
	nick := ""
	user := ""
	channel := ""
	fmt.Println("Enter IRC Server Information:")

	fmt.Print("Server (Ex: irc.geekshed.net): ")
	fmt.Scanf("%s", &server)
	out = out + "\"" + server + "\",\n" +
		ircServerConfPort
	server = ""

	fmt.Print("Port (Ex: 6667): ")
	fmt.Scanf("%s", &port)
	out = out + "\"" + port + "\",\n" +
		ircServerConfNick
	port = ""

	fmt.Print("Nick (Ex: Stuart): ")
	fmt.Scanf("%s", &nick)
	out = out + "\"" + nick + "\",\n" +
		ircServerConfUser
	nick = ""

	fmt.Print("User (Ex: stuart): ")
	fmt.Scanf("%s", &user)
	out = out + "\"" + user + "\",\n" +
		ircServerConfChannelsStart
	user = ""

	fmt.Print("Channel Name (Ex: #jupitercolony): ")
	fmt.Scanf("%s", &channel)
	out = out +
		"          \"" + channel + "\""
	channel = ""

	userAddChannel := "y"
	fmt.Printf("Add another channel for this Server? (y/n) [%s]: ", userAddChannel)
	fmt.Scanf("%s", &userAddChannel)

	for "y" == userAddChannel {
		out = out + ",\n"
		fmt.Print("Channel Name (Ex: #jupitercolony): ")
		fmt.Scanf("%s", &channel)
		out = out +
			"          \"" + channel + "\""
		channel = ""

		fmt.Printf("Add another channel for this Server? (y/n) [%s]: ", userAddChannel)
		fmt.Scanf("%s", &userAddChannel)
	}
	out = out + "\n" +
		ircServerConfChannelsEnd

	userAddNewServerConf := "y"
	fmt.Printf("Add new IRC Server? (y/n) [%s]: ", userAddNewServerConf)
	fmt.Scanf("%s", &userAddNewServerConf)

	for "y" == userAddNewServerConf {
		out = out +
			ircServerConfEndComma +
			ircServerConfStart +
			ircServerConfServer

		fmt.Print("Server (Ex: irc.geekshed.net): ")
		fmt.Scanf("%s", &server)
		out = out + "\"" + server + "\",\n" +
			ircServerConfPort
		server = ""

		fmt.Print("Port (Ex: 6667): ")
		fmt.Scanf("%s", &port)
		out = out + "\"" + port + "\",\n" +
			ircServerConfNick
		port = ""

		fmt.Print("Nick (Ex: Stuart): ")
		fmt.Scanf("%s", &nick)
		out = out + "\"" + nick + "\",\n" +
			ircServerConfUser
		nick = ""

		fmt.Print("User (Ex: stuart): ")
		fmt.Scanf("%s", &user)
		out = out + "\"" + user + "\",\n" +
			ircServerConfChannelsStart
		user = ""

		fmt.Print("Channel Name (Ex: #jupitercolony): ")
		fmt.Scanf("%s", &channel)
		out = out +
			"          \"" + channel + "\""
		channel = ""

		userAddChannel := "y"
		fmt.Printf("Add another channel for this Server? (y/n) [%s]: ", userAddChannel)
		fmt.Scanf("%s", &userAddChannel)

		for "y" == userAddChannel {
			out = out + ",\n"
			fmt.Print("Channel Name (Ex: #jupitercolony): ")
			fmt.Scanf("%s", &channel)
			out = out +
				"          \"" + channel + "\""
			channel = ""

			fmt.Printf("Add another channel for this Server? (y/n) [%s]: ", userAddChannel)
			fmt.Scanf("%s", &userAddChannel)
		}
		out = out + "\n" +
			ircServerConfChannelsEnd
	}
	out = out +
		ircServerConfEnd +
		ircConfEnd +
		mainConfEnd

	_ = ioutil.WriteFile(userDataDir + "/main-conf.json", []byte(out), 0644)
}
