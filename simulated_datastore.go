package fh4server

import (
	"time"

	"github.com/golang/glog"
)

// SimulatedDataStore implements the PacketStore interface but does not do
// anything meaningful with the packets it stores (and has no external
// dependencies). It keeps the last `packetsToStore` received.
type SimulatedDataStore struct {
	packetsToStore int
	packets        []Packet
}

// NewSimulatedDataStore sets up and returns a handle to the SimulatedDataStore.
func NewSimulatedDataStore(packetsToStore int) *SimulatedDataStore {
	return &SimulatedDataStore{
		packetsToStore: packetsToStore, packets: make([]Packet, 0)}
}

// WritePacket writes a packet to the database
func (dataStore *SimulatedDataStore) WritePacket(packet Packet, timestamp time.Time) {
	dataStore.packets = append(dataStore.packets, packet)
	// trim down to max length
	dataStore.packets = dataStore.packets[len(dataStore.packets)-dataStore.packetsToStore : len(dataStore.packets)]

	glog.Infof("simulated writing packet at %s", timestamp)
	glog.Infof("packet contents: %v", packet)
}
