//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    flow_counter.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    08/08/2018
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

package sflow

import (
	"github.com/VerizonDigital/vflow/vlogger"
	"io"
)

const (
	// SFGenericInterfaceCounters is Generic interface counters - see RFC 2233
	SFGenericInterfaceCounters = 1

	// SFEthernetInterfaceCounters is Ethernet interface counters - see RFC 2358
	SFEthernetInterfaceCounters = 2

	// SFTokenRingInterfaceCounters is Token ring counters - see RFC 1748
	SFTokenRingInterfaceCounters = 3

	// SF100BaseVGInterfaceCounters is 100 BaseVG interface counters - see RFC 2020
	SF100BaseVGInterfaceCounters = 4

	// SFVLANCounters is VLAN counters
	SFVLANCounters = 5

	// SFProcessorCounters is processor counters
	SFProcessorCounters = 1001

	SFLCOUNTERS_80211         = 6
	SFLCOUNTERS_LACP          = 7
	SFLCOUNTERS_SFP           = 10
	SFLCOUNTERS_RADIO         = 1002
	SFLCOUNTERS_OFPORT        = 1004
	SFLCOUNTERS_PORTNAME      = 1005
	SFLCOUNTERS_HOST_HID      = 2000 /* host id */
	SFLCOUNTERS_ADAPTORS      = 2001 /* host adaptors */
	SFLCOUNTERS_HOST_PAR      = 2002 /* host parent */
	SFLCOUNTERS_HOST_CPU      = 2003 /* host cpu  */
	SFLCOUNTERS_HOST_MEM      = 2004 /* host memory  */
	SFLCOUNTERS_HOST_DSK      = 2005 /* host storage I/O  */
	SFLCOUNTERS_HOST_NIO      = 2006 /* host network I/O */
	SFLCOUNTERS_HOST_IP       = 2007
	SFLCOUNTERS_HOST_ICMP     = 2008
	SFLCOUNTERS_HOST_TCP      = 2009
	SFLCOUNTERS_HOST_UDP      = 2010
	SFLCOUNTERS_HOST_VRT_NODE = 2100 /* host virt node */
	SFLCOUNTERS_HOST_VRT_CPU  = 2101 /* host virt cpu */
	SFLCOUNTERS_HOST_VRT_MEM  = 2102 /* host virt mem */
	SFLCOUNTERS_HOST_VRT_DSK  = 2103 /* host virt storage */
	SFLCOUNTERS_HOST_VRT_NIO  = 2104 /* host virt network I/O */
	SFLCOUNTERS_JVM           = 2105 /* java runtime */
	SFLCOUNTERS_JMX           = 2106 /* java JMX stats */
	SFLCOUNTERS_MEMCACHE      = 2200 /* memcached (deprecated) */
	SFLCOUNTERS_HTTP          = 2201 /* http */
	SFLCOUNTERS_APP           = 2202
	SFLCOUNTERS_APP_RESOURCE  = 2203
	SFLCOUNTERS_MEMCACHE2     = 2204 /* memcached */
	SFLCOUNTERS_VDI           = 2205
	SFLCOUNTERS_APP_WORKERS   = 2206
	SFLCOUNTERS_OVSDP         = 2207
	SFLCOUNTERS_HOST_GPU_NVML = (5703 << 12) + 1 /* = 23359489 */
	SFLCOUNTERS_BCM_TABLES    = (4413 << 12) + 3
)

// GenericInterfaceCounters represents Generic Interface Counters RFC2233
type GenericInterfaceCounters struct {
	Index               uint32
	Type                uint32
	Speed               uint64
	Direction           uint32
	Status              uint32
	InOctets            uint64
	InUnicastPackets    uint32
	InMulticastPackets  uint32
	InBroadcastPackets  uint32
	InDiscards          uint32
	InErrors            uint32
	InUnknownProtocols  uint32
	OutOctets           uint64
	OutUnicastPackets   uint32
	OutMulticastPackets uint32
	OutBroadcastPackets uint32
	OutDiscards         uint32
	OutErrors           uint32
	PromiscuousMode     uint32
}

// EthernetInterfaceCounters represents Ethernet Interface Counters RFC2358
type EthernetInterfaceCounters struct {
	AlignmentErrors           uint32
	FCSErrors                 uint32
	SingleCollisionFrames     uint32
	MultipleCollisionFrames   uint32
	SQETestErrors             uint32
	DeferredTransmissions     uint32
	LateCollisions            uint32
	ExcessiveCollisions       uint32
	InternalMACTransmitErrors uint32
	CarrierSenseErrors        uint32
	FrameTooLongs             uint32
	InternalMACReceiveErrors  uint32
	SymbolErrors              uint32
}

