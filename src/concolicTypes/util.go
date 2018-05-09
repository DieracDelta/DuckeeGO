package concolicTypes

import "log"

func reportError(message string, a ...interface{}) {
  log.Printf(message, a...)
}
