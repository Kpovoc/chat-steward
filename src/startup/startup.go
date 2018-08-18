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
  var config mainConf
  config.Discord = handleDiscordConf()
  config.IRC = handleIrcConf()
  config.WebSitePort = handleWebSitePortConf()
  writeConfigToFile(config, userDataDir)

}

func handleDiscordConf() discordbot.DiscordConf {
  ret := discordbot.DiscordConf{
    Admins: []discordbot.DiscordAdmin{},
    BotToken: "",
  }

  userAddDiscordBotToken := "y"
  fmt.Printf("Add Discord Bot Token? (y/n) [%s]: ", userAddDiscordBotToken)
  fmt.Scanf("%s", &userAddDiscordBotToken)

  if "y" == userAddDiscordBotToken {
    fmt.Printf("Bot Token: ")
    fmt.Scanf("%s", &ret.BotToken)

    ret.Admins = handleDiscordAdminConf()
  }

  return ret
}

func handleIrcConf() ircbot.IrcConf {
  var ret ircbot.IrcConf

  ret.Servers = append(ret.Servers, createIrcServerConf())

  userAddNewServerConf := "y"
  fmt.Printf("Add new IRC Server? (y/n) [%s]: ", userAddNewServerConf)
  fmt.Scanf("%s", &userAddNewServerConf)

  for "y" == userAddNewServerConf {
    ret.Servers = append(ret.Servers, createIrcServerConf())

    fmt.Printf("Add new IRC Server? (y/n) [%s]: ", userAddNewServerConf)
    fmt.Scanf("%s", &userAddNewServerConf)
  }

  return ret
}

func handleWebSitePortConf() string {
  port := "8080"
  fmt.Printf("Select a port number for the website [%s]: ", port)
  fmt.Scanf("%s", &port)
  if port == "" || len(port) > 5 {
    fmt.Printf("\"%s\" is an incorrect port. Defaulting to \"8080\"\n", port)
    port = "8080"
  }

  return port
}

func writeConfigToFile(config mainConf, userDataDir string) {
  js, err := json.MarshalIndent(config, "", "  ")
  if err != nil {
    log.Fatalln("An error occured while Marshalling the user configuration. Error: ", err)
  }

  _ = ioutil.WriteFile(userDataDir + "/main-conf.json", js, 0644)
}

func createIrcServerConf() ircbot.IrcServerConf {
  var ircServConf ircbot.IrcServerConf

  fmt.Println("Enter IRC Server Information:")
  fmt.Print("Server (Ex: irc.geekshed.net): ")
  fmt.Scanf("%s", &ircServConf.Server)

  fmt.Print("Port (Ex: 6667): ")
  fmt.Scanf("%s", &ircServConf.Port)

  fmt.Print("Nick (Ex: Stuart): ")
  fmt.Scanf("%s", &ircServConf.Nick)

  fmt.Print("User (Ex: stuart): ")
  fmt.Scanf("%s", &ircServConf.User)

  fmt.Print("Password: ")
  fmt.Scanf("%s", &ircServConf.Password)

  ircServConf.Channels = handleIrcServerChannelsConf()

  ircServConf.AdminNicks = handleIrcServerAdminNicksConf()

  return ircServConf
}

func handleIrcServerChannelsConf() []string {
  channels := []string{}
  channel := ""
  fmt.Print("Channel Name (Ex: #jupitercolony): ")
  fmt.Scanf("%s", &channel)
  channels = append(channels, channel)

  userAddChannel := "y"
  fmt.Printf("Add another channel for this Server? (y/n) [%s]: ", userAddChannel)
  fmt.Scanf("%s", &userAddChannel)

  for "y" == userAddChannel {
    channel = ""
    fmt.Print("Channel Name (Ex: #jupitercolony): ")
    fmt.Scanf("%s", &channel)
    channels = append(channels, channel)

    fmt.Printf("Add another channel for this Server? (y/n) [%s]: ", userAddChannel)
    fmt.Scanf("%s", &userAddChannel)
  }

  return channels
}