// TokenRingCounters represents Token Ring Counters - see RFC 1748
type TokenRingCounters struct {
	LineErrors         uint32
	BurstErrors        uint32
	ACErrors           uint32
	AbortTransErrors   uint32
	InternalErrors     uint32
	LostFrameErrors    uint32
	ReceiveCongestions uint32
	FrameCopiedErrors  uint32
	TokenErrors        uint32
	SoftErrors         uint32
	HardErrors         uint32
	SignalLoss         uint32
	TransmitBeacons    uint32
	Recoveries         uint32
	LobeWires          uint32
	Removes            uint32
	Singles            uint32
	FreqErrors         uint32
}

// VGCounters represents 100 BaseVG interface counters - see RFC 2020
type VGCounters struct {
	InHighPriorityFrames    uint32
	InHighPriorityOctets    uint64
	InNormPriorityFrames    uint32
	InNormPriorityOctets    uint64
	InIPMErrors             uint32
	InOversizeFrameErrors   uint32
	InDataErrors            uint32
	InNullAddressedFrames   uint32
	OutHighPriorityFrames   uint32
	OutHighPriorityOctets   uint64
	TransitionIntoTrainings uint32
	HCInHighPriorityOctets  uint64
	HCInNormPriorityOctets  uint64
	HCOutHighPriorityOctets uint64
}

// VlanCounters represents VLAN Counters
type VlanCounters struct {
	ID               uint32
	Octets           uint64
	UnicastPackets   uint32
	MulticastPackets uint32
	BroadcastPackets uint32
	Discards         uint32
}

// ProcessorCounters represents Processor Information
type ProcessorCounters struct {
	CPU5s       uint32
	CPU1m       uint32
	CPU5m       uint32
	TotalMemory uint64
	FreeMemory  uint64
}

// CounterSample represents the periodic sampling or polling of counters associated with a Data Source
type CounterSample struct {
	SequenceNo   uint32
	SourceIDType byte
	SourceIDIdx  uint32
	RecordsNo    uint32
	Records      map[string]Record
}

func decodeFlowCounter(r io.ReadSeeker) (*CounterSample, error) {
	var (
		cs          = new(CounterSample)
		rTypeFormat uint32
		rTypeLength uint32
		err         error
	)

	if err = cs.unmarshal(r); err != nil {
		return nil, err
	}

	cs.Records = make(map[string]Record)

	for i := uint32(0); i < cs.RecordsNo; i++ {
		if err = read(r, &rTypeFormat); err != nil {
			return nil, err
		}
		if err = read(r, &rTypeLength); err != nil {
			return nil, err
		}

		switch rTypeFormat {

		case SFGenericInterfaceCounters:
			d, err := decodeGenericIntCounters(r)
			if err != nil {
				return cs, err
			}
			cs.Records["GenericInt"] = d
		case SFEthernetInterfaceCounters:
			d, err := decodeEthIntCounters(r)
			if err != nil {
				return cs, err
			}
			cs.Records["EthernetInt"] = d
		case SFTokenRingInterfaceCounters:
			d, err := decodeTokenRingCounters(r)
			if err != nil {
				return cs, err
			}
			cs.Records["TokenRingInt"] = d
		case SF100BaseVGInterfaceCounters:
			d, err := decodeVGCounters(r)
			if err != nil {
				return cs, err
			}
			cs.Records["SF100VgInt"] = d
		case SFVLANCounters:
			d, err := decodeVlanCounters(r)
			if err != nil {
				return cs, err
			}
			cs.Records["VLAN"] = d
		case SFProcessorCounters:
			d, err := decodedProcessorCounters(r)
			if err != nil {
				return cs, err
			}
			cs.Records["Processor"] = d
		default:
			vlogger.Logger.Printf("unknow counter, typeFormat is %d", rTypeFormat)

			r.Seek(int64(rTypeLength), 1)
		}
	}

	return cs, nil
}

func decodeGenericIntCounters(r io.Reader) (*GenericInterfaceCounters, error) {
	var gic = new(GenericInterfaceCounters)

	if err := gic.unmarshal(r); err != nil {
		return nil, err
	}

	return gic, nil
}

