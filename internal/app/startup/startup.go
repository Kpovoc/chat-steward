package startup

import (
  "bufio"
  "encoding/json"
  "fmt"
  "gitlab.com/Kpovoc/chat-steward/internal/app/adapter/discordbot"
  "gitlab.com/Kpovoc/chat-steward/internal/app/adapter/ircbot"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin/info"
  "gitlab.com/Kpovoc/chat-steward/internal/app/util"
  "io"
  "io/ioutil"
  "log"
  "os"
)

type MainConf struct {
  Discord discordbot.DiscordConf
  IRC ircbot.IrcConf
  WebSitePort string
  Plugins plugin.PluginConf
}

var conf MainConf

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

  cw := ConfigWizard{}
  cw.Run()
  writeConfigToFile(cw.GetConf(), userDataDir)
}

type ConfigWizard struct {
  conf   MainConf
  reader io.Reader
  writer io.Writer
  scanner *bufio.Scanner
}

func (c *ConfigWizard) GetConf() MainConf {
  return c.conf
}

func (c *ConfigWizard) getUserInput() string {
  c.scanner.Scan()
  return c.scanner.Text()
}

func (c *ConfigWizard) Run() MainConf {
  if c.reader == nil {
    c.reader = os.Stdin
  }
  if c.writer == nil {
    c.writer = os.Stdout
  }

  c.scanner = bufio.NewScanner(c.reader)

  c.conf = MainConf {
    Discord: discordbot.DiscordConf{},
    IRC: ircbot.IrcConf{},
    WebSitePort: "",
    Plugins: plugin.PluginConf{},
  }

  c.runDiscordSetup()
  c.runIrcSetup()
  c.runWebSetup()
  c.runPluginSetup()

  return c.conf
}

func (c *ConfigWizard) runDiscordSetup() {
  dc := discordbot.DiscordConf {
    Admins: []discordbot.DiscordAdmin{},
    BotToken: "",
  }

  userAddDiscordBotToken := "y"
  fmt.Fprintf(c.writer, "Add Discord Bot Token? (y/n) [%s]: ", userAddDiscordBotToken)
  userAddDiscordBotToken = c.getUserInput()

  if "y" == userAddDiscordBotToken {
    fmt.Fprintf(c.writer, "Bot Token: ")
    dc.BotToken = c.getUserInput()

    dc.Admins = c.createDiscordAdminConf()
  }

  c.conf.Discord = dc
}

func (c *ConfigWizard) runIrcSetup() {
  ic := ircbot.IrcConf{}

  ic.Servers = append(ic.Servers, c.createIrcServerConf())

  userAddNewServerConf := "y"
  fmt.Fprintf(c.writer, "Add new IRC Server? (y/n) [%s]: ", userAddNewServerConf)
  userAddNewServerConf = c.getUserInput()

  for "y" == userAddNewServerConf {
    ic.Servers = append(ic.Servers, c.createIrcServerConf())

    fmt.Fprintf(c.writer, "Add new IRC Server? (y/n) [%s]: ", userAddNewServerConf)
    userAddNewServerConf = c.getUserInput()
  }

  c.conf.IRC = ic
}

func (c *ConfigWizard) runWebSetup() {
  port := "8080"
  fmt.Fprintf(c.writer, "Select a port number for the website [%s]: ", port)
  port = c.getUserInput()
  if port == "" || len(port) > 5 {
    fmt.Fprintf(c.writer, "\"%s\" is an incorrect port. Defaulting to \"8080\"\n", port)
    port = "8080"
  }

  c.conf.WebSitePort = port
}

func (c *ConfigWizard) runPluginSetup() {
  infoConf := []info.InfoConfItem{}

  addInfoConfItem := "y"
  fmt.Fprint(c.writer, "Add item to Info Plugin? (y/n): ")
  addInfoConfItem = c.getUserInput()

  for addInfoConfItem == "y" {
    fmt.Fprint(c.writer, "Term to fetch the desired information (Ex: irc): ")
    termStr := c.getUserInput()

    fmt.Fprint(c.writer, "Info to return when Term is entered\n(Ex: IRC Info - Server: irc.geekshed.net | Channel: #jupiterbroadcasting):\n")
    infoStr := c.getUserInput()

    infoConf = append(infoConf, info.InfoConfItem{termStr, infoStr})

    fmt.Fprint(c.writer, "Add item to Info Plugin? (y/n): ")
    addInfoConfItem = c.getUserInput()
  }

  c.conf.Plugins.Info = infoConf
}

