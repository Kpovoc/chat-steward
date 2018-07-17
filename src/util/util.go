package util

import "log"

func ExitWithError(err error) {
	ExitWithErrorMsg(err.Error())
}
func ExitWithErrorMsg(errMsg string) {
	log.Fatalf("An error has occurred:\n  %s\n  Application will now exit.", errMsg)
}