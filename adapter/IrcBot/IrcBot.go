package IrcBot

import (
	"net"
	"bufio"
	"net/textproto"
	"log"
	"fmt"
	"math/rand"
)

type Bot struct{
	server string
	port string
	nick string
	user string
	channels []string
	channel string
	pass string
	pread, pwrite chan string
	readNick, readText chan []byte
	conn net.Conn
	ioReader *bufio.Reader
	tpReader *textproto.Reader
	errChan chan error
}

//func (bot *Bot) SetChannel(channelName string) {
//	bot.channel = channelName
//}

func (bot *Bot) Connect() (err error) {
	conn, err := net.Dial("tcp", bot.server + ":" + bot.port)
	if err != nil {
		log.Fatal("Unable to connect ", err)
		return err
	}
	bot.conn = conn
	log.Printf("Connected to Server %s (%s)\n", bot.server, bot.conn.RemoteAddr())
	return nil
}

func (bot *Bot) JoinChannel() (err error) {
	fmt.Fprintf(bot.conn, "NICK %s\r\n", bot.nick)
	fmt.Fprintf(bot.conn, "USER %s 0 * :%s\r\n", bot.user, bot.user)
	bot.ioReader = bufio.NewReader(bot.conn)
	bot.tpReader = textproto.NewReader(bot.ioReader)
	var line string

	// Wait for PING
	for line, err = bot.tpReader.ReadLine();
		line[0:4] != "PING" && err == nil;
		line, err = bot.tpReader.ReadLine() {
		fmt.Printf("%s\n", line)
	}

	if err != nil {
		bot.conn.Close()
		return err
	}
	fmt.Printf("%s\r\n", line)
	var response = "PONG" + line[4:]
	fmt.Fprintf(bot.conn, "%s\r\n", response)
	fmt.Fprintf(bot.conn, "JOIN %s\r\n", bot.channel)
	return nil
}

func (bot *Bot) ListenForChannelOutput() (err error) {
	defer bot.conn.Close()
	for {
		line, err := bot.tpReader.ReadLine()
		if err != nil {
			break
		}
		if line[0:4] == "PING" {
			var response = "PONG" + line[4:]
			fmt.Fprintf(bot.conn, "%s\r\n", response)
		} else {
			bot.ParseChannelLine(line)
		}
	}
	return err
}

func (bot *Bot) ParseChannelLine(line string) {
	var infoB []byte
	var textB []byte
	var nickB []byte

	infoColonPassed := false
	nickBangReached := false
	firstSpaceReached := false

	for i := 1; i < len(line); i++ {
		c := line[i]
		if !infoColonPassed {
			if c == ':' {
				infoColonPassed = true
				continue
			} else if !nickBangReached {
				if !firstSpaceReached {
					if c == ' ' {
						firstSpaceReached = true
						continue
					} else if c == '!' {
						nickBangReached = true
						continue
					} else {
						nickB = append(nickB, c)
					}
				} else {
					nickB = append(nickB, c)
				}
			} else {
				infoB = append(infoB, c)
			}
		} else {
			textB = append(textB, c)
		}
	}
	if nickBangReached {
		fmt.Printf("%s: %s\n", nickB, textB)
		bot.checkForBang(nickB, textB)
	}
}

func (bot *Bot) checkForBang(nick []byte, text []byte) {
	if text[0] == '!' {
		var commandB []byte

		for i := 1; i < len(text) && text[i] != ' '; i++ {
			commandB = append(commandB, text[i])
		}

		command := string(commandB)

		switch(command) {
		case "8ball":
			eightBallModule(bot)
		default:
			fmt.Printf("< User %s invoked unknown \"%s\" function >\n", nick, command)
		}
	}
}

func (bot *Bot) Launch() (err error) {
	err = bot.Connect()
	if err != nil {
		return err
	}

	err = bot.JoinChannel()
	if err != nil {
		return err
	}

	err = bot.ListenForChannelOutput()
	return err
}

func eightBallModule(bot *Bot) {
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes, definitely",
		"You may rely on it",
		"As I see it, yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy, try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"Outlook not so good",
		"Very doubtful",
	}
	response := answers[rand.Intn(len(answers))]
	fmt.Printf("Response was \"%s\"\n",response)
	// fmt.Fprintf(bot.conn, "PRIVMSG %s :%s\r\n", bot.channel, response)
}