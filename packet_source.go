package fh4server

import (
	"bytes"
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"
)

var (
	udpListenPort = flag.Int("udp_listen_port", 10001, "port to listen for game messages on.")
)

// PacketSource is implemented by anything that can provide game packets
// (real, or fake) to the microservice.
type PacketSource interface {
	ReadNextPacket() *bytes.Buffer
}

// FH4Game implements PacketSource and uses Forza Horizon 4's Data Out setting
// as a source of data.
type FH4Game struct {
	udpConn *net.UDPConn
	buf     []byte
}

// NewFH4Game sets up the UDP server for listening to game data messages being
// broadcast by the game.
func NewFH4Game() *FH4Game {

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", *udpListenPort))
	if err != nil {
		glog.Fatalf("failed to resolve udp address: %v", err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		glog.Fatalf("failed to listen to udp: %v", err)
	}

	glog.Infof("listening on %v", udpAddr.String())

	buf := make([]byte, 1024)
	return &FH4Game{udpConn, buf}
}

// Close should be called when the FH4Game is no longer needed. This closes the
// udp connection.
func (fh4Game *FH4Game) Close() {
	fh4Game.udpConn.Close()
	glog.Infof("done")
}

// ReadNextPacket blocks until receiving a UDP packet, reads it, and returns it.
func (fh4Game *FH4Game) ReadNextPacket() *bytes.Buffer {

	n, _, err := fh4Game.udpConn.ReadFromUDP(fh4Game.buf)
	if err != nil {
		glog.Errorf("failed to read from udp: %v", err)
	}

	packetBytes := fh4Game.buf[0:n]
	return bytes.NewBuffer(packetBytes)
}
