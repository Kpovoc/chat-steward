package show

import (
  "testing"
)

// TODO: '!show start <show_title>' should do the following:
// if user is not admin
// - return something that informs the room the user does not have permission to run start command
// if user is admin and arg is invalid or empty
// - return something that informs the room the show title is invalid
// if user is admin and arg is valid
// - set the title plugin's showTitle to 'arg'
// - erase the title plugin's current suggestions
// * return a message stating that 'arg' show has started

func TestShowCmd_Start(t *testing.T) {
  t.Run("returns show started message", func(t *testing.T) {
    title := "Linux Unplugged"
    spy := &TitleFuncSpy{}

    want := title + StartShowSuffix
    got := StartShow(title, spy)

    if want != got {
      t.Errorf("want '%s', got '%s'", want, got)
    }
  })

  t.Run("calls StartTitle no error", func(t *testing.T) {
    title := "Linux Unplugged"
    spy := &TitleFuncSpy{}

    StartShow(title, spy)

    want := true
    got := spy.startTitleCalled

    if want != got {
      t.Errorf("failed to call StartTitle")
    }
  })
}

type TitleFuncSpy struct {
  startTitleCalled bool
}

func (t *TitleFuncSpy) StartTitle(title string) error {
  t.startTitleCalled = true
  return nil
}
