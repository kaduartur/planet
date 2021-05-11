package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/kaduartur/planet"
	"log"
)

type ConsumerManager struct {
	kafka     *kafka.Consumer
	processor planet.EventsProcessor
}

func NewConsumer(processor planet.EventsProcessor) ConsumerManager {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     "localhost",
		"broker.address.family": "v4",
		"group.id":              "planet-processor",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "latest"})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}

	log.Printf("Created Consumer %v\n", c)
	return ConsumerManager{kafka: c, processor: processor}
}

func (c ConsumerManager) Read() {
	err := c.kafka.SubscribeTopics([]string{"planet-processor"}, nil)
	if err != nil {
		log.Fatal()
	}

	run := true

	for run == true {
		ev := c.kafka.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			log.Printf("Message on %s:\n%s\n", e.TopicPartition, string(e.Value))

			var eventType string
			for _, header := range e.Headers {
				if header.Key == "EVENT_TYPE" {
					eventType = string(header.Value)
				}
			}

			event, exist := c.processor[eventType]
			if !exist {
				continue
			}

			event.Process(e.Value)

		case kafka.Error:
			log.Printf("Error: %v: %v\n", e.Code(), e)
			if e.Code() == kafka.ErrAllBrokersDown {
				run = false
			}
		default:
			fmt.Printf("Ignored %v\n", e)
		}

	}

	fmt.Printf("Closing consumer\n")
	c.kafka.Close()
}
