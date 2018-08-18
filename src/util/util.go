package util

import "log"

func ExitWithError(err error) {
  ExitWithErrorMsg(err.Error())
}
func ExitWithErrorMsg(errMsg string) {
  log.Fatalf("An error has occurred:\n  %s\n  Application will now exit.", errMsg)
}

func IsStringInArray(str string, strs []string) bool {
  for i:=0;i<len(strs);i++ {
    if str == strs[i] {
      return true
    }
  }
  return false
}
