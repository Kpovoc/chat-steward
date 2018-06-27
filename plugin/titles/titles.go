package titles

import (
	"gitlab.com/Kpovoc/chat-steward/core/message"
	"time"
	"fmt"
)

type TitleSuggestion struct{
	ID int
	Title string
	Author string
	CreatedOn time.Time
	TotalVotes int
}

var titles []*TitleSuggestion

var showTitle string

// Main plugin call
func Plugin(m *message.Message, msgContent string) string {
	author := m.Sender.DiscordUserName
	if "" == author {
		author = m.Sender.IrcID
	}

	t := &TitleSuggestion{
		ID: len(titles),
		Title: msgContent,
		Author: author,
		CreatedOn: m.CreatedOn,
		TotalVotes: 0,
	}
	t.logToConsole()
	titles = append(titles, t)
	return ""
}

func (t *TitleSuggestion) logToConsole() {
	fmt.Printf(
		"TitleSuggestion {" +
			"\n  ID: %d," +
			"\n  Title: %s," +
			"\n  Author: %s," +
			"\n  CreatedOn: %s," +
			"\n  TotalVotes: %d," +
			"\n}\n",
		t.ID, t.Title, t.Author, t.CreatedOn, t.TotalVotes)
}

func GetTitles() []*TitleSuggestion {
	return titles
}

func addVoteToTitle(id int) {
	titles[id].TotalVotes += 1
}

func StartShow(st string) {
	showTitle = st
	titles = nil
	ipsThatVoted = make(map[string]bool)
}