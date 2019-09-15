package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	"github.com/narrative/fh4server"
)

var (
	// Note: defaulting these to true temporarily during development
	simulatePacketSource = flag.Bool("simulate_packet_source", false, "when enabled, incoming data is artificial.")
	simulateDataStore    = flag.Bool("simulate_data_store", false, "when enabled, outgoing data is artificial.")
)

func main() {
	flag.Parse()
	glog.Infof("Running with the following options set: ")
	flag.VisitAll(func(f *flag.Flag) {
		glog.Infof("%s = %s", f.Name, f.Value)
	})

	var packetSource fh4server.PacketSource
	if *simulatePacketSource {
		packetSource = fh4server.NewSimulatedPacketSource(2 * time.Second)
	} else {
		fh4Game := fh4server.NewFH4Game()
		defer fh4Game.Close()
		packetSource = fh4Game
	}

	var store fh4server.PacketStore
	if *simulateDataStore {
		store = fh4server.NewSimulatedDataStore(50)
	} else {
		influx, err := fh4server.NewInfluxStore()
		if err != nil {
			glog.Fatalf("failed to connect to db: %v", err)
		}
		defer influx.Close()
		store = influx
	}

	fh4server.Run(fh4server.AllowAll(), packetSource, store)
}
