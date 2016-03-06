package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
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

func streamFrames(frames chan Wireless80211Frame, endpoint, origin string, client http.Client) {
	ws, err := websocket.Dial(endpoint, "", origin)
	if err != nil {
		fmt.Println("failed ws dail: ", err)
		return // sleep, discard until redail success
	}
	for {
		frame := <-frames
		flat, err := json.Marshal(frame)
		if err != nil {
			log.Println(err)
		}
		if _, err := ws.Write([]byte(flat)); err != nil {
			log.Println(err) // discard and rebuild
		}
	}
}