func (gic *GenericInterfaceCounters) unmarshal(r io.Reader) error {
	var err error

	fields := []interface{}{
		&gic.Index,
		&gic.Type,
		&gic.Speed,
		&gic.Direction,
		&gic.Status,
		&gic.InOctets,
		&gic.InUnicastPackets,
		&gic.InMulticastPackets,
		&gic.InBroadcastPackets,
		&gic.InDiscards,
		&gic.InErrors,
		&gic.InUnknownProtocols,
		&gic.OutOctets,
		&gic.OutUnicastPackets,
		&gic.OutMulticastPackets,
		&gic.OutBroadcastPackets,
		&gic.OutDiscards,
		&gic.OutErrors,
		&gic.PromiscuousMode,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}
func decodeEthIntCounters(r io.Reader) (*EthernetInterfaceCounters, error) {
	var eic = new(EthernetInterfaceCounters)

	if err := eic.unmarshal(r); err != nil {
		return nil, err
	}

	return eic, nil
}

func (eic *EthernetInterfaceCounters) unmarshal(r io.Reader) error {
	var err error

	fields := []interface{}{
		&eic.AlignmentErrors,
		&eic.FCSErrors,
		&eic.SingleCollisionFrames,
		&eic.MultipleCollisionFrames,
		&eic.SQETestErrors,
		&eic.DeferredTransmissions,
		&eic.LateCollisions,
		&eic.ExcessiveCollisions,
		&eic.InternalMACTransmitErrors,
		&eic.CarrierSenseErrors,
		&eic.FrameTooLongs,
		&eic.InternalMACReceiveErrors,
		&eic.SymbolErrors,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}
func decodeTokenRingCounters(r io.Reader) (*TokenRingCounters, error) {
	var tr = new(TokenRingCounters)

	if err := tr.unmarshal(r); err != nil {
		return nil, err
	}

	return tr, nil
}

func (tr *TokenRingCounters) unmarshal(r io.Reader) error {
	var err error

	fields := []interface{}{
		&tr.LineErrors,
		&tr.BurstErrors,
		&tr.ACErrors,
		&tr.AbortTransErrors,
		&tr.InternalErrors,
		&tr.LostFrameErrors,
		&tr.ReceiveCongestions,
		&tr.FrameCopiedErrors,
		&tr.TokenErrors,
		&tr.SoftErrors,
		&tr.HardErrors,
		&tr.SignalLoss,
		&tr.TransmitBeacons,
		&tr.Recoveries,
		&tr.LobeWires,
		&tr.Removes,
		&tr.Singles,
		&tr.FreqErrors,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}

func decodeVGCounters(r io.Reader) (*VGCounters, error) {
	var vg = new(VGCounters)

	if err := vg.unmarshal(r); err != nil {
		return nil, err
	}

	return vg, nil
}

func (vg *VGCounters) unmarshal(r io.Reader) error {
	var err error

	fields := []interface{}{
		&vg.InHighPriorityFrames,
		&vg.InHighPriorityOctets,
		&vg.InNormPriorityFrames,
		&vg.InNormPriorityOctets,
		&vg.InIPMErrors,
		&vg.InOversizeFrameErrors,
		&vg.InDataErrors,
		&vg.InNullAddressedFrames,
		&vg.OutHighPriorityFrames,
		&vg.OutHighPriorityOctets,
		&vg.TransitionIntoTrainings,
		&vg.HCInHighPriorityOctets,
		&vg.HCInNormPriorityOctets,
		&vg.HCOutHighPriorityOctets,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}

func decodeVlanCounters(r io.Reader) (*VlanCounters, error) {
	var vc = new(VlanCounters)

	if err := vc.unmarshal(r); err != nil {
		return nil, err
	}

	return vc, nil
}

func (vc *VlanCounters) unmarshal(r io.Reader) error {
	var err error
	fields := []interface{}{
		&vc.ID,
		&vc.Octets,
		&vc.UnicastPackets,
		&vc.MulticastPackets,
		&vc.BroadcastPackets,
		&vc.Discards,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}

func decodedProcessorCounters(r io.Reader) (*ProcessorCounters, error) {
	var pc = new(ProcessorCounters)

	if err := pc.unmarshal(r); err != nil {
		return nil, err
	}

	return pc, nil
}

func (pc *ProcessorCounters) unmarshal(r io.Reader) error {
	var err error
	fields := []interface{}{
		&pc.CPU5s,
		&pc.CPU1m,
		&pc.CPU5m,
		&pc.TotalMemory,
		&pc.FreeMemory,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}

func (cs *CounterSample) unmarshal(r io.Reader) error {

	var err error

	if err = read(r, &cs.SequenceNo); err != nil {
		return err
	}

	if err = read(r, &cs.SourceIDType); err != nil {
		return err
	}

	buf := make([]byte, 3)
	if err = read(r, &buf); err != nil {
		return err
	}
	cs.SourceIDIdx = uint32(buf[2]) | uint32(buf[1])<<8 | uint32(buf[0])<<16

	err = read(r, &cs.RecordsNo)

	return err
}