func (c *ConfigWizard) createDiscordAdminConf() []discordbot.DiscordAdmin {
  admins := []discordbot.DiscordAdmin{}
  admin := discordbot.DiscordAdmin{}
  fmt.Fprint(c.writer, "Bot Admin Username for Discord (Ex: John): ")
  admin.Username = c.getUserInput()
  fmt.Fprint(c.writer, "Bot Admin Discriminator for Discord (Ex: 1234): ")
  admin.Discriminator = c.getUserInput()
  if admin.Username != "" && admin.Discriminator != "" {
    admins = append(admins, admin)
  } else {
    fmt.Fprintln(c.writer, "Having no Bot Admins will disable the functionality of some plugins. " +
      "It is strongly advised that you have at least one registered nickname entered as " +
      "a Bot Administrator.")
  }

  userAddAdmin := "y"
  fmt.Fprintf(c.writer, "Add a Bot Admin for Discord? (y/n) [%s]: ", userAddAdmin)
  userAddAdmin = c.getUserInput()

  for "y" == userAddAdmin {
    admin = discordbot.DiscordAdmin{}
    fmt.Fprint(c.writer, "Bot Admin Username for Discord (Ex: John): ")
    admin.Username = c.getUserInput()
    fmt.Fprint(c.writer, "Bot Admin Discriminator for Discord (Ex: 1234): ")
    admin.Discriminator = c.getUserInput()
    if admin.Username != "" && admin.Discriminator != "" {
      if !isDiscordAdminInArray(admin, admins) {
        admins = append(admins, admin)
      } else {
        fmt.Fprintln(c.writer, admin.Username + "#" + admin.Discriminator + " is already registered.")
      }
    } else if len(admins) <= 0 {
      fmt.Fprintln(c.writer, "Having no Bot Admins will disable the functionality of some plugins. " +
        "It is strongly advised that you have at least one registered nickname entered as " +
        "a Bot Administrator.")
    }

    fmt.Fprintf(c.writer, "Add a Bot Admin for Discord? (y/n) [%s]: ", userAddAdmin)
    userAddAdmin = c.getUserInput()
  }

  return admins
}

func (c *ConfigWizard) createIrcServerConf() ircbot.IrcServerConf {
  ircServConf := ircbot.IrcServerConf{}

  fmt.Fprintln(c.writer, "Enter IRC Server Information:")
  fmt.Fprint(c.writer, "Server (Ex: irc.geekshed.net): ")
  ircServConf.Server = c.getUserInput()

  fmt.Fprint(c.writer, "Port (Ex: 6667): ")
  ircServConf.Port = c.getUserInput()

  fmt.Fprint(c.writer, "Nick (Ex: Stuart): ")
  ircServConf.Nick = c.getUserInput()

  fmt.Fprint(c.writer, "User (Ex: stuart): ")
  ircServConf.User = c.getUserInput()

  fmt.Fprint(c.writer, "Password: ")
  ircServConf.Password = c.getUserInput()

  ircServConf.Channels = c.handleIrcServerChannelsConf()

  ircServConf.AdminNicks = c.handleIrcServerAdminNicksConf()

  return ircServConf
}

func (c *ConfigWizard) handleIrcServerChannelsConf() []string {
  channels := []string{}

  fmt.Fprint(c.writer, "Channel Name (Ex: #jupitercolony): ")
  channel := c.getUserInput()
  channels = append(channels, channel)

  userAddChannel := "y"
  fmt.Fprintf(c.writer, "Add another channel for this Server? (y/n) [%s]: ", userAddChannel)
  userAddChannel = c.getUserInput()

  for "y" == userAddChannel {
    channel = ""
    fmt.Fprint(c.writer, "Channel Name (Ex: #jupitercolony): ")
    channel = c.getUserInput()
    channels = append(channels, channel)

    fmt.Fprintf(c.writer, "Add another channel for this Server? (y/n) [%s]: ", userAddChannel)
    userAddChannel = c.getUserInput()
  }

  return channels
}

func (c *ConfigWizard) handleIrcServerAdminNicksConf() []string {
  nicks := []string{}

  fmt.Fprint(c.writer, "Bot Admin Nickname for Server: ")
  nick := c.getUserInput()
  if nick != "" {
    nicks = append(nicks, nick)
  } else {
    fmt.Fprintln(c.writer, "Having no Bot Admins will disable the functionality of some plugins. " +
      "It is strongly advised that you have at least one registered nickname entered as " +
      "a Bot Administrator.")
  }

  userAddAdmin := "y"
  fmt.Fprintf(c.writer, "Add a Bot Admin Nickname for this server? (y/n) [%s]: ", userAddAdmin)
  userAddAdmin = c.getUserInput()

  for "y" == userAddAdmin {
    nick = ""
    fmt.Fprint(c.writer, "Bot Admin Nickname for Server: ")
    nick = c.getUserInput()

    if nick != "" {
      if !util.IsStringInArray(nick, nicks) {
        nicks = append(nicks, nick)
      } else {
        fmt.Fprintln(c.writer, nick + " is already registered.")
      }
    } else if len(nicks) <= 0 {
      fmt.Fprintln(c.writer, "Having no Bot Admins will disable the functionality of some plugins. " +
        "It is strongly advised that you have at least one registered nickname entered as " +
        "a Bot Administrator.")
    }

    fmt.Fprintf(c.writer, "Add a Bot Admin Nickname for this server? (y/n) [%s]: ", userAddAdmin)
    userAddAdmin = c.getUserInput()
  }

  return nicks
}


func writeConfigToFile(config MainConf, userDataDir string) {
  js, err := json.MarshalIndent(config, "", "  ")
  if err != nil {
    log.Fatalln("An error occured while Marshalling the user configuration. Error: ", err)
  }

  _ = ioutil.WriteFile(userDataDir + "/main-conf.json", js, 0644)
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
