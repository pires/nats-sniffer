# nats-sniffer
A simple snigger for [NATS](https://nats.io), the cloud native messaging system.

## Pre-requisites

* Go 1.5.x or newer
* `make`
* Optionally, [`gb`](http://getgb.io)

## Build

The following will clean any previously built artifacts, run tests and generate a binary for your platform:
```
make
```

If you're on Mac OS but want to build for Linux:
```
make linux
```

## Usage

### Sniffer

```
bin/nats-sniffer --help
```

Example:
```
bin/nats-sniffer -port 8080 -nats 192.168.99.100:4222
```

### Client

```
curl "<HOST>:<PORT>/sniff/?subject=<SUBJECT>"
```

Example:
```
curl "localhost:8080/sniff/?subject=device.*.connection"
{"device": {"id": "simulator-1","mac": "simulator-1","firmware": "1.0.0","eventType": "CONNECTED"}
{"device": {"id": "simulator-1","mac": "simulator-1"},"eventType": "DISCONNECTED"}
```

## Vendored Dependencies

* `github.com/nats-io/nats`
