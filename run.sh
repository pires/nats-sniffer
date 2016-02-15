#!/bin/sh

# set environment
export NATS_SERVER=${NATS_SERVER:-nats:4222}

/nats-sniffer -nats $NATS_SERVER
