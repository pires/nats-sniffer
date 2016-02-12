package sniffer

import (
	"errors"
	"time"

	"github.com/nats-io/nats"
	"github.com/satori/go.uuid"
)

var (
	ERR_NATS_CONN_CLOSED = errors.New("NATS connection is closed.")
)

type SniffedMessageHandler func(msg string)

// Sniffer subscribes to client-request NATS subjects and let them know
// when new messages arrive.
type Sniffer struct {
	natsURL                 string
	natsConn                *nats.Conn
	subjectSubscriptionsMap *SubjectSubscriptionsMap
	subjectHandlersMap      *SubjectHandlersMap
	Quit                    chan struct{}
}

// NewSniffer returns a new sniffer instance.
func NewSniffer(url string) *Sniffer {
	return &Sniffer{
		natsURL:                 url,
		subjectSubscriptionsMap: NewSubjectSubscriptionMap(),
		subjectHandlersMap:      NewSubjectHandlersMap(),
		Quit:                    make(chan struct{}),
	}
}

// Start connects to the specified NATS server and starts managing
// subject subscription and pumping data to registered clients
func (s *Sniffer) Start() error {
	var err error

	// prepare NATS urls
	servers := []string{"nats://" + s.natsURL}

	// setup options to include all servers in the cluster
	opts := nats.DefaultOptions

	// configure internal buffered chan; the default is 64k and
	// we end up using too much memory because of that;
	opts.SubChanLen = 8192
	opts.Servers = servers
	opts.MaxReconnect = 5
	opts.ReconnectWait = (5 * time.Second)

	// connect
	s.natsConn, err = opts.Connect()
	if err != nil {
		return err
	}

	go s.run()

	return nil
}

func (s *Sniffer) run() {
	// run cleanup periodically
	cleanupTick := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-cleanupTick.C:
			// look for subscribed subjects no one cares about anymore
			for subject, subscribers := range s.subjectHandlersMap.Iter() {
				if subscribers.Count() == 0 {
					// remove subscription to subject
					if subscription, ok := s.subjectSubscriptionsMap.Get(subject); ok {
						subscription.Unsubscribe()
						s.subjectSubscriptionsMap.Remove(subject)
						s.subjectHandlersMap.Remove(subject)
					}
				}
			}
			break
		case <-s.Quit:
			cleanupTick.Stop()
			s.natsConn.Close()
			return
		}
	}
}

// Sniff makes sure there's only on subscription per subject
// Every time a new message arrives at one sniffed subject, registered
// client handlers shall be called
func (s *Sniffer) Sniff(subject string, msgHandler SniffedMessageHandler) (string, error) {
	// is subject already being sniffed?
	if !s.subjectSubscriptionsMap.Has(subject) {
		// subject not being sniffed, so sniff
		if s.natsConn.IsClosed() {
			return "", ERR_NATS_CONN_CLOSED
		}

		// subscribe subject
		subscription, err := s.natsConn.Subscribe(subject, func(m *nats.Msg) {
			// call all message handlers interested in the incoming message
			handlers, ok := s.subjectHandlersMap.Get(subject)
			if ok {
				for _, handlerFn := range handlers.Values() {
					handlerFn(string(m.Data))
				}
			}
		})
		if err != nil {
			return "", err
		}

		// store subject subscription
		s.subjectSubscriptionsMap.Set(subject, subscription)
	}

	// store subject subscription handler
	handlers, ok := s.subjectHandlersMap.Get(subject)
	if !ok {
		handlers = NewHandlerMap()
	}
	randomId := uuid.NewV4().String()
	handlers.Set(randomId, msgHandler)
	s.subjectHandlersMap.Set(subject, handlers)

	return randomId, nil
}

func (s *Sniffer) Unsniff(subject string, handlerId string) {
	handlers, ok := s.subjectHandlersMap.Get(subject)
	if ok {
		handlers.Remove(handlerId)
	}
}
