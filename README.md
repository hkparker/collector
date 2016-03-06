collector
=========

collector is part of the [Wave](https://github.com/hkparker/Wave) wireless IDS and visualizer.  It listens for 802.11 frames on wireless interfaces and streams the JSON formatted frames to Wave via websocket.

Installing
----------

```
go get github.com/hkparker/collector
```

Building
--------

**Install deps:**

```
go get github.com/google/gopacket
go get github.com/google/gopacket/layers
go get github.com/google/gopacket/pcap
go get golang.org/x/net/websocket
```

**Build:**

```
go build
```

Usage
-----

```
Usage of ./collector:
  -certificate string
    	path to a TLS client certificate to present to Wave (default "collector.pem")
  -interfaces string
    	comma-separated list of network interfaces to collect (default "mon0")
  -port int
    	port the Wave server is accessible on (default 443)
  -wave string
    	hostname of Wave server to stream frames to (default "127.0.0.1")
```

License
-------

This project is licensed under the MIT license, see LICENSE for more information.
