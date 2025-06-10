package pkg

import "log"

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
)

var logger = log.Default()

func Default(msg string) {
	logger.Println(msg)
}

func Error(msg string) {
	logger.Println(Red + msg + Reset)
}
