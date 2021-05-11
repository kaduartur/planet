package mock

import "github.com/kaduartur/planet"

type KafkaWrite struct {
	Err error
}

func (k KafkaWrite) Write(topic string, eventType planet.EventType, msg []byte) error {
	return k.Err
}
