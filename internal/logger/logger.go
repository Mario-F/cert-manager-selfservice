package logger

import (
	"log"
)

var (
	DebugMode   bool = false
	VerboseMode bool = false
)

// Debugf logs formatted debug messages
func Debugf(msg string, vars ...interface{}) {
	if DebugMode {
		log.Printf(msg, vars...)
	}
}

// Verbosef logs formatted info messages
func Verbosef(msg string, vars ...interface{}) {
	if VerboseMode {
		log.Printf(msg, vars...)
	}
}

// Infof logs formatted info messages
func Infof(msg string, vars ...interface{}) {
	log.Printf(msg, vars...)
}