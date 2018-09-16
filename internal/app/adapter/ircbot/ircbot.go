package ircbot

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"time"

	"gitlab.com/Kpovoc/chat-steward/internal/app/core"
	"gitlab.com/Kpovoc/chat-steward/internal/app/core/message"
	"gitlab.com/Kpovoc/chat-steward/internal/app/core/user"

	)

type IrcConf struct {
	Servers []IrcServerConf
}

type IrcServerConf struct {
	Server string
	Port string
	Nick string
	User string
	Password string
	Channels []string
	AdminNicks []string
}

type Bot struct{
	server string
	port string
	nickname string
	username string
	password string
	ircChannels []string
	adminNicks []string
	connection net.Conn
	ioReader *bufio.Reader
	tpReader *textproto.Reader
	errChan chan error
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
		password: conf.Password,
		ircChannels: conf.Channels,
		adminNicks: conf.AdminNicks,
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

	// Handle Credentials
	fmt.Fprintf(bot.connection, "NICK %s\r\n", bot.nickname)
	fmt.Fprintf(bot.connection, "USER %s 0 * :%s\r\n", bot.username, bot.username)

	// Handle first ping
	err = bot.handleInitialPing()
	if err != nil {
		return err
	}

	// Join Channels
	for i := 0; i < len(bot.ircChannels); i++ {
		fmt.Fprintf(bot.connection, "JOIN %s\r\n", bot.ircChannels[i])
	}

	// Identify with the NickServ
	fmt.Fprintf(bot.connection, "PRIVMSG NickServ :IDENTIFY %s\r\n", bot.password)

	return nil
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
	nick, msg, channel := parseChannelLine(line)
	if "" == nick || 		// Line is not a message
		bot.nickname == nick {	// Line is a message from our bot
		// Ignore and move on
		return;
	}

	fmt.Printf("%s <- %s: %s\n",channel, nick, msg)

	coreMsg := bot.convertToCoreMessage(nick, msg)

	response := core.GenerateResponse(coreMsg)
	if "" != response {
		responseChan := channel
		if bot.nickname == channel {
			responseChan = nick
		}

		fmt.Printf("%s <- %s: %s\n", responseChan, bot.nickname, response)
		fmt.Fprintf(bot.connection, "PRIVMSG %s :%s\r\n", responseChan, response)
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
			// Reached second space. Next character should be start of channel/user
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

	//if '#' != line[i] {
	//	// Something went wrong. Should be channelHash. Ignore whatever this line is.
	//	return "","",""
	//}

	// Still here, so got PRIVMSG. Now get channel the msg was sent on
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

func (bot *Bot) convertToCoreMessage(nick string, msg string) *message.Message {
	sender := user.New( // Call Read Function later
	"",
	"",
	"",
	nick,
	"",
	"",
	bot.isAdmin(nick))

	return &message.Message {
		Sender: sender,
		CreatedOn: time.Now(),
		Content: msg,
	}
}

func (bot *Bot) isAdmin(nick string) bool {
	for i:=0; i<len(bot.adminNicks); i++ {
		if nick == bot.adminNicks[i] {
			return true
		}
	}

	return false
}
