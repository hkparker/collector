package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net"
)

func dialWave(wave_host string) net.Conn {
	// discard frames while dialing every 5 seconds until success
	endpoint := endpoint_uri(wave_host)
	origin := origin_uri(wave_host)
	log.WithFields(log.Fields{
		"wave_host": wave_host,
	}).Info("dialing Wave")
	ws, err := websocket.Dial(endpoint, "", origin)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"wave_host": wave_host,
		}).Error("failed to dial wave")
		//discardUntil(, frames)
	} else {
		log.WithFields(log.Fields{
			"wave_host": wave_host,
		}).Info("success dialing Wave, sending frames")
	}
	return ws

}

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

func streamFrames(frames chan Wireless80211Frame, wave_host string) {
	for {
		ws := dialWave(wave_host)
		for frame := range frames {
			flat, err := json.Marshal(frame)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Warn("failed to marshal frame")
				continue
			}
			if rateLimit(frame) {
				continue
			}
			if _, err := ws.Write([]byte(flat)); err != nil {
				ws.Close()
				log.WithFields(log.Fields{
					"error":     err,
					"wave_host": wave_host,
				}).Error("failed to send frame, redailing Wave")
				break
			}
		}
	}
}
