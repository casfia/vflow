//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    network.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    02/01/2017
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------

package packet

import (
	"errors"
	"net"
)

// IPv4Header represents an IPv4 header
type IPv4Header struct {
	Version  int    `json:"version"`         // protocol version
	TOS      int    `json:"tos"`             // type-of-service
	TotalLen int    `json:"total_len"`       // packet total length
	ID       int    `json:"id"`              // identification
	Flags    int    `json:"flags"`           // flags
	FragOff  int    `json:"fragment_offset"` // fragment offset
	TTL      int    `json:"ttl"`             // time-to-live
	Protocol int    `json:"protocol"`        // next protocol
	Checksum int    `json:"check_sum"`       // checksum
	Src      string `json:"src_ip"`          // source address
	Dst      string `json:"dst_ip"`          // destination address
}

// IPv6Header represents an IPv6 header
type IPv6Header struct {
	Version      int    `json:"version"`        // protocol version
	TrafficClass int    `json:"traffic_class"`  // traffic class
	FlowLabel    int    `json:"flow_label"`     // flow label
	PayloadLen   int    `json:"payload_length"` // payload length
	NextHeader   int    `json:"next_layer"`     // next header
	HopLimit     int    `json:"hop_limit"`      // hop limit
	Src          string `json:"src_ip"`         // source address
	Dst          string `json:"dst_ip"`         // destination address
}

const (
	// IPv4HLen is IPv4 header length size
	IPv4HLen = 20

	// IPv6HLen is IPv6 header length size
	IPv6HLen = 40
)

var (
	errShortIPv4HeaderLength = errors.New("short ipv4 header length")
	errShortIPv6HeaderLength = errors.New("short ipv6 header length")
	errShortEthernetLength   = errors.New("short ethernet header length")
	errUnknownTransportLayer = errors.New("unknown transport layer")
	errUnknownL3Protocol     = errors.New("unknown network layer protocol")
)

func (p *Packet) decodeNextLayer() error {

	var (
		proto int
		len   int
	)

	switch p.L3.(type) {
	case IPv4Header:
		proto = p.L3.(IPv4Header).Protocol
	case IPv6Header:
		proto = p.L3.(IPv6Header).NextHeader
	default:
		return errUnknownL3Protocol
	}

	switch proto {
	case IANAProtoICMP:
		icmp, err := decodeICMP(p.Data)
		if err != nil {
			return err
		}

		p.L4 = icmp
		len = 4
	case IANAProtoTCP:
		tcp, err := decodeTCP(p.Data)
		if err != nil {
			return err
		}

		p.L4 = tcp
		len = 20
	case IANAProtoUDP:
		udp, err := decodeUDP(p.Data)
		if err != nil {
			return err
		}

		p.L4 = udp
		len = 8
	default:
		return errUnknownTransportLayer
	}

	p.Data = p.Data[len:]

	return nil
}

func (p *Packet) decodeIPv6Header() error {
	if len(p.Data) < IPv6HLen {
		return errShortIPv6HeaderLength
	}

	var (
		src net.IP = p.Data[8:24]
		dst net.IP = p.Data[24:40]
	)

	p.L3 = IPv6Header{
		Version:      int(p.Data[0]) >> 4,
		TrafficClass: int(p.Data[0]&0x0f)<<4 | int(p.Data[1])>>4,
		FlowLabel:    int(p.Data[1]&0x0f)<<16 | int(p.Data[2])<<8 | int(p.Data[3]),
		PayloadLen:   int(uint16(p.Data[4])<<8 | uint16(p.Data[5])),
		NextHeader:   int(p.Data[6]),
		HopLimit:     int(p.Data[7]),
		Src:          src.String(),
		Dst:          dst.String(),
	}

	p.Data = p.Data[IPv6HLen:]

	return nil
}

func (p *Packet) decodeIPv4Header() error {
	if len(p.Data) < IPv4HLen {
		return errShortIPv4HeaderLength
	}

	var (
		src net.IP = p.Data[12:16]
		dst net.IP = p.Data[16:20]
	)

	p.L3 = IPv4Header{
		Version:  int(p.Data[0] & 0xf0 >> 4),
		TOS:      int(p.Data[1]),
		TotalLen: int(p.Data[2])<<8 | int(p.Data[3]),
		ID:       int(p.Data[4])<<8 | int(p.Data[5]),
		Flags:    int(p.Data[6] & 0x07),
		TTL:      int(p.Data[8]),
		Protocol: int(p.Data[9]),
		Checksum: int(p.Data[10])<<8 | int(p.Data[11]),
		Src:      src.String(),
		Dst:      dst.String(),
	}

	p.Data = p.Data[IPv4HLen:]

	return nil
}
