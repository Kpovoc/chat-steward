package titles

import (
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/response"
  "gitlab.com/Kpovoc/chat-steward/internal/app/sqlite"
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/message"
  "log"
)

var showName string

func StartShow(m *message.Message, msgContent string) *response.Response {
  author := m.Sender.DiscordUserName
  if "" == author {
    author = m.Sender.IrcID
  }

  if !m.Sender.IsBotAdmin {
    return &response.Response {
      Content: author + " is not a registered Bot Admin. Could not comply with request.",
      ForceWhisper: false,
    }
  }

  updateShowName(msgContent)
  resetTitles()
  return &response.Response {
    Content: "Started \"" + msgContent + "\"",
    ForceWhisper: false,
  }
}

func getShowName() string {
  if showName == "" {
    showName = readShowName()
  }
  return showName
}

func readShowName() string {
  db := sqlite.GetInstance()
  rows, err := db.Query("SELECT * FROM show")
  if err != nil {
    log.Println("Failed to fetch show name: " + err.Error())
    return ""
  }

  defer rows.Close()
  for rows.Next() {
    name := ""
    err = rows.Scan(&name)
    if err != nil {
      log.Println("Failed to fetch show name: " + err.Error())
      return ""
    }
    return name
  }

  return ""
}

func updateShowName(sn string) {
  db := sqlite.GetInstance()
  tx, err := db.Begin()
  if err != nil {
    log.Println( "Failed to update show name: " + err.Error())
    return
  }

  defer tx.Rollback()

  delStmt, err := tx.Prepare("DELETE FROM show")
  if err != nil {
    log.Println( "Failed to update show name: " + err.Error())
    return
  }

  defer delStmt.Close()

  _, err = delStmt.Exec()
  if err != nil {
    log.Println( "Failed to update show name: " + err.Error())
    return
  }

  insStmt, err := tx.Prepare("INSERT INTO show VALUES (?)")
  if err != nil {
    log.Println( "Failed to update show name: " + err.Error())
    return
  }

  defer insStmt.Close()

  _, err = insStmt.Exec(sn)
  if err != nil {
    log.Println( "Failed to update show name: " + err.Error())
    return
  }

  err = tx.Commit()
  if err != nil {
    log.Println( "Failed to update show name: " + err.Error())
    return
  }
  showName = sn
}
