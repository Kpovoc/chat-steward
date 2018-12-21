package show

import (
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/response"
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/user"
  "gitlab.com/Kpovoc/chat-steward/internal/app/plugin/titles"
  "strings"
)

const (
  StartShowSuffix = " starts now!"
  CmdStart = "start"
  CmdEnd = "end"
  CmdNext = "next"
  CmdNow = "now"
  RspHelp = "TODO: add help info"
  RspNotAllowed = " is not a registered Bot Admin. Could not comply with request."
)
type TitleFuncs interface {
  StartTitle(title string) error
}


func Plugin(u *user.User, msgContent string) *response.Response {
  cmd, content := parseCmd(msgContent)

  switch cmd {
  case CmdStart:
    return showResponse(StartShow(u, content, TempTitleFuncs{}))
  default:
    return showResponse("Unknown Command")
  }
}

func StartShow(u *user.User, title string, titleFuncs TitleFuncs) string {
  if !u.IsBotAdmin {
    return userNotAllowedRsp(u.GetDisplayName())
  }

  err := titleFuncs.StartTitle(title)
  if err != nil {
    return "Error starting show: " + err.Error()
  }

  return title + StartShowSuffix
}

type TempTitleFuncs struct {}

func (t TempTitleFuncs) StartTitle(title string) error {
  return titles.StartShow(title)
}

func parseCmd(str string) (command, arguments string) {
  if str == "" {
    return
  }

  // Trim begin and end spaces
  str = strings.TrimSpace(str)

  firstSpace := strings.Index(str, " ")

  // No extra variables. The str is the command
  if firstSpace == -1 {
    return str, ""
  }


  return str[:firstSpace], str[firstSpace+1:]
}

func userNotAllowedRsp(username string) string {
  return username + RspNotAllowed
}

func showResponse(msg string) *response.Response {
  return &response.Response {
    Content: msg,
    ForceWhisper: false,
  }
}
