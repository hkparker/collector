package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net/http"
)

func PostFrames(frames chan Wireless80211Frame, endpoint string, client http.Client) {
	for {
		frame := <-frames
		flat, _ := json.Marshal(frame)
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(flat))
		if err != nil {
			log.Println(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		resp.Body.Close()
	}
}

func RateLimitFrame(frame Wireless80211Frame) bool {
	return false
}

func main() {
	var iface = flag.String("interface", "mon0", "interface to sniff")
	var endpoint = flag.String("endpoint", "http://127.0.0.1:4567/frames", "server to post packet data to")
	flag.Parse()
	frames := make(chan Wireless80211Frame, 100)
	go PostFrames(frames, *endpoint, http.Client{})

	if handle, err := pcap.OpenLive(*iface, 1600, true, 1); err != nil {
		log.Fatal(err)
	} else if err := handle.SetBPFFilter(""); err != nil {
		log.Fatal(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			radio, ok := packet.Layer(layers.LayerTypeRadioTap).(*layers.RadioTap)
			if !ok {
				log.Println("packet could not be asserted as radio tap")
				continue
			}
			ether, ok := packet.Layer(layers.LayerTypeDot11).(*layers.Dot11)
			if !ok {
				log.Println("packet could not be asserted as 802.11 frame")
				continue
			}
			frame := Wireless80211Frame{
				Length:           radio.Length,
				TSFT:             radio.TSFT,
				FlagsRadio:       radio.Flags,
				Rate:             radio.Rate,
				ChannelFrequency: radio.ChannelFrequency,
				ChannelFlags:     radio.ChannelFlags,
				FHSS:             radio.FHSS,
				DBMAntennaSignal: radio.DBMAntennaSignal,
				DBMAntennaNoise:  radio.DBMAntennaNoise,
				LockQuality:      radio.LockQuality,
				TxAttenuation:    radio.TxAttenuation,
				DBTxAttenuation:  radio.DBTxAttenuation,
				DBMTxPower:       radio.DBMTxPower,
				Antenna:          radio.Antenna,
				DBAntennaSignal:  radio.DBAntennaSignal,
				DBAntennaNoise:   radio.DBAntennaNoise,
				RxFlags:          radio.RxFlags,
				TxFlags:          radio.TxFlags,
				RtsRetries:       radio.RtsRetries,
				DataRetries:      radio.DataRetries,
				MCS:              radio.MCS,
				AMPDUStatus:      radio.AMPDUStatus,
				VHT:              radio.VHT,
				Type:             fmt.Sprint(ether.Type),
				Proto:            ether.Proto,
				Flags80211:       ether.Flags,
				DurationID:       ether.DurationID,
				Address1:         fmt.Sprint(ether.Address1),
				Address2:         fmt.Sprint(ether.Address2),
				Address3:         fmt.Sprint(ether.Address3),
				Address4:         fmt.Sprint(ether.Address4),
				SequenceNumber:   ether.SequenceNumber,
				FragmentNumber:   ether.FragmentNumber,
				Checksum:         ether.Checksum,
			}

			if _, ok := packet.Layer(layers.LayerTypeDot11MgmtBeacon).(*layers.Dot11MgmtBeacon); ok {
			} else if probe_req, ok := packet.Layer(layers.LayerTypeDot11MgmtProbeReq).(*layers.Dot11MgmtProbeReq); ok {
				frame.Elements = ParseFrameElements(probe_req.LayerContents())
			}

			if RateLimitFrame(frame) {
				continue
			}
			frames <- frame
		}
	}
}
