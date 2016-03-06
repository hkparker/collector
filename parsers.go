package main

import (
	log "github.com/Sirupsen/logrus"
)

func ParseFrameElements(stream []byte) (elements map[string][]byte) {
	for len(stream) > 0 {
		field_id, remainder := stream[0], stream[1:]
		stream = remainder

		field, ok := ELEMENT_IDS[field_id]
		if !ok {
			log.WithFields(log.Fields{
				"id": field_id,
			}).Warn("unknown element id")
		}

		field_len, remainder := stream[0], stream[1:]
		stream = remainder
		if field_len == 0 {
			continue
		}

		field_data, remainder := stream[:field_len], stream[field_len:]
		stream = remainder

		elements[field] = field_data
	}
	return
}
