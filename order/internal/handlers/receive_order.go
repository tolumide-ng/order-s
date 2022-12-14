package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/config"
	"order-service/events"
	"order-service/models"
	"order-service/publisher"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// ReceiveOrder handler will accept an order, validate the payload and publish an OrderReceived event to kafka.
// returns a HRRP 201 status code indicating an order was created
//
// Example cURL payload (localhost)
// $ curl -v -H "Content-Type: application/json -d ”{"id":"6e042f29-350b-4d51-8849-5e36456dfa48","products":[{"productCode":"12345","quantity":2}],"customer":{"firstName":"Tom","lastName":"Hardy","emailAddress":"tom.hardy@email.com","shippingAddress":{"line1":"123 Anywhere St","city":"Anytown","state":"AL","postalCode":"12345"}}}' http://localhost:8080/orders'"
func ReceiveOrder(w http.ResponseWriter, r *http.Request) {
	var o models.Order
	// Create a new ID for the order
	o.ID = uuid.New()

	var err error

	if err = json.NewDecoder(r.Body).Decode(&o); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	log.WithField("order", o).Info("received new order")

	if err = validate(o); err != nil {
		log.WithField("orderID", o.ID).Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	e := translateOrderToEvent(o)

	log.WithField("event", e).Info("transformed order to event")

	if err = publisher.PublishEvent(e, config.OrderReceivedTopicName); err != nil {
		log.WithField("event", e).Info("published event")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	log.WithField("event", e).Info("")
	w.WriteHeader(http.StatusCreated)
}

// Validates the order payload has the necessary information and returns an error if it is invalid
func validate(o models.Order) error {
	if len(o.Products) == 0 {
		return fmt.Errorf("there are no productsd in the order")
	}

	for i, p := range o.Products {
		if len(p.ProductCode) == 0 {
			return fmt.Errorf("product code is required for product [%d]", i)
		}

		if len(o.Customer.ShippingAddress.Line1) == 0 {
		}

		if p.Quantity <= 0 {
			return fmt.Errorf("quantity should be greater than zero for product [%s]", p.ProductCode)
		}
	}

	if len(o.Customer.EmailAddress) == 0 {
		return fmt.Errorf("email address is required")
	}

	if len(o.Customer.ShippingAddress.Line1) == 0 {
		return fmt.Errorf("shipping address city is required")
	}

	if len(o.Customer.ShippingAddress.PostalCode) == 0 {
		return fmt.Errorf("shipping address postal code is required")
	}

	return nil
}

func translateOrderToEvent(o models.Order) events.Event {
	var event = events.OrderReceived{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: o,
	}

	return event
}
