package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"

	"github.com/pires/nats-sniffer/sniffer"
)

var (
	port = flag.Int("port", 8080, "Port to listen to for client requests")
	nats = flag.String("nats", "localhost:4222", "NATS address (user:pass@host:port) to connect to for sniffing")
)

// Broker handles message delivery to all connected clients
type Broker struct {
	sniffer *sniffer.Sniffer
}

// ServeHTTP handles GET /sniff/ URL.
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// make sure that the writer supports flushing.
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// get the subject from path
	subject := r.URL.Query().Get("subject")
	fmt.Printf("Incoming client [%s].\n", subject)

	// incoming message handler
	handlerFn := func(msg string) {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		f.Flush()
	}

	// sniff
	handlerId, err := b.sniffer.Sniff(subject, handlerFn)
	if err != nil {
		fmt.Fprintf(w, "There was an error while sniffing subject [%s]: %s\n", subject, err.Error())
		f.Flush()
		return
	}

	// wait until client disconnects
	closing := w.(http.CloseNotifier).CloseNotify()
	<-closing
	b.sniffer.Unsniff(subject, handlerId)
	fmt.Printf("Client gone [%s].\n", subject)
}

// MainPageHandler renders the main page.
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// read in the template with our SSE JavaScript code.
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Errorf("WTF dude, error parsing your template.")
		return
	}

	// Render the template, writing to `w`.
	t.Execute(w, "Duder")

	// Done.
	fmt.Println("Finished HTTP request at ", r.URL.Path)
}

func main() {
	flag.Parse()

	s := sniffer.NewSniffer(*nats)
	if err := s.Start(); err != nil {
		panic(err)
	}

	// Make a new Broker instance
	b := &Broker{sniffer: s}

	// handlers
	http.Handle("/sniff/", b)
	http.Handle("/", http.HandlerFunc(MainPageHandler))

	// wait for Ctrl-c to stop server
	fmt.Println("Service is running, press CTRL+C or CTRL+Z to quit...")
	if err := http.ListenAndServe(fmt.Sprint(":", *port), nil); err != nil {
		panic(err)
	}

	fmt.Println("Sniffer terminated.")
}
