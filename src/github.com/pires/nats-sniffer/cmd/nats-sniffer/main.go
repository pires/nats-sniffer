package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats"
)

var (
	server = flag.String("server", "localhost:4222", "NATS adress (host:port) to connect to for sniffing")
	//tls      = flag.Bool("tls", false, "Enable TLS connection to NATS server")
	//duration = flag.Duration("duration", (10 * time.Second), "How long should we sniff")
	subject = flag.String("subject", "*", "The subject to sniff")
)

func main() {
	flag.Parse()

	// connect to specified NATS server
	servers := []string{"nats://" + *server}
	// Setup options to include all servers in the cluster
	opts := nats.DefaultOptions
	// configure internal buffered chan; the default is 64k and
	// we end up using too much memory because of that;
	opts.SubChanLen = 4096
	opts.Servers = servers
	opts.MaxReconnect = 5
	opts.ReconnectWait = (5 * time.Second)
	// go
	nc, err := opts.Connect()
	if err != nil {
		panic(err)
	}

	// wait for explicit termination
	fmt.Printf("Sniffing on [%s] for subject [%s]. press CTRL+C or CTRL+Z to quit...\n", *server, *subject)

	// subscribe to subject
	subscription, err := nc.Subscribe(*subject, func(m *nats.Msg) {
		fmt.Printf("[%s] -> %s\n", m.Subject, string(m.Data))
	})
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	subscription.Unsubscribe()
	nc.Close()

	fmt.Println("Sniffer terminated.")
}
