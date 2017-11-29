package discordbot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/Kpovoc/JBot-Go/core"
	"github.com/Kpovoc/JBot-Go/core/user"
	"github.com/Kpovoc/JBot-Go/core/message"
	"github.com/satori/go.uuid"
)

type DiscordBot struct {
	Token string
	Session *discordgo.Session
}

func New(t string) (*DiscordBot, error) {
	s, err := discordgo.New("Bot " + t)

	var bot = &DiscordBot{
		Token: t,
		Session: s,
	}

	return bot, err
}

func (b *DiscordBot) Run() error {
	session := b.Session

	// Register the messageCreate func as a callback for MessageCreate events.
	session.AddHandler(messageCreate)

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
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := convertToCoreMessage(s, m)

	response := core.GenerateResponse(message)

	if response != "" {
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

func convertToCoreMessage(s *discordgo.Session, dMsg *discordgo.MessageCreate) *message.Message{
	sender := &user.User{ // Call Read Function later
		ID: uuid.NewV4(),
		JBID: "",
		DiscordID: dMsg.Author.ID,
		DiscordUserName: dMsg.Author.Username,
		IrcID: "",
		TwitchID: "",
		TelegramID: "",
		DiscordUser: dMsg.Author,
	}

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
