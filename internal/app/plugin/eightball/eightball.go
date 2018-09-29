package eightball

import "gitlab.com/Kpovoc/chat-steward/internal/app/core/response"

type RandInt func(int) int

const noQuestionResponse = "You must first ask a question, before you can receive the answer."

func Plugin(msgContent string, randFn RandInt) *response.Response {
  if len(msgContent) <= 0 {
    return &response.Response {
      Content: noQuestionResponse,
      ForceWhisper: false,
    }
  }

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

  selectedIndex := randFn(len(answers))

  return &response.Response {
    Content: answers[selectedIndex],
    ForceWhisper: false,
  }
}