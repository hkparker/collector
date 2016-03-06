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

`go build`

Usage
-----

```

```

License
-------

This project is licensed under the MIT license, see LICENSE for more information.
