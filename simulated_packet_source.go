package fh4server

import (
	"bytes"
	"encoding/base64"
	"time"

	"github.com/golang/glog"
)

const (
	fakePacket = "AQAAAHGpEAD2//lF+P9HREg17ETOUljBxAq0vp6jGk" +
		"AHMB9BgzqMPgAOVj/faPM9fXSzv8Yygry3H/w/9+gKvWx7s7vyegc/NSD9P" +
		"lpD7D7ZpEo/43IfQG+03z/mgKtAzkzePwZWCz8L5qFASKvOQXFkYUEAAAAA" +
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJmZGT+ZmRk/mZkZP5mZGT8" +
		"b06vAHRbewF+x7MBXmtjAWWu9QBYF5UDkJRJBbp7fQNDBPjzQQQg8IFO7O6" +
		"QsFz0kBAAAAwAAACADAAACAAAABAAAACMAAAAAAAAAAAAAAJeqw8Lv3WtDx" +
		"9qxxRPPH0EAAAAAAAAAAGVnm0LYSZpC/FaeQvxWnkLMdrLAAACAPwAAAAAA" +
		"AAAAAAAAAAAAAABs2xdEAAAA/wAAAAGBAAAA"
)

// SimulatedPacketSource implements PacketSource but does not have any external
// dependencies and always sends an identical packet at a regular interval.
//
// Note: could add some randomization to the packet values in here, but
// honestly this probably isn't going to be seeing much use now.
type SimulatedPacketSource struct {
	interval    time.Duration
	packetBytes []byte
}

// NewSimulatedPacketSource sets up a SimulatedPacketSource to send the packet
// at regular interval.
func NewSimulatedPacketSource(interval time.Duration) *SimulatedPacketSource {
	packetBytes, err := base64.StdEncoding.DecodeString(fakePacket)
	if err != nil {
		glog.Errorf("invalid input for fake packet source: %v", err)
	}
	return &SimulatedPacketSource{
		interval:    interval,
		packetBytes: packetBytes,
	}
}

// ReadNextPacket blocks for a period of time defined by `interval` before
// returning the test packet.
func (packetSource *SimulatedPacketSource) ReadNextPacket() *bytes.Buffer {
	time.Sleep(packetSource.interval)
	return bytes.NewBuffer(packetSource.packetBytes)
}
