package dii2p

import (
    "log"
)

var Verbose bool = false

func Log(msg ...string) {
	if Verbose {
		log.Println("LOG: ", msg)
	}
}

func LogA(msg []string) {
	if Verbose {
		log.Println("LOG: ", msg)
	}
}

func Warn(err error, errmsg string, msg ...string) (bool, error) {
	LogA(msg)
	if err != nil {
		log.Println("WARN: ", errmsg, err)
		return false, nil
	}
	return true, nil
}

func Fatal(err error, errmsg string, msg ...string) (bool, error) {
	LogA(msg)
	if err != nil {
		log.Fatal("FATAL: ", errmsg, err)
		return false, err
	}
	return true, nil
}
