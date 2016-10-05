collector
=========

Collector is part of the [Wave](https://github.com/hkparker/Wave) wireless IDS and visualizer.  It listens for 802.11 frames on wireless interfaces and streams the JSON formatted frames to Wave via websocket.

Installing
----------

```
go get github.com/hkparker/collector
```

Usage
-----

```
$ collector -help
Usage of ./collector:
  -certificate string
        path to a TLS client certificate to present to Wave (default "collector.pem")
  -interfaces string
        comma-separated list of network interfaces to collect (default "mon0")
  -local
        collect frames without streaming them to wave (use with -print)
  -port int
        port the Wave server is accessible on (default 443)
  -print
        print the frames to standard output
  -wave string
        hostname of Wave server to stream frames to (default "127.0.0.1")
```
