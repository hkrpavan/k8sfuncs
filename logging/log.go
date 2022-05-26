package logging

import (
	"log"
	"os"
)

func getLogParams(debug bool) *int {
	var params int = 0
	if debug {
		params = log.Ldate | log.Ltime | log.Lshortfile
	}
	return &params
}

var (
	Debug         = true
	InfoLogger    = log.New(os.Stdout, "INFO: ", *getLogParams(Debug))
	WarningLogger = log.New(os.Stderr, "WARNING: ", *getLogParams(Debug))
	ErrorLogger   = log.New(os.Stderr, "ERROR: ", *getLogParams(Debug))
)
