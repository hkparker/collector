package main

import (
	"flag"
	"strings"
)

var print bool
var port int
var local bool
var certificate string
var key string

func main() {
	var interface_list string
	var host string
	flag.StringVar(&interface_list, "interfaces", "mon0", "comma-separated list of network interfaces to collect")
	flag.StringVar(&host, "wave", "127.0.0.1", "hostname of Wave server to stream frames to")
	flag.IntVar(&port, "port", 444, "port the Wave server is accessible on")
	flag.StringVar(&certificate, "certificate", "collector.crt", "path to a TLS client certificate to present to Wave")
	flag.StringVar(&key, "key", "collector.key", "path to a TLS client certificate private key")
	flag.BoolVar(&print, "print", false, "print the frames to standard output")
	flag.BoolVar(&local, "local", false, "collect frames without streaming them to wave (use with -print)")
	flag.Parse()

	interfaces := strings.Split(interface_list, ",")
	frames := make(chan Wireless80211Frame, 100)
	go streamFrames(frames, host)
	for _, iface := range interfaces {
		go sniffInterface(iface, frames)
	}
	select {}
}
