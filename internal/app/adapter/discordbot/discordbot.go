package discordbot

import (
  "fmt"
  "os"
  "os/signal"
  "syscall"

  "gitlab.com/Kpovoc/chat-steward/internal/app/core"
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/message"
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/user"
  "github.com/bwmarrin/discordgo"
)

type DiscordConf struct {
  BotToken string
  Admins []DiscordAdmin
}

type DiscordAdmin struct {
  Username string
  Discriminator string
}

type DiscordBot struct {
  Token string
  Admins []DiscordAdmin
  Session *discordgo.Session
  ErrChan chan error
}

func CreateAndStartSession(conf DiscordConf, fatalErr chan error) {
  if conf.BotToken == "" {
    fmt.Println("No Bot Token given.")
    return
  }

  session, err := discordgo.New("Bot " + conf.BotToken)
  if err != nil {
    fmt.Printf(
      "An Error has occured:\nMsg: %s\nErr: %s\n",
      "Could not create Discord session",
      err)
    return
  }

  admins := conf.Admins
  if admins == nil {
    admins = []DiscordAdmin{}
  }

  discordBot := &DiscordBot{
    Token: conf.BotToken,
    Admins: admins,
    Session: session,
    ErrChan: fatalErr,
  }

  err = discordBot.Run()
  if err != nil {
    fmt.Printf(
      "An Error has occured:\nMsg: %s\nErr: %s\n",
      "Something went wrong during Discord run",
      err)
  }
  discordBot.ErrChan <- err
}

func (b *DiscordBot) Run() error {
  session := b.Session

  // Register the messageCreate func as a callback for MessageCreate events.
  session.AddHandler(b.messageCreate)

  // Open a websocket connection to Discord and begin listening.
  err := session.Open()

  if err != nil {
    fmt.Println("error opening connection,", err)
    return err
  }

  // Wait here until CTRL-C or other term signal is received.
  fmt.Println("Bot is now running.  Press CTRL-C to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc

  // Cleanly close down the Discord session.
  session.Close()

  return err
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (bot *DiscordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

  // Ignore all messages created by the bot itself
  // This isn't required in this specific example but it's a good practice.
  if m.Author.ID == s.State.User.ID {
    return
  }

  coreMsg := bot.convertToCoreMessage(s, m)

  response := core.GenerateResponse(coreMsg)

  if response != "" {
    s.ChannelMessageSend(m.ChannelID, response)
  }
}

func (bot *DiscordBot) convertToCoreMessage(s *discordgo.Session, dMsg *discordgo.MessageCreate) *message.Message {
  sender := user.New( // Call Read Function later
    "",
    dMsg.Author.ID,
    dMsg.Author.Username,
    "",
    "",
    "",
    bot.isAdmin(dMsg.Author.Username, dMsg.Author.Discriminator))

  // testPrintDGOMessage(s, dMsg)

  content := dMsg.Content
  createdOn, _ := dMsg.Timestamp.Parse()

  return &message.Message{
    Sender: sender,
    CreatedOn: createdOn,
    Content: content,
  }
}

func testPrintDGOMessage(s *discordgo.Session, dMsg *discordgo.MessageCreate) {
  // discordgo.Message struct
  // strings
  id := dMsg.ID
  channelId := dMsg.ChannelID
  content := dMsg.Content
  cwmr := dMsg.ContentWithMentionsReplaced()
  cwmmr, _ := dMsg.ContentWithMoreMentionsReplaced(s)
  fmt.Println("Strings:")
  fmt.Println("msgId       : ", id)
  fmt.Println("channelId   : ", channelId)
  fmt.Println("content     : ", content)
  fmt.Println("contentWMR  : ", cwmr)
  fmt.Println("contentWMMR : ", cwmmr)
  fmt.Println()

  // bool
  tts := dMsg.Tts // bool
  mentionEveryone := dMsg.MentionEveryone // bool
  fmt.Println("Booleans:")
  fmt.Println("Tts             : ", tts)
  fmt.Println("MentionEveryone : ", mentionEveryone)
  fmt.Println()

  // []string
  mentionRoles := dMsg.MentionRoles
  fmt.Println("[]string MentionRoles:")
  for index, str := range mentionRoles {
    fmt.Println(index, ": ", str)
  }
  fmt.Println()

  // discordgo.Timestamp
  timestamp := dMsg.Timestamp
  editedTimeStamp := dMsg.EditedTimestamp
  fmt.Println("Timestamps:")
  fmt.Println("timestamp       : ", timestamp)
  fmt.Println("editedTimeStamp : ", editedTimeStamp)

  // discordgo.User
  author := dMsg.Author // *discordgo.User
  testPrintDGOUser(author, "Author")
  mentions := dMsg.Mentions // []*discordgo.User
  for index, dUser := range mentions {
    testPrintDGOUser(dUser, "Mentioned " + string(index))
  }

  // discordgo.MessageAttachment
  attachments := dMsg.Attachments // []*discordgo.MessageAttachment
  for _, attachment := range attachments {
    testPrintDGOAttachment(attachment)
  }
  // embeds := dMsg.Embeds // []*discordgo.MessageEmbed
  //reactions := dMsg.Reactions // []*discordgo.MessageReactions
  //dmsgType := dMsg.Type  // discordgo.MessageType
}

func testPrintDGOUser(user *discordgo.User, title string) {
  fmt.Println("User ", title, ":")
  fmt.Println("userId        : ", user.ID)
  fmt.Println("email         : ", user.Email)
  fmt.Println("username      : ", user.Username)
  fmt.Println("avatar        : ", user.Avatar)
  fmt.Println("discriminator : ", user.Discriminator)
  fmt.Println("token         : ", user.Token)
  fmt.Println("verified      : ", user.Verified)
  fmt.Println("mfaEnabled    : ", user.MFAEnabled)
  fmt.Println("bot           : ", user.Bot)
  fmt.Println()
}

func testPrintDGOAttachment(at *discordgo.MessageAttachment) {
  fmt.Println("Attachment:")
  fmt.Println("attachmentId : ", at.ID)
  fmt.Println("url          : ", at.URL)
  fmt.Println("proxyURL     : ", at.ProxyURL)
  fmt.Println("filename     : ", at.Filename)
  fmt.Println("width        : ", at.Width)
  fmt.Println("height       : ", at.Height)
  fmt.Println("size         : ", at.Size)
  fmt.Println()
}

func (bot *DiscordBot) isAdmin(username string, discriminator string) bool {
  for i:=0;i<len(bot.Admins);i++ {
    admin := bot.Admins[i]
    if username == admin.Username && discriminator == admin.Discriminator {
      return true
    }
  }

  return false
}
