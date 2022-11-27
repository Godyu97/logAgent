package main

import "log"

func DebugLog(a ...interface{}) {
	if cfg.DebugMode {
		log.Println(a)
	}
}
