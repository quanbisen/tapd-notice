package logger

import "log"

func Setup() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}
