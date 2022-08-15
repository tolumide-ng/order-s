package main

import (
	"order-service/config"
	"order-service/events"
	"order-service/models"
	"order-service/publisher"
	"os"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(config.LogLevel())
}

func main() {
	var err error
	var order = models.Order{
		ID: uuid.New(),
	}

	var event = events.OrderReceived{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: order,
	}

	if err = publisher.PublishEvent(event, config.OrderReceivedTopicName); err != nil {
		log.WithField("orderID", order.ID).Error(err.Error())
	} else {
		log.WithField("event", event).Info("published event")
	}
}
