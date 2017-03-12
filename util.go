package main

import (
	"github.com/hkparker/Wave/models"
	"strconv"
	"time"
)

var lastBeacons map[string]time.Time

func discardUntil(done chan bool, channel chan models.Wireless80211Frame) {
	for {
		select {
		case _ = <-done:
			return
		case _ = <-channel:
		}
	}

}

func rateLimit(frame models.Wireless80211Frame) bool {
	if frame.Type == "MgmtBeacon" {
		// Can be used to filter beacon frames that are
		// not useful to reduce bandwidth requirements
		return false
	}
	return false
}

func origin_uri(wave_host string) string {
	return "https://" + wave_host + ":" + strconv.Itoa(port) + "/frames"
}

func endpoint_uri(wave_host string) string {
	return "wss://" + wave_host + ":" + strconv.Itoa(port) + "/frames"
}
