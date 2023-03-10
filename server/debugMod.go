package main

import "log"

func DebugErr(err error) bool {
	if err != nil {
		return true
	}
	return false
}
func DebugLog(a ...interface{}) {
	if cfg.DebugMode {
		log.Println(a)
	}
}
