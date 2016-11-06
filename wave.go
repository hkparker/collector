package main

import (
	"crypto/tls"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net"
	"net/url"
	"time"
)

func dialWave(wave_host string, frames chan Wireless80211Frame) net.Conn {
	log.WithFields(log.Fields{
		"wave_host": wave_host,
	}).Info("dialing Wave")

	done := make(chan bool)
	go discardUntil(done, frames)

	endpoint, err := url.Parse(endpoint_uri(wave_host))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("unable to parse wave uri")
	}
	origin, err := url.Parse(origin_uri(wave_host))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("unable to parse wave uri")
	}

	var ws net.Conn
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to load client certificate")
	}
	for {
		config := websocket.Config{
			Location: endpoint,
			Origin:   origin,
			TlsConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
		var err error
		ws, err = websocket.DialConfig(&config)
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
