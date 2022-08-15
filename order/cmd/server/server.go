package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/bpmericle/Asynchronous-Event-Handling-Using-Microservices-and-Kafka/milestone2/code/logger"
	"github.com/bpmericle/Asynchronous-Event-Handling-Using-Microservices-and-Kafka/milestone2/code/order/internal/handlers"
)

type Server struct {
	Port int
}

// ListenAndServe will start the web server and listen for requests
func (s *Server) ListenAndServe() error {
	// setup CHI router
	r := chi.NewRouter()

	// setup middlewares
	r.Use(middleware.Heartbeat("/ping")) // allows LB to verify service up
	r.Use(middleware.RequestID)          // ensures a request ID is logged
	r.Use(logger.NewStructuredLogger())  // uses structured logging like our app (logs only at debug level)
	r.Use(middleware.Recoverer)

	// setup supported routes
	r.Get("/", handlers.Root)
	r.Get("/health", handlers.Health)

	address := fmt.Sprintf(":%d", s.Port)
	log.WithField("address", address).Info("server starting")

	return http.ListenAndServe(address, r)
}
