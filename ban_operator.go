package main

import "time"

var banList = make(map[string]time.Time, 100)

const banDuration = 5 //minutes

func isHostBanned(host string) bool {
	expireTime, ok := banList[host]
	if !ok {
		return false
	}
	if time.Now().After(expireTime) {
		delete(banList, host)
		return false
	}
	return true
}

func banHost(host string) {
	banList[host] = time.Now().Add(time.Minute * banDuration)
}

func unbanExpired() {
	for host, expireTime := range banList {
		if time.Now().After(expireTime) {
			delete(banList, host)
		}
	}
}
