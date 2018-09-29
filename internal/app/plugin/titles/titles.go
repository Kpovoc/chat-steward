package titles

import (
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/message"
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/response"
  "time"
  "fmt"
  "gitlab.com/Kpovoc/chat-steward/internal/app/sqlite"
  "log"
)

type TitleSuggestion struct{
  ID int
  Title string
  Author string
  CreatedOn time.Time
  TotalVotes int
}

var nextId int

func Init() {
  db := sqlite.GetInstance()
  rows, err := db.Query("SELECT count(*) FROM titles")
  if err != nil {
    log.Println("Error getting initial count of titles table: " + err.Error())
    nextId = 0
    return
  }

  defer rows.Close()
  if rows.Next() {
    err = rows.Scan(&nextId)
    if err != nil {
      log.Println("Error getting initial count of titles table: " + err.Error())
      nextId = 0
      return
    }
    return
  }

  nextId = 0
}

// Main plugin call
func Plugin(m *message.Message, msgContent string) *response.Response {
  author := m.Sender.DiscordUserName
  if "" == author {
    author = m.Sender.IrcID
  }

  t := &TitleSuggestion{
    ID: nextId,
    Title: msgContent,
    Author: author,
    CreatedOn: m.CreatedOn,
    TotalVotes: 0,
  }
  t.logToConsole()
  err := upsertTitle(t)
  if err != nil {
    log.Println("Failed to add title: " + err.Error())
    return &response.Response {
      Content: "Failed to add title.",
      ForceWhisper: true,
    }
  }
  nextId += 1
  return &response.Response {
    Content: "\"" + msgContent + "\" was added to the title suggestions.",
    ForceWhisper: true,
  }

}

func GetTitles() []*TitleSuggestion {
  var titles []*TitleSuggestion
  db := sqlite.GetInstance()
  rows, err := db.Query("SELECT * FROM titles ORDER BY total_votes DESC")
  defer rows.Close()
  if err != nil { return nil }

  for rows.Next() {
    t := TitleSuggestion{}
    err = rows.Scan(&t.ID, &t.Title, &t.Author, &t.CreatedOn, &t.TotalVotes)
    if err != nil { return nil }

    titles = append(titles, &t)
  }

  return titles
}

// Show plugin integration
func resetTitles() {
  nextId = 0
  ipsThatVoted = make(map[string]bool)
  deleteTitles()
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

func readTitle(id int) *TitleSuggestion {
  db := sqlite.GetInstance()
  rows, err := db.Query("SELECT * FROM titles WHERE id=?", id)
  defer rows.Close()
  if err != nil { return nil }
  for rows.Next() {
    t := &TitleSuggestion{}
    err = rows.Scan(&t.ID, &t.Title, &t.Author, &t.CreatedOn, &t.TotalVotes)
    if err != nil { return nil }
    // Should just be one
    return t
  }
  return nil
}

func upsertTitle(t *TitleSuggestion) error {
  query :=
      "INSERT OR REPLACE INTO titles(" +
          "id, title, author, created_on, total_votes" +
          ") values(?, ?, ?, ?, ?)"

  db := sqlite.GetInstance()
  _, err := db.Exec(query, t.ID, t.Title, t.Author, t.CreatedOn, t.TotalVotes)
  if err != nil { return err }
  return nil
}

func deleteTitles() {
  db := sqlite.GetInstance()
  db.Exec("delete from titles")
}

func addVoteToTitle(id int) {
  t := readTitle(id)
  if t == nil { return }
  t.TotalVotes += 1
  err := upsertTitle(t)
  if err != nil {
    log.Println("Error adding vote to title: " + err.Error())
  }
}
