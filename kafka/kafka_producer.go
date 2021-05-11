package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/kaduartur/planet"
	"log"
)

type ProducerManager struct {
	kafka *kafka.Producer
}

func NewProducer() *ProducerManager {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}

	return &ProducerManager{kafka: p}
}

func (pm *ProducerManager) Write(topic string, eventType planet.EventType, msg []byte) error {
	deliveryChan := make(chan kafka.Event)

	err := pm.kafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
		Headers:        []kafka.Header{{Key: "EVENT_TYPE", Value: []byte(eventType.String())}},
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return nil
}
