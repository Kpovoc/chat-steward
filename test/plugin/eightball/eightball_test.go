package eightball

import (
	"testing"
	"gitlab.com/Kpovoc/chat-steward/internal/app/plugin/eightball"
)

func TestPlugin_BlankMsg(t *testing.T) {
	exptd := "You must first ask a question, before you can receive the answer."

	result := eightball.Plugin("", testRand)
	if result != exptd { t.Fail() }
}

func TestPlugin_NonBlankMsg(t *testing.T) {
	exptd := "It is certain"

	result := eightball.Plugin("Will this test pass?", testRand)
	if result != exptd { t.Fail() }
}

func testRand(x int) int {
	return 0
}
