package fh4server

import (
	"flag"
	"time"
)

var (
	filterPause = flag.Bool("filter_pause", true, "when enabled, incoming data is ignored while the game is paused.")
)

// Run is the main entry point for the service. `store` and `packetSource` are
// interfaces that represent the service's source of input, and output
// destination.
func Run(whitelist Whitelist, packetSource PacketSource, store PacketStore) {
	for {
		packetBuf := packetSource.ReadNextPacket()
		timestamp := time.Now()
		packet := ParseBuf(packetBuf, whitelist)

		// TODO: once finishers are implemented, make is_race_on a bool.
		// might not actually be super easy way to do this because tags are a
		// map[string]string
		if *filterPause && packet.Tags["is_race_on"] == "0" {
			continue
		}
		go store.WritePacket(packet, timestamp)
	}
}
