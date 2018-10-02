package info

import (
  "fmt"
  "gitlab.com/Kpovoc/chat-steward/internal/app/core/response"
)

const (
  InfoNotFoundFormat = "Info could not be found for '%s'."
  InfoHelpFormat = "To obtain information, use !info followed by one of the following: %s"
)

var (
  infoMap = map[string]string{}
  validTerms = ""
  infoHelpStr = ""
)

type InfoConfItem struct {
  Term string
  Info string
}

func Init(confItems []InfoConfItem) {
  for i:=0;i<len(confItems);i++ {
    confItem := confItems[i]

    // TODO: add check to ensure Term is alphabetic, lower-case, and contains no spaces
    // Skip if term or info is blank
    if confItem.Term == "" || confItem.Info == "" {
      continue
    }

    infoMap[confItem.Term] = confItem.Info
    validTerms += " " + confItem.Term
  }

  infoHelpStr = fmt.Sprintf(InfoHelpFormat, validTerms)
}

func Plugin(msgContent string) *response.Response {
  resp := &response.Response {
    Content: infoMap[msgContent],
    ForceWhisper: true,
  }
  if "" == resp.Content {
    resp.Content = fmt.Sprintf(InfoNotFoundFormat + " " + infoHelpStr, msgContent)
  }

  return resp
}
