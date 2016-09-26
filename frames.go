package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
)

func sniffInterface(iface string, frames chan Wireless80211Frame) {
	if handle, err := pcap.OpenLive(iface, 1600, true, 1); err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"interface": iface,
		}).Fatal("failed to open pcap handler")
		//} else if err := handle.SetBPFFilter(""); err != nil {
		//	log.WithFields(log.Fields{
		//		"error":     err,
		//		"interface": iface,
		//	}).Fatal("failed to set network filter")
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			radio, ok := packet.Layer(layers.LayerTypeRadioTap).(*layers.RadioTap)
			if !ok {
				log.Info("frame could not be asserted as radio tap")
				continue
			}

			ether, ok := packet.Layer(layers.LayerTypeDot11).(*layers.Dot11)
			if !ok {
				log.Info("packet could not be asserted as 802.11 frame")
				continue
			}

			frame := createFrame(packet, radio, ether)

			if rateLimit(frame) {
				continue
			}

			frames <- frame
		}
	}
}

func createFrame(packet gopacket.Packet, radio *layers.RadioTap, ether *layers.Dot11) Wireless80211Frame {
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

	frame.parseElements(packet, ether)

	return frame
}

func (frame *Wireless80211Frame) parseElements(packet gopacket.Packet, ether *layers.Dot11) {
	// since ether.Type is known, try to lookup the right one
	//fmt.Println(ether.Type, uint8(ether.Type))
	//fmt.Println(packet.Layers())
	switch ether.Type {
	case 0: //layers.LayerTypeDot11MgmtBeacon:
		//fmt.Println("its a mgt")
		//if beacon, ok := packet.Layer(layers.LayerTypeDot11MgmtBeacon).(*layers.Dot11MgmtBeacon); ok {
		//	fmt.Println("got a beacon", beacon)
		//	elements = ParseBeaconFrame()
		//}
	}
	if probegreq, ok := packet.Layer(layers.LayerTypeDot11MgmtProbeReq).(*layers.Dot11MgmtProbeReq); ok {
		//fmt.Println(ether.NextLayerType())
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
}

func streamFrames(frames chan Wireless80211Frame, wave_host string) {
	var ws net.Conn
	for {
		ws = dialWave(wave_host, frames)
		for frame := range frames {
			if rateLimit(frame) {
				continue
			}

			flat, err := json.Marshal(frame)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Warn("failed to marshal frame")
				continue
			}

			if _, err := ws.Write([]byte(flat)); err != nil {
				ws.Close()
				log.WithFields(log.Fields{
					"error":     err,
					"wave_host": wave_host,
				}).Error("failed to send frame, redailing Wave")
				break
			}
		}
	}
}
