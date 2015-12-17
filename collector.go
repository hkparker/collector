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
	// define the http client here, also use websockets instead of posting
	for {
		frame := <-frames
		flat, err := json.Marshal(frame)
		if err != nil {
			log.Println(err)
		}
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(flat))
		if err != nil {
			log.Println(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err) // recreate the client for a few attempts.  sleep to avoid DDoSing it?  dump the frame chan meanwhile
		} else {
			resp.Body.Close()
		}
	}
}

func RateLimitFrame(frame Wireless80211Frame) bool {
	if frame.Type == "MgmtBeacon" {
		return true // only send an exact match for some beacon properties every second
	}
	return false
}

func main() {
	var iface = flag.String("interface", "mon0", "interface to sniff")
	var endpoint = flag.String("endpoint", "http://127.0.0.1:8080/frames", "server to post packet g to")
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

			// since ether.Type is known, try to lookup the right one
			if _, ok := packet.Layer(layers.LayerTypeDot11MgmtBeacon).(*layers.Dot11MgmtBeacon); ok {
			} else if probegreq, ok := packet.Layer(layers.LayerTypeDot11MgmtProbeReq).(*layers.Dot11MgmtProbeReq); ok {
				frame.Elements = ParseFrameElements(probegreq.LayerContents())
			} else if _, ok := packet.Layer(layers.LayerTypeDot11Data).(*layers.Dot11Data); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataCFAck).(*layers.Dot11DataCFAck); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataCFAckNoData).(*layers.Dot11DataCFAckNoData); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataCFAckPoll).(*layers.Dot11DataCFAckPoll); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataCFAckPollNoData).(*layers.Dot11DataCFAckPollNoData); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataCFPoll).(*layers.Dot11DataCFPoll); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataCFPollNoData).(*layers.Dot11DataCFPollNoData); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataNull).(*layers.Dot11DataNull); ok {
				//} else if _, ok := packet.Layer(layers.LayerTypeDot11DataQOS).(*layers.Dot11DataQOS); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataQOSCFAckPollNoData).(*layers.Dot11DataQOSCFAckPollNoData); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataQOSCFPollNoData).(*layers.Dot11DataQOSCFPollNoData); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataQOSData).(*layers.Dot11DataQOSData); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataQOSDataCFAck).(*layers.Dot11DataQOSDataCFAck); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataQOSDataCFAckPoll).(*layers.Dot11DataQOSDataCFAckPoll); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataQOSDataCFPoll).(*layers.Dot11DataQOSDataCFPoll); ok {
			} else if _, ok := packet.Layer(layers.LayerTypeDot11DataQOSNull).(*layers.Dot11DataQOSNull); ok {
			}

			if RateLimitFrame(frame) {
				continue
			}
			frames <- frame
		}
	}
}
