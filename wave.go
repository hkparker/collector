package main

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net"
	"time"
)

func dialWave(wave_host string, frames chan Wireless80211Frame) net.Conn {
	log.WithFields(log.Fields{
		"wave_host": wave_host,
	}).Info("dialing Wave")

	done := make(chan bool)
	go discardUntil(done, frames)

	endpoint := endpoint_uri(wave_host)
	origin := origin_uri(wave_host)

	var ws net.Conn
	for {
		var err error
		ws, err = websocket.Dial(endpoint, "", origin)
		if err != nil {
			log.WithFields(log.Fields{
				"error":     err,
				"wave_host": wave_host,
			}).Error("failed to dial wave")
			time.Sleep(5 * time.Second)
		} else {
			log.WithFields(log.Fields{
				"wave_host": wave_host,
			}).Info("success dialing Wave, sending frames")
			done <- true
			break
		}
	}
	return ws

}
