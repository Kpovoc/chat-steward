package discordbot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/Kpovoc/JBot-Go/core/user"
	"github.com/satori/go.uuid"
	"github.com/Kpovoc/JBot-Go/core"
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

	message := convertToCoreMessage(m)

	response := core.GenerateResponse(message)

	if response != "" {
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

func convertToCoreMessage(dMsg *discordgo.MessageCreate) *core.Message{
	sender := &user.User{ // Call Read Function later
		ID: uuid.NewV4(),
		JBID: "",
		DiscordID: dMsg.Author.ID,
		IrcID: "",
		TwitchID: "",
		TelegramID: "",
		DiscordUser: dMsg.Author,
	}

	content := dMsg.Content

	return &core.Message{
		Sender: sender,
		Content: content,
	}
}
