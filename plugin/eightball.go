package plugin

import (
	"math/rand"
)

func EightBall(msgContent string) string {
	if len(msgContent) <= 0 {
		return "You must first ask a question, before you can receive the answer."
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

	selectedIndex := rand.Intn(len(answers))

	return answers[selectedIndex]
}