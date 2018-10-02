package startup

import (
  "bytes"
  "gitlab.com/Kpovoc/chat-steward/internal/app/adapter/discordbot"
  "gitlab.com/Kpovoc/chat-steward/internal/app/adapter/ircbot"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin/info"
  "bufio"
  "reflect"
  "testing"
)

func TestConfigWizard_RunDiscordSetup(t *testing.T) {
  userInputs := bytes.NewBufferString(`y
000000000000000000000000.000000.00_000000000000000000000000
TestUser
1234
n`)
  out := &bytes.Buffer{}
  cw := ConfigWizard{
    conf: MainConf{},
    reader: userInputs,
    writer: out,
    scanner: bufio.NewScanner(userInputs),
  }
  cw.runDiscordSetup()

  want := discordbot.DiscordConf {
    BotToken: "000000000000000000000000.000000.00_000000000000000000000000",
    Admins: []discordbot.DiscordAdmin{
      {Username: "TestUser", Discriminator: "1234"},
    },
  }
  got := cw.GetConf().Discord

  if !reflect.DeepEqual(want, got) {
    t.Errorf("\ngot: '%v'\nwant: '%v'", got, want)
  }
}

func TestConfigWizard_RunIrcSetup(t *testing.T) {
  userInputs := bytes.NewBufferString(`irc.test.net
6667
TestNick
testuser
password
#channel1
y
#channel2
n
TestAdminNick
n
n`)
  out := &bytes.Buffer{}
  cw := ConfigWizard{
    conf: MainConf{},
    reader: userInputs,
    writer: out,
    scanner: bufio.NewScanner(userInputs),
  }
  cw.runIrcSetup()

  want := ircbot.IrcConf{
    Servers: []ircbot.IrcServerConf{
      {
        Server: "irc.test.net",
        Port: "6667",
        Nick: "TestNick",
        User: "testuser",
        Password: "password",
        Channels: []string{
          "#channel1",
          "#channel2",
        },
        AdminNicks: []string {
          "TestAdminNick",
        },
      },
    },
  }
  got := cw.GetConf().IRC

  if !reflect.DeepEqual(want, got) {
    t.Errorf("\ngot: '%v'\nwant: '%v'", got, want)
  }
}

func TestConfigWizard_RunWebSetup(t *testing.T) {
  userInputs := bytes.NewBufferString("9000\n")
  out := &bytes.Buffer{}
  cw := ConfigWizard{
    conf: MainConf{},
    reader: userInputs,
    writer: out,
    scanner: bufio.NewScanner(userInputs),
  }

  cw.runWebSetup()

  want := "9000"
  got := cw.GetConf().WebSitePort

  if want != got {
    t.Errorf("\ngot: '%s'\nwant: '%v'", got, want)
  }
}

func TestConfigWizard_RunPluginSetup(t *testing.T) {
  userInputs := bytes.NewBufferString(`y
discord
Discord Info - Join at https://discord.me/jupitercolony
y
irc
IRC Info - Server: irc.geekshed.net | Channel: #jupiterbroadcasting
y
mumble
Mumble Info - Server: mumble.jupiterbroadcasting.org | Port: 64734
n`)
  out := &bytes.Buffer{}
  cw := ConfigWizard{
    conf: MainConf{},
    reader: userInputs,
    writer: out,
    scanner: bufio.NewScanner(userInputs),
  }
  cw.runPluginSetup()

  want := plugin.PluginConf{
    Info: []info.InfoConfItem{
      {Term: "discord", Info: "Discord Info - Join at https://discord.me/jupitercolony"},
      {Term: "irc", Info: "IRC Info - Server: irc.geekshed.net | Channel: #jupiterbroadcasting"},
      {Term: "mumble", Info: "Mumble Info - Server: mumble.jupiterbroadcasting.org | Port: 64734"},
    },
  }
  got := cw.GetConf().Plugins

  if !reflect.DeepEqual(want, got) {
    t.Errorf("\ngot: '%v'\nwant: '%v'", got, want)
  }
}
