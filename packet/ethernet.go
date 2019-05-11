//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ethernet.go
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
	"fmt"
)

// Datalink represents layer two IEEE 802.11
type Datalink struct {
	// SrcMAC represents source MAC address
	SrcMAC string `json:"src_mac"`

	// DstMAC represents destination MAC address
	DstMAC string `json:"dst_mac"`

	// Vlan represents VLAN value
	Vlan int `json:"vlan_id"`

	// EtherType represents upper layer type value
	EtherType uint16 `json:"ether_type"`
}

const (
	// EtherTypeARP is Address Resolution Protocol EtherType value
	EtherTypeARP = 0x0806

	// EtherTypeIPv4 is Internet Protocol version 4 EtherType value
	EtherTypeIPv4 = 0x0800

	// EtherTypeIPv6 is Internet Protocol Version 6 EtherType value
	EtherTypeIPv6 = 0x86DD

	// EtherTypeLACP is Link Aggregation Control Protocol EtherType value
	EtherTypeLACP = 0x8809

	// EtherTypeIEEE8021Q is VLAN-tagged frame (IEEE 802.1Q) EtherType value
	EtherTypeIEEE8021Q = 0x8100
)

var (
	errShortEthernetHeaderLength = errors.New("the ethernet header is too small")
)

func (p *Packet) decodeEthernet() error {
	var (
		d   Datalink
		err error
	)

	if len(p.Data) < 14 {
		return errShortEthernetHeaderLength
	}

	d, err = decodeIEEE802(p.Data)
	if err != nil {
		return err
	}

	if d.EtherType == EtherTypeIEEE8021Q {
		vlan := int(p.Data[14])<<8 | int(p.Data[15])
		p.Data[12], p.Data[13] = p.Data[16], p.Data[17]
		p.Data = append(p.Data[:14], p.Data[18:]...)

		d, err = decodeIEEE802(p.Data)
		if err != nil {
			return err
		}
		d.Vlan = vlan
	}

	p.L2 = d
	p.Data = p.Data[14:]

	return nil
}

func decodeIEEE802(b []byte) (Datalink, error) {
	var d Datalink

	if len(b) < 14 {
		return d, errShortEthernetLength
	}

	d.EtherType = uint16(b[13]) | uint16(b[12])<<8

	hwAddrFmt := "%0.2x:%0.2x:%0.2x:%0.2x:%0.2x:%0.2x"

	if d.EtherType != EtherTypeIEEE8021Q {
		d.DstMAC = fmt.Sprintf(hwAddrFmt, b[0], b[1], b[2], b[3], b[4], b[5])
		d.SrcMAC = fmt.Sprintf(hwAddrFmt, b[6], b[7], b[8], b[9], b[10], b[11])
	}

	return d, nil
}
