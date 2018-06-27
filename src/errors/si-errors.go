package dii2perrs

import (
	"log"
)

//Verbose determines the log verbosity
var Verbose = false

//Log wraps logging
func Log(msg ...interface{}) {
	if Verbose {
		log.Println("LOG: ", msg)
	}
}

//LogA wraps logging arrays
func LogA(msg []interface{}) {
	if Verbose {
		log.Println("LOG: ", msg)
	}
}

//Warn checks for non-fatal errors and re-sets them
func Warn(err error, errmsg interface{}, msg ...interface{}) (bool, error) {
	LogA(msg)
	if err != nil {
		log.Println("WARN: ", errmsg, err)
		return false, nil
	}
	return true, nil
}

//Fatal prints results of fatal errors and exits
func Fatal(err error, errmsg interface{}, msg ...interface{}) (bool, error) {
	LogA(msg)
	if err != nil {
		log.Fatal("FATAL: ", errmsg, err)
		return false, err
	}
	return true, nil
}
