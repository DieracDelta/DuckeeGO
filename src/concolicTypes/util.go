package concolicTypes

import "log"

func reportError(message string, a ...interface{}) {
  log.printf(message, a...)
}
