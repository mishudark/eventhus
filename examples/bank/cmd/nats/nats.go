package main

import (
	"log"

	n "github.com/mishudark/eventhus/eventbus/nats"

	"github.com/nats-io/go-nats"
)

func main() {
	end := make(chan bool)
	client, err := n.NewClient("nats://ruser:T0pS3cr3t@localhost:6222", false)

	nc, err := client.Options.Connect()
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	subj := "bank.account"
	nc.Subscribe(subj, func(msg *nats.Msg) {
		log.Println(msg)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]\n", subj)
	<-end
	log.Printf("Listening on [%s]\n", subj)

}
