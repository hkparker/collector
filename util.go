package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
)

func connectToWave() (ws net.Conn) {
	return
}

func discardUntil(done chan bool, channel chan interface{}) {

}

func rateLimit(frame Wireless80211Frame) bool {
	if frame.Type == "MgmtBeacon" {
		return true // only send an exact match for some beacon properties every second
	}
	return false
}

func streamFrames(frames chan Wireless80211Frame, endpoint string, client http.Client) {
	// Use gorilla client https://github.com/gorilla/websocket/blob/master/client.go
	origin := endpoint // wont work but gorilla will only need endpoint
	ws, err := websocket.Dial(endpoint, "", origin)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to dial wave")
		return // sleep, discard until redail success
	}
	for {
		frame := <-frames
		flat, err := json.Marshal(frame)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warn("failed to marshal frame")
			continue
		}
		if _, err := ws.Write([]byte(flat)); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("failed to send frame")
			// discard until rebuilt
		}
	}
}
