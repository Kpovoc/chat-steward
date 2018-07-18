package startup

import (
	"gitlab.com/Kpovoc/chat-steward/src/adapter/ircbot"
	"gitlab.com/Kpovoc/chat-steward/src/adapter/discordbot"
	"io/ioutil"
	"fmt"
	"log"
	"encoding/json"
	"gitlab.com/Kpovoc/chat-steward/src/util"
)

type mainConf struct {
	Discord discordbot.DiscordConf
	IRC ircbot.IrcConf
	WebSitePort string
}

var conf mainConf

func InitializeConfig(userDataDir string) {
	data, err := ioutil.ReadFile(userDataDir + "/main-conf.json")
	if nil != err {
		// Either error with current main-conf.json file, or file does not exist.
		askUserToStartConfigWizard(userDataDir)

		data, err = ioutil.ReadFile(userDataDir + "/main-conf.json")
		if nil != err {
			util.ExitWithError(err)
		}
	}
	if err = json.Unmarshal(data, &conf); err != nil {
		util.ExitWithError(err)
	}
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
		"  },\n"
	websitePortConfStart :=
		"  \"WebSitePort\": \""
	websitePortConfEnd :=   "\"\n"
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

	sitePort := "8080"
	fmt.Printf("Select a port number for the website [%s]: ", sitePort)
	fmt.Scanf("%s", &sitePort)

	if sitePort == "" || len(sitePort) > 5 {
		fmt.Printf("\"%s\" is an incorrect port. Defaulting to \"8080\"\n", sitePort)
		sitePort = "8080"
	}


	out = out +
		ircServerConfEnd +
		ircConfEnd +
		websitePortConfStart + sitePort + websitePortConfEnd +
		mainConfEnd

	_ = ioutil.WriteFile(userDataDir + "/main-conf.json", []byte(out), 0644)
}