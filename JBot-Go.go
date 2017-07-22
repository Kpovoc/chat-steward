package main

import (
	"net"
	"log"
	"bufio"
	"net/textproto"
	"fmt"
)

type Bot struct{
	server string
	port string
	nick string
	user string
	channel string
	pass string
	pread, pwrite chan string
	conn net.Conn
	ioReader *bufio.Reader
	tpReader *textproto.Reader
}

func NewJbIrcBot() *Bot {
	return &Bot{
		server: "irc.geekshed.net",
		port: "6667",
		nick: "Japawig",
		user: "Japawig",
		channel: "#jupiterbroadcasting",
		pass: "",
		conn: nil,
		ioReader: nil,
		tpReader: nil,
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
		}
		fmt.Printf("%s\n", line)
	}
	return err
}

func main() {
	japawig := NewJbIrcBot()
	err := japawig.Launch()
	if err != nil {
		log.Printf("Exited with error: %s\n", err)
	} else {
		log.Println("Exited with no problems.")
	}
}
