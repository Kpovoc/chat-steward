package suggesttitle

import (
	"github.com/Kpovoc/JBot-Go/core/message"
	"time"
	"fmt"
)

type TitleSuggestion struct{
	Title string
	Author string
	CreatedOn time.Time
	TotalVotes int
}

var titles []*TitleSuggestion

// Main plugin call
func Plugin(m *message.Message, msgContent string) string {
	t := &TitleSuggestion{
		Title: msgContent,
		Author: m.Sender.DiscordUserName,
		CreatedOn: m.CreatedOn,
		TotalVotes: 0,
	}
	t.logToConsole()
	titles = append(titles, t)
	return ""
}

func (t *TitleSuggestion) logToConsole() {
	fmt.Printf("TitleSuggestion {\n  Title: %s,\n  Author: %s,\n  CreatedOn: %s,\n  TotalVotes: %d,\n}\n",
		t.Title, t.Author, t.CreatedOn.String(), t.TotalVotes)
}