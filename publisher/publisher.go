package publisher

import (
	"encoding/json"

	"order-service/config"
	"order-service/events"

	log "github.com/sirupsen/logrus"
	kafka "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// PublisherEvent will publish the specified event to the messaging system (currently running on localhost)
func PublisherEvent(event events.Event, topic string) error {
	log.WithField("event", event).Info("attempting to publish event")

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   config.BrokerAddress(),
		"socket.timeout.ms":   3000,
		"delivery.timeout.ms": 3000,
	})

	if err != nil {
		return err
	}

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used
	deliveryChan := make(chan kafka.Event)

	var value []byte
	if value, err = json.Marshal(event); err != nil {
		return err
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	log.WithField("Name", *m.TopicPartition.Topic).
		WithField("Partition", m.TopicPartition.Partition).
		WithField("PartitionOffset", m.TopicPartition.Offset).
		Infof("Delivered message to topic")

	close(deliveryChan)

	return nil
}
