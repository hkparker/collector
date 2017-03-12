package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/hkparker/Wave/models"
	"net"
)

func sniffInterface(iface string, frames chan models.Wireless80211Frame) {
	if handle, err := pcap.OpenLive(iface, 1600, true, 1); err != nil {
		log.WithFields(log.Fields{
			"at":        "sniffInterface",
			"error":     err,
			"interface": iface,
		}).Fatal("failed to open pcap handler")
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			radio, ok := packet.Layer(layers.LayerTypeRadioTap).(*layers.RadioTap)
			if !ok {
				log.WithFields(log.Fields{
					"at": "sniffInterface",
				}).Info("frame could not be asserted as radio tap")
				continue
			}

			ether, ok := packet.Layer(layers.LayerTypeDot11).(*layers.Dot11)
			if !ok {
				log.WithFields(log.Fields{
					"at": "sniffInterface",
				}).Info("packet could not be asserted as 802.11 frame")
				continue
			}

			frame := createFrame(packet, radio, ether, iface)

			if rateLimit(frame) {
				continue
			}

			frames <- frame
		}
	}
}

func createFrame(packet gopacket.Packet, radio *layers.RadioTap, ether *layers.Dot11, iface string) models.Wireless80211Frame {
	frame := models.Wireless80211Frame{
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
		Interface:        iface,
	}

	frame.ParseElements(packet, ether)

	return frame
}

func streamFrames(frames chan models.Wireless80211Frame, wave_host string) {
	var ws net.Conn
	for {
		if !local {
			ws = dialWave(wave_host, frames)
		}
		for frame := range frames {
			flat, err := json.Marshal(frame)
			if err != nil {
				log.WithFields(log.Fields{
					"at":    "streamFrames",
					"error": err,
				}).Warn("failed to marshal frame")
				continue
			}

			if print {
				fmt.Println(string(flat))
			}
			if !local {
				if _, err := ws.Write([]byte(flat)); err != nil {
					ws.Close()
					log.WithFields(log.Fields{
						"at":        "streamFrames",
						"error":     err,
						"wave_host": wave_host,
					}).Error("failed to send frame, redailing Wave")
					break
				}
			}
		}
	}
}
