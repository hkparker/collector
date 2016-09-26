package main

import (
	"time"
)

var lastBeacons map[string]time.Time

func discardUntil(done chan bool, channel chan Wireless80211Frame) {
	for {
		select {
		case _ = <-done:
			return
		case _ = <-channel:
		}
	}

}

func rateLimit(frame Wireless80211Frame) bool {
	if frame.Type == "MgmtBeacon" {
		return true // only send an exact match for some beacon properties every second
	}
	return false
}

func origin_uri(wave_host string) string {
	return ""
}

func endpoint_uri(wave_host string) string {
	return ""
}
