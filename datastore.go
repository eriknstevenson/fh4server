package fh4server

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/influxdata/influxdb-client-go"
)

var (
	influxAddr  = flag.String("influx_addr", "http://influxdb:9999", "address of InfluxDB server")
	influxToken = flag.String("influx_token", "", "token to use to connect to InfluxDB server.")
	bucketName  = flag.String("influx_bucket_name", "data", "packets will be written to db with this bucket name")
	orgName     = flag.String("influx_org_name", "fh4server", "packets will be written to db with this org name")
	timeout     = flag.Int("influx_timeout", 2, "writes to db will timeout after this amount of seconds")
)

// PacketStore is a general purpose interface for anything that can store packets.
type PacketStore interface {
	WritePacket(packet Packet, timestamp time.Time)
}

// InfluxStore implements PacketStore and uses InfluxDB as a backend.
type InfluxStore struct {
	influx *influxdb.Client
}

// NewInfluxStore connects to influx and returns a handle to the interface.
func NewInfluxStore() (*InfluxStore, error) {
	influx, err := influxdb.New(*influxAddr, *influxToken)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to influxdb: %v", err)
	}
	glog.Infof("connected to influxdb: %s", *influxAddr)
	return &InfluxStore{influx}, nil
}

// Close should be called when the database is no longer needed.
func (dataStore *InfluxStore) Close() {
	dataStore.influx.Close()
}

// WritePacket writes a packet to the database
func (dataStore *InfluxStore) WritePacket(packet Packet, timestamp time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	// TODO(erik): double check that this is being done correctly here. What happens if the timeout is exceeded?
	row := influxdb.NewRowMetric(packet.Fields, "fh4", packet.Tags, timestamp)

	// The actual write..., this method can be called concurrently.
	err := dataStore.influx.Write(ctx, *bucketName, *orgName, row)
	if err != nil {
		glog.Warningf("failed to write packet to db: %v", err)
	}
}
