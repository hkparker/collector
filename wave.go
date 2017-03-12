package main

import (
	"crypto/tls"
	"crypto/x509"
	log "github.com/Sirupsen/logrus"
	"github.com/hkparker/Wave/models"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net"
	"net/url"
	"time"
)

func dialWave(wave_host string, frames chan models.Wireless80211Frame) net.Conn {
	log.WithFields(log.Fields{
		"at":        "dialWave",
		"wave_host": wave_host,
	}).Info("dialing Wave")

	done := make(chan bool)
	go discardUntil(done, frames)

	endpoint, err := url.Parse(endpoint_uri(wave_host))
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "dialWave",
			"error": err.Error(),
		}).Fatal("unable to parse wave uri")
	}
	origin, err := url.Parse(origin_uri(wave_host))
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "dialWave",
			"error": err.Error(),
		}).Fatal("unable to parse wave uri")
	}

	var ws net.Conn
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "dialWave",
			"error": err.Error(),
		}).Fatal("failed to load client certificate")
	}
	for {
		tls_config := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		if ca != "" {
			wave_pool := x509.NewCertPool()
			wave_ca, err := ioutil.ReadFile(ca)
			if err != nil {
				log.WithFields(log.Fields{
					"at":      "dialWave",
					"ca_file": ca,
				}).Fatal("Could not load wave ca")
			}
			wave_pool.AppendCertsFromPEM(wave_ca)
			tls_config.RootCAs = wave_pool
		}
		config := websocket.Config{
			Location:  endpoint,
			Origin:    origin,
			TlsConfig: tls_config,
			Version:   13,
		}
		var err error
		ws, err = websocket.DialConfig(&config)
		if err != nil {
			log.WithFields(log.Fields{
				"at":        "dialWave",
				"error":     err,
				"wave_host": wave_host,
			}).Error("failed to dial wave")
			time.Sleep(5 * time.Second)
		} else {
			log.WithFields(log.Fields{
				"at":        "dialWave",
				"wave_host": wave_host,
			}).Info("success dialing Wave, sending frames")
			done <- true
			break
		}
	}
	return ws

}
