package main

import (
	"flag"
	"strings"
)

var print bool
var port int
var local bool

func main() {
	var interface_list string
	var host string
	var certificate string
	flag.StringVar(&interface_list, "interfaces", "mon0", "comma-separated list of network interfaces to collect")
	flag.StringVar(&host, "wave", "127.0.0.1", "hostname of Wave server to stream frames to")
	flag.IntVar(&port, "port", 443, "port the Wave server is accessible on")
	flag.StringVar(&certificate, "certificate", "collector.pem", "path to a TLS client certificate to present to Wave")
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
