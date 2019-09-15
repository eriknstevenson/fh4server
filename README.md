# Forza Horizon 4 Telemetry Dashboard

## Current Status

Currently, the UDP server is implemented and able to receive/parse messages from the game. These messages are written to InfluxDB where they are able to be queried and visualized.

Next step is to come up with a better story around configuration and influx setup.

## How to run

1. Setup and run Influx.
2. go run cmd/fh4server.go [options]

#### With docker compose

> **Note**: docker compose configuration is still a work in progress and not fully functional.

```
$ docker-compose up -d
```

### Configure your game to send out data.

Configure the "data out" settings in Forza Horizon 4 as follows:

```
DATA OUT: on
IP address: <your computer's local ip address>
UDP port: 10001
```

### Configure your instance of Influx

TODO

### That's it

In a browser, go to http://localhost:9999 and log in to view the dashboards.

## How to develop and run tests

Build and run tests as per usual with built-in `go` commands.

### If Playing Forza on Windows!

This has only been tested while playing on Xbox One. Using this while playing Forza on Windows may require adjustment of firewall settings.