func handleIrcServerAdminNicksConf() []string {
  nicks := []string{}
  nick := ""
  fmt.Print("Bot Admin Nickname for Server: ")
  fmt.Scanf("%s", &nick)
  if nick != "" {
    nicks = append(nicks, nick)
  } else {
    fmt.Println("Having no Bot Admins will disable the functionality of some plugins. " +
      "It is strongly advised that you have at least one registered nickname entered as " +
      "a Bot Administrator.")
  }

  userAddAdmin := "y"
  fmt.Printf("Add a Bot Admin Nickname for this server? (y/n) [%s]: ", userAddAdmin)
  fmt.Scanf("%s", &userAddAdmin)

  for "y" == userAddAdmin {
    nick = ""
    fmt.Print("Bot Admin Nickname for Server: ")
    fmt.Scanf("%s", &nick)

    if nick != "" {
      if !util.IsStringInArray(nick, nicks) {
        nicks = append(nicks, nick)
      } else {
        fmt.Println(nick + " is already registered.")
      }
    } else if len(nicks) <= 0 {
      fmt.Println("Having no Bot Admins will disable the functionality of some plugins. " +
        "It is strongly advised that you have at least one registered nickname entered as " +
        "a Bot Administrator.")
    }

    fmt.Printf("Add a Bot Admin Nickname for this server? (y/n) [%s]: ", userAddAdmin)
    fmt.Scanf("%s", &userAddAdmin)
  }

  return nicks
}

func handleDiscordAdminConf() []discordbot.DiscordAdmin {
  admins := []discordbot.DiscordAdmin{}
  admin := discordbot.DiscordAdmin{}
  fmt.Print("Bot Admin Username for Discord (Ex: John): ")
  fmt.Scanf("%s", &admin.Username)
  fmt.Print("Bot Admin Discriminator for Discord (Ex: 1234): ")
  fmt.Scanf("%s", &admin.Discriminator)
  if admin.Username != "" && admin.Discriminator != "" {
    admins = append(admins, admin)
  } else {
    fmt.Println("Having no Bot Admins will disable the functionality of some plugins. " +
      "It is strongly advised that you have at least one registered nickname entered as " +
      "a Bot Administrator.")
  }

  userAddAdmin := "y"
  fmt.Printf("Add a Bot Admin for Discord? (y/n) [%s]: ", userAddAdmin)
  fmt.Scanf("%s", &userAddAdmin)

  for "y" == userAddAdmin {
    admin = discordbot.DiscordAdmin{}
    fmt.Print("Bot Admin Username for Discord (Ex: John): ")
    fmt.Scanf("%s", &admin.Username)
    fmt.Print("Bot Admin Discriminator for Discord (Ex: 1234): ")
    fmt.Scanf("%s", &admin.Discriminator)
    if admin.Username != "" && admin.Discriminator != "" {
      if !isDiscordAdminInArray(admin, admins) {
        admins = append(admins, admin)
      } else {
        fmt.Println(admin.Username + "#" + admin.Discriminator + " is already registered.")
      }
    } else if len(admins) <= 0 {
      fmt.Println("Having no Bot Admins will disable the functionality of some plugins. " +
        "It is strongly advised that you have at least one registered nickname entered as " +
        "a Bot Administrator.")
    }

    fmt.Printf("Add a Bot Admin for Discord? (y/n) [%s]: ", userAddAdmin)
    fmt.Scanf("%s", &userAddAdmin)
  }

  return admins
}

func isDiscordAdminInArray(admin discordbot.DiscordAdmin, admins []discordbot.DiscordAdmin) bool {
  for i:=0;i<len(admins);i++ {
    a := admins[i]
    if admin.Username == a.Username && admin.Discriminator == a.Discriminator {
      return true
    }
  }
  return false
}
