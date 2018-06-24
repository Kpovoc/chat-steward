package ircbot

import (
	"net"
	"bufio"
	"net/textproto"
	"log"
	"fmt"
	"gitlab.com/Kpovoc/JBot-Go/core"
	"gitlab.com/Kpovoc/JBot-Go/core/message"
	"github.com/satori/go.uuid"
	"time"
	"gitlab.com/Kpovoc/JBot-Go/core/user"
)

type IrcConf struct {
	Servers []IrcServerConf
}

type IrcServerConf struct {
	Server string
	Port string
	Nick string
	User string
	Channels []string
}

type Bot struct{
	server string
	port string
	nickname string
	username string
	password string
	ircChannels []string
	connection net.Conn
	ioReader *bufio.Reader
	tpReader *textproto.Reader
	errChan chan error
	// TODO: Don't think we'll need these. If not, remove them.
	//pread, pwrite chan string
	//readNick, readText chan []byte
}

func Start(conf IrcConf, fatalErr chan error) {
	for i := 0; i < len(conf.Servers); i++ {
		go launchBotForServer(conf.Servers[i], fatalErr)
	}
}

func launchBotForServer(conf IrcServerConf, fatalErr chan error) {
	// TODO: Error handling for missing info
	bot := Bot {
		server: conf.Server,
		port: conf.Port,
		nickname: conf.Nick,
		username: conf.User,
		password: "",
		ircChannels: conf.Channels,
		connection: nil,
		ioReader: nil,
		tpReader: nil,
		errChan: fatalErr,
	}

	bot.launch()
}

func (bot *Bot) launch() {
	err := bot.connect()
	if err != nil {
		bot.errChan <- err
		return
	}

	bot.handleCredentials()

	err = bot.handleInitialPing()
	if err != nil {
		bot.errChan <- err
		return
	}

	bot.joinChannels()

	bot.listen()
}

func (bot *Bot) connect() (err error) {
	conn, err := net.Dial("tcp", bot.server + ":" + bot.port)
	if err != nil {
		log.Fatal("Unable to connect ", err)
		return err
	}
	bot.connection = conn
	log.Printf("Connected to Server %s (%s)\n", bot.server, bot.connection.RemoteAddr())
	return nil
}

func (bot *Bot) handleCredentials() {
	fmt.Fprintf(bot.connection, "NICK %s\r\n", bot.nickname)
	fmt.Fprintf(bot.connection, "USER %s 0 * :%s\r\n", bot.username, bot.username)
}

func (bot *Bot) handleInitialPing() (err error) {
	bot.ioReader = bufio.NewReader(bot.connection)
	bot.tpReader = textproto.NewReader(bot.ioReader)
	var line string

	// Wait for PING
	for line, err = bot.tpReader.ReadLine();
		line[0:4] != "PING" && err == nil;
	line, err = bot.tpReader.ReadLine() {
		fmt.Printf("%s\n", line)
	}

	if err != nil {
		bot.connection.Close()
		return err
	}

	fmt.Printf("%s\r\n", line)
	var response = "PONG" + line[4:]
	fmt.Fprintf(bot.connection, "%s\r\n", response)
	return nil
}

func (bot *Bot) joinChannels() {
	for i := 0; i < len(bot.ircChannels); i++ {
		fmt.Fprintf(bot.connection, "JOIN %s\r\n", bot.ircChannels[i])
	}
}

func (bot *Bot) listen() {
	// Listen for channel output, and handle it
	defer bot.connection.Close()
	for {
		line, err := bot.tpReader.ReadLine()
		if err != nil {
			break;
		}

		if line[0:4] == "PING" {
			var response = "PONG" + line[4:]
			fmt.Fprintf(bot.connection, "%s\r\n", response)
			continue;
		}

		bot.handleLine(line)
	}
}

func (bot *Bot) handleLine(line string) {
	// & = Operator, Admin
	// @ = Operator
	// ~ = Owner, Operator
	// + = Has Voice
	nick, msg, channel := parseChannelLine(line)
	if "" == nick || 		// Line is not a message
		bot.nickname == nick {	// Line is a message from our bot
		// Ignore and move on
		return;
	}

	fmt.Printf("%s %s: %s\n",channel, nick, msg)

	coreMsg := convertToCoreMessage(nick, msg)

	response := core.GenerateResponse(coreMsg)

	if "" != response {
		fmt.Printf("%s %s: %s\n",channel, bot.nickname, response)
		fmt.Fprintf(bot.connection, "PRIVMSG %s :%s\r\n", channel, response)
	}
}

func parseChannelLine(line string) (nick string, msg string, channel string) {
	// Message we care about looks like the following:
	// :nick!user@host PRIVMSG #channel :msg
	var nickB []byte
	var typeB []byte
	var chanB []byte
	var mesgB []byte

	i := 1 // skipping the very first colon

	for ; i < len(line); i++ {
		c := line[i]

		if ' ' == c {
			// Reached first space before nickBang. Not msg we're interested in.
			return "","",""
		}

		if '!' == c {
			// End of Nickname reached.
			i++
			break
		}

		nickB = append(nickB, c)
	}

	// Still here, means we got nickname. Skip forward until 1st space
	for ; i < len(line); i++ {
		c := line[i]

		if ' ' == c {
			// Reached first space. Next character will start message type.
			i++
			break
		}
	}

	// Continue until second space for message type
	for ; i < len(line); i++ {
		c := line[i]

		if ' ' == c {
			// Reached second space. Next character should be channelHash
			i++
			break
		}

		typeB = append(typeB, c)
	}

	msgType := string(typeB)
	if "PRIVMSG" != msgType {
		// Not a type we're interested in.
		return "","",""
	}

	// Still here, so got PRIVMSG. Now get channel the msg was sent on
	if '#' != line[i] {
		// Something went wrong. Should be channelHash. Ignore whatever this line is.
		return "","",""
	}

	for ; i < len(line); i++ {
		c := line[i]

		if ' ' == c {
			// Reached third space. Next character should be messageColon
			i++
			break
		}

		chanB = append(chanB, c)
	}

	// Still here, so got Channel. Now get message.
	if ':' != line[i] {
		// Something went wrong. Should be beginning of message. Ignore whatever this line is.
		return "","",""
	}

	i++ // Now at the start of the msg
	for ; i < len(line); i++ {
		mesgB = append(mesgB, line[i])
	}

	return string(nickB), string(mesgB), string(chanB)
}

func convertToCoreMessage(nick string, msg string) *message.Message {
	id,_ := uuid.NewV4()
	sender := &user.User{ // Call Read Function later
		ID: id,
		JBID: "",
		DiscordID: "",
		DiscordUserName: "",
		IrcID: nick,
		TwitchID: "",
		TelegramID: "",
		DiscordUser: nil,
	}

	return &message.Message {
		Sender: sender,
		CreatedOn: time.Now(),
		Content: msg,
	}
}