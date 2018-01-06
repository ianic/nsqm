package main

import (
	"fmt"
	"log"
	"time"

	"github.com/minus5/nsqm"
	"github.com/minus5/nsqm/discovery/consul"
	nsq "github.com/nsqio/go-nsq"
)

const (
	topic   = "hello_world"
	channel = "app"
)

func main() {
	// discovery
	dcy, err := consul.New("127.0.0.1:8500")
	if err != nil {
		log.Fatal(err)
	}
	// show discovered configuration
	// la, _ := dcy.NSQLookupdAddresses()
	// na, _ := dcy.NSQDAddress()
	// fmt.Printf("config from consul:\n\tnsqd: %s,\n\tnsqlookupds:%v\n", na, la)

	// configuration with discovery
	cfgr := nsqm.WithDiscovery(dcy)
	// create producer
	producer, err := nsqm.NewProducer(cfgr)
	if err != nil {
		log.Fatal(err)
	}
	// create consumer
	h := &handler{msgs: make(chan string)}
	consumer, err := nsqm.NewConsumer(cfgr, topic, channel, h)
	if err != nil {
		log.Fatal(err)
	}
	// send message with producer
	msg := fmt.Sprintf("Hello Word at %v", time.Now())
	if err := producer.Publish(topic, []byte(msg)); err != nil {
		log.Fatal(err)
	}

	// waith for consumer to receive a message
	log.Printf("received: %s\n", <-h.msgs)

	// cleanup
	producer.Stop()
	consumer.Stop()
}

type handler struct {
	msgs chan string
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	h.msgs <- string(m.Body)
	return nil
}
