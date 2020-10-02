package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"os"
	"strconv"
)

var (
	broker     = "kafka-broker:9092"
	numServers = 1000
	reports    = "reports"
	repair     = "repair"
)

func setup() {
	if brokerAddress, ok := os.LookupEnv("KAFKA_ADDRESS"); ok {
		broker = brokerAddress
	}

	servers, ok := os.LookupEnv("NUM_SERVERS")
	if !ok {
		return
	}

	var err error
	numServers, err = strconv.Atoi(servers)

	if err != nil {
		log.Fatalf("failed to parse configured number of servers %s %v", servers, err)
	}
}

func main() {
	setup()
	log.Printf("starting data center simulator")
	log.Printf("connecting to kafka %s", broker)
	dc := newDataCenter(numServers)

	go dc.run()
	go consumeRepairRequests(dc.fixRequests)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer p.Close()

	for {
		select {
		case report := <-dc.reports:
			any := anypb.Any{}
			_ = any.MarshalFrom(report)
			value, _ := proto.Marshal(&any)
			_ = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &reports, Partition: kafka.PartitionAny},
				Key:            []byte(report.ServerId),
				Value:          value,
			}, nil)
		}
	}
}

func consumeRepairRequests(workOrders chan<- string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        broker,
		"go.events.channel.enable": true,
		"group.id":                 "fix-requests",
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := c.Subscribe(repair, nil); err != nil {
		log.Fatal(err)
	}

	for {
		ev := c.Poll(100)
		if ev == nil {
			continue
		}

		switch event := ev.(type) {
		case *kafka.Message:
			workOrders <- string(event.Key)
		default:
			fmt.Printf("Read invalid message %v", event)
			continue
		}
	}
}
